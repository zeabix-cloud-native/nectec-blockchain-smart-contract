package models

import "time"

type TransactionNectecStaff struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	ProfileImg    string    `json:"profileImg"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	DocType DocType `json:"docType"`
}

type FilterGetAllNectecStaff struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type NectecStaffTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	ProfileImg    string    `json:"profileImg"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type GetAllNectecStaffResponse struct {
	Data  string                `json:"data"`
	Obj   []*NectecStaffTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

