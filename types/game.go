package types

type GetResult struct {
	ItemId uint   `json:"itemId"`
	Answer string `json:"answer"`
}

type GameResponse struct {
	ItemId  uint     `json:"itemId"`
	Image   string   `json:"image"`
	Answers []string `json:"answers"`
}
