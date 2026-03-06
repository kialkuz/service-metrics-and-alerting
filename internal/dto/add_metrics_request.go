package dto

type AddMetricsRequest struct {
	Type  string `json:"type" binding:"required"`
	Name  string `json:"name" binding:"required"`
	Value string `json:"value" binding:"required"`
}
