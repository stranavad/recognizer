package types

type CreateExamDto struct {
	Name string `json:"name" binding:"required"`
}
