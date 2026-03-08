package dto

type CreateSubjectRequest struct {
	SubjectGroup string `json:"subjectGroup" binding:"required,min=1"`
	SubjectName  string `json:"subjectName" binding:"required,min=1"`
	SubjectCode  string `json:"subjectCode"`
	GradeLevel   string `json:"gradeLevel" binding:"required,min=1"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`
}

type SubjectResponse struct {
	ID           string `json:"id"`
	OwnerID      string `json:"ownerId"`
	SubjectGroup string `json:"subjectGroup"`
	SubjectName  string `json:"subjectName"`
	SubjectCode  string `json:"subjectCode"`
	GradeLevel   string `json:"gradeLevel"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`
	Status       string `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
