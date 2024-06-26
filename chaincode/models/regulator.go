package models

import "time"

type TransactionRegulator struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type FilterGetAllRegulator struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type RegulatorTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegulatorGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*RegulatorTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}
