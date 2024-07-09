package models

type TransactionRegulator struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type FilterGetAllRegulator struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type RegulatorTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type RegulatorGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*RegulatorTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}
