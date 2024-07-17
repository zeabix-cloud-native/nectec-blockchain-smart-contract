package models

type TransactionExporter struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	Owner     string    `json:"owner"`
	PlantType     string    `json:"plantType"`
	OrgName   string    `json:"orgName"`
	PlantTypeDetail PlantTypeModel `json:"plantTypeDetail"`
	DocType   DocType   `json:"docType"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type ExporterFilterGetAll struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Search             *string  `json:"search"`
	Province             *string  `json:"province"`
	District             *string  `json:"district"`
	CreatedAtFrom  *string `json:"createdAtFrom"`
	CreatedAtTo *string `json:"createdAtTo"`
	ExpireDateFrom     *string    `json:"ExpireDateFrom"`
	ExpireDateTo     *string    `json:"ExpireDateTo"`
}

type ExporterTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	PlantType     string    `json:"plantType"`
	PostCode     string    `json:"postCode"`
	Address     string    `json:"address"`
	District     string    `json:"district"`
	Province     string    `json:"province"`
	Email     string    `json:"email"`
	IssueDate     string    `json:"issueDate"`
	ExpiredDate     string    `json:"expiredDate"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
	PlantTypeDetail PlantTypeModel `json:"plantTypeDetail"`
}

type ExporterGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*ExporterTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}