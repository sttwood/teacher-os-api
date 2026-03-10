package dto

type LessonPlanDocument struct {
	PlanID         string `json:"planId"`
	Title          string `json:"title"`
	SubjectID      string `json:"subjectId"`
	LearningUnitID string `json:"learningUnitId"`
	LessonNo       int    `json:"lessonNo"`
	LessonTitle    string `json:"lessonTitle"`
	LessonHours    int    `json:"lessonHours"`

	SubjectGroup string `json:"subjectGroup"`
	GradeLevel   string `json:"gradeLevel"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academicYear"`
	SchoolName   string `json:"schoolName"`
	TeacherName  string `json:"teacherName"`

	BasicInfo         map[string]interface{} `json:"basicInfo"`
	Objective         map[string]interface{} `json:"objective"`
	StandardIndicator map[string]interface{} `json:"standardIndicator"`
	Activity          map[string]interface{} `json:"activity"`
	MediaMaterial     map[string]interface{} `json:"mediaMaterial"`
	Assessment        map[string]interface{} `json:"assessment"`
	Homework          map[string]interface{} `json:"homework"`
}
