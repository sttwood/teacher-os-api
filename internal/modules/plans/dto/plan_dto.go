package dto

type CreatePlanRequest struct {
	SubjectID      string `json:"subjectId" binding:"required"`
	LearningUnitID string `json:"learningUnitId" binding:"required"`
	LessonNo       int    `json:"lessonNo"`
	LessonTitle    string `json:"lessonTitle" binding:"required,min=1"`
	LessonHours    int    `json:"lessonHours"`

	Title        string `json:"title" binding:"required,min=1"`
	SubjectGroup string `json:"subjectGroup"`
	GradeLevel   string `json:"gradeLevel"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`
}

type PlanResponse struct {
	ID             string `json:"id"`
	OwnerID        string `json:"ownerId"`
	SubjectID      string `json:"subjectId,omitempty"`
	LearningUnitID string `json:"learningUnitId,omitempty"`
	LessonNo       int    `json:"lessonNo"`
	LessonTitle    string `json:"lessonTitle"`
	LessonHours    int    `json:"lessonHours"`

	Title        string `json:"title"`
	SubjectGroup string `json:"subjectGroup"`
	GradeLevel   string `json:"gradeLevel"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

type ListPlansQuery struct {
	Status       string `form:"status"`
	SubjectGroup string `form:"subjectGroup"`
	GradeLevel   string `form:"gradeLevel"`
}

type ListPlansMeta struct {
	Total int `json:"total"`
}
