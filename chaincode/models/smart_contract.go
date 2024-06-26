package models

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}