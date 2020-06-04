package model

type Response struct {
	TotalCount int `json:"total_count"`
	Items []Repo `json:"items"`
	Error error
}
