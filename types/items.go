package types

type CreateItem struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	ExamId  uint   `json:"examId"`
	GroupId uint   `json:"groupId"`
}

type UpdateItem struct {
	Name    string `json:"name"`
	Image   string `json:"image"`
	GroupId uint   `json:"groupId"`
}
