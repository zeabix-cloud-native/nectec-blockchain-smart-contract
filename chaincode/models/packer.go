package models

import "time"

type TransactionPacker struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UserId    string    `json:"userId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	PackerGmp PackerGmp `json:"packerGmp"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type FilterGetAllPacker struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	PackerGmp string `json:"packerGmp"`
}

type PackerGmp struct {
	Id                         string    `json:"id"`
	PackerId 				   string    `json:"packerId"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	PackingHouseName           string    `json:"packingHouseName"`
	UpdatedDate                string    `json:"updatedDate"`
	Source                     string    `json:"source"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  time.Time `json:"updatedAt"`
	CreatedAt                  time.Time `json:"createdAt"`
}

type PackerTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UserId    string    `json:"userId"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type PackerGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*PackerTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}
