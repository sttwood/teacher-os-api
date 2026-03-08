package dto

type CreateLearningUnitRequest struct {
	UnitNo      int    `json:"unitNo" binding:"required"`
	UnitTitle   string `json:"unitTitle" binding:"required,min=1"`
	Description string `json:"description"`
	TotalHours  int    `json:"totalHours"`
	OrderIndex  int    `json:"orderIndex"`
}

type LearningUnitResponse struct {
	ID          string `json:"id"`
	SubjectID   string `json:"subjectId"`
	UnitNo      int    `json:"unitNo"`
	UnitTitle   string `json:"unitTitle"`
	Description string `json:"description"`
	TotalHours  int    `json:"totalHours"`
	OrderIndex  int    `json:"orderIndex"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
