package models

type Pagination struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type DocType string

// Doc Type
const (
	Farmer DocType = "farmer"
	Gap DocType = "gap"
)