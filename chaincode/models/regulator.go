package models

type TransactionRegulator struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	ProfileImg    string    `json:"profileImg"`
	UserId    string    `json:"userId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	DocType   DocType `json:"docType"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type FilterGetAllRegulator struct {
	UserId string `json:"userId"`
	Skip   int 	  `json:"skip"`
	Limit  int     `json:"limit"`
}

type RegulatorTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type RegulatorGetAllResponse struct {
	Obj   []*TransactionRegulator `json:"obj"`
	Total int                   `json:"total"`
}

type RegulatorByIdResponse struct {
	Data string              `json:"data"`
	Obj  *TransactionRegulator `json:"obj"`
}