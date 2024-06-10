package types

type CreateGroupDto struct {
	Name   string `json:"name"`
	ExamID uint   `json:"examId"`
}

type UpdateGroupDto struct {
	Name string `json:"name"`
}
