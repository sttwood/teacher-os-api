package dto

type PlanWidgetResponse struct {
	ID          string                 `json:"id"`
	PlanID      string                 `json:"planId"`
	WidgetType  string                 `json:"widgetType"`
	OrderIndex  int                    `json:"orderIndex"`
	Title       string                 `json:"title"`
	ContentJSON map[string]interface{} `json:"contentJSON"`
	IsCollapsed bool                   `json:"isCollapsed"`
	CreatedAt   string                 `json:"createdAt"`
	UpdatedAt   string                 `json:"updatedAt"`
}

type UpdatePlanWidgetsRequest struct {
	Widgets []UpdatePlanWidgetItem `json:"widgets" binding:"required"`
}

type UpdatePlanWidgetItem struct {
	ID          string                 `json:"id"`
	WidgetType  string                 `json:"widgetType" binding:"required"`
	OrderIndex  int                    `json:"orderIndex"`
	Title       string                 `json:"title"`
	ContentJSON map[string]interface{} `json:"contentJSON"`
	IsCollapsed bool                   `json:"isCollapsed"`
}
