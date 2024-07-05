package models

import "time"

type TransactionExporter struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	PlantType     string    `json:"plantType"`
	Name     string    `json:"name"`
	Address     string    `json:"address"`
	District     string    `json:"district"`
	Province     string    `json:"province"`
	PostCode     string    `json:"postCode"`
	Email     string    `json:"email"`
	IssueDate     string    `json:"issueDate"`
	ExpiredDate     string    `json:"expiredDate"`
	OrgName   string    `json:"orgName"`
	DocType   DocType   `json:"docType"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExporterFilterGetAll struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Search             *string  `json:"search"`
	Province             *string  `json:"province"`
	District             *string  `json:"district"`
	IssueDateFrom     *string    `json:"issueDateFrom"`
	IssueDateTo     *string    `json:"issueDateTo"`
	ExpireDateFrom     *string    `json:"ExpireDateFrom"`
	ExpireDateTo     *string    `json:"ExpireDateTo"`
}

type ExporterTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	PlantType     string    `json:"plantType"`
	PostCode     string    `json:"postCode"`
	Name     string    `json:"name"`
	Address     string    `json:"address"`
	District     string    `json:"district"`
	Province     string    `json:"province"`
	Email     string    `json:"email"`
	IssueDate     string    `json:"issueDate"`
	ExpiredDate     string    `json:"expiredDate"`
	OrgName   string    `json:"orgName"`
	DocType   DocType   `json:"docType"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type ExporterGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*ExporterTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}
