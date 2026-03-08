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
	exportHandler "teacher-os-api/internal/modules/export/handler"
	exportService "teacher-os-api/internal/modules/export/service"
	planHandler "teacher-os-api/internal/modules/plans/handler"
	planWidgetHandler "teacher-os-api/internal/modules/plans/handler"
	planModel "teacher-os-api/internal/modules/plans/model"
	planRepository "teacher-os-api/internal/modules/plans/repository"
	planService "teacher-os-api/internal/modules/plans/service"
	planWidgetService "teacher-os-api/internal/modules/plans/service"
	"teacher-os-api/internal/shared/httpx"
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
		&planModel.Plan{},
		&planModel.PlanWidget{},
		&planModel.PlanVersion{},
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

	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			httpx.Error(c, http.StatusInternalServerError, "DB_ERROR", "db error")
			return
		}

		if err := sqlDB.Ping(); err != nil {
			httpx.Error(c, http.StatusInternalServerError, "DB_UNREACHABLE", "db unreachable")
			return
		}

		httpx.Success(c, http.StatusOK, gin.H{
			"message": "ok",
			"service": "teacher-os-api",
		})
	})

	// auth
	authRepo := authRepository.NewAuthRepository(db)
	authSvc := authService.NewAuthService(authRepo, cfg.JWTSecret)
	authH := authHandler.NewAuthHandler(authSvc)
	authMW := authMiddleware.NewJWTMiddleware(authSvc)

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

	// plans
	planRepo := planRepository.NewPlanRepository(db)
	planSvc := planService.NewPlanService(planRepo)
	planH := planHandler.NewPlanHandler(planSvc)
	planWidgetSvc := planWidgetService.NewPlanWidgetService(planRepo)
	planWidgetH := planWidgetHandler.NewPlanWidgetHandler(planWidgetSvc)

	exportSvc := exportService.NewExportService(planRepo)
	exportH := exportHandler.NewExportHandler(exportSvc)

	plans := r.Group("/plans", authMW.RequireAuth())
	{
		plans.GET("", planH.ListPlans)
		plans.POST("", planH.CreatePlan)
		plans.GET("/:id/widgets", planWidgetH.GetWidgets)
		plans.PUT("/:id/widgets", planWidgetH.SaveWidgets)

		plans.GET("/:id/export/preview", exportH.PreviewLessonPlan)
		plans.GET("/:id/export/docx", exportH.ExportLessonPlanDOCX)
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
