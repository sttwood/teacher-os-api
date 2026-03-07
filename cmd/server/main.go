package main

import (
	"net/http"
	"teacher-os-api/internal/config"
	"teacher-os-api/internal/database"
	authHandler "teacher-os-api/internal/modules/auth/handler"
	authMiddleware "teacher-os-api/internal/modules/auth/middleware"
	authModel "teacher-os-api/internal/modules/auth/model"
	authRepository "teacher-os-api/internal/modules/auth/repository"
	authService "teacher-os-api/internal/modules/auth/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	db := database.NewPostgres(cfg.DBUrl)

	if err := db.AutoMigrate(
		&authModel.User{},
		&authModel.EmailVerificationToken{},
		&authModel.PasswordResetToken{},
		&authModel.RefreshToken{},
	); err != nil {
		panic(err)
	}

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authRepo := authRepository.NewAuthRepository(db)
	authSvc := authService.NewAuthService(authRepo, cfg.JWTSecret)
	authH := authHandler.NewAuthHandler(authSvc)
	authMW := authMiddleware.NewJWTMiddleware(authSvc)

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "db error",
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "db unreachable",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"service": "teacher-os-api",
		})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/logout", authH.Logout)

		auth.GET("/me", authMW.RequireAuth(), authH.Me)

		auth.POST("/verifyEmail/resend", authH.ResendVerifyEmail)
		auth.POST("/verifyEmail/confirm", authH.ConfirmVerifyEmail)

		auth.POST("/forgotPassword", authH.ForgotPassword)
		auth.POST("/resetPassword", authH.ResetPassword)

		auth.POST("/refresh", authH.Refresh)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
