package dto

type CreatePlanRequest struct {
	Title        string `json:"title" binding:"required,min=1"`
	SubjectGroup string `json:"subjectGroup"`
	GradeLevel   string `json:"gradeLevel"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`
}

type PlanResponse struct {
	ID           string `json:"id"`
	OwnerID      string `json:"ownerId"`
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
