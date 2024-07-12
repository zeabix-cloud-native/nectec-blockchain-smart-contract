package models

type TransactionGap struct {
	Id          string    `json:"id"`
	CertID      string    `json:"certId"`
	DisplayCertID      string    `json:"displayCertId"`
	AreaCode    string    `json:"areaCode"`
	AreaRai     float32   `json:"areaRai"`
	AreaStatus  string    `json:"areaStatus"`
	OldAreaCode string    `json:"oldAreaCode"`
	IssueDate   string    `json:"issueDate"`
	ExpireDate  string    `json:"expireDate"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	Source      string    `json:"source"`
	FarmerID    string    `json:"farmerId"`
	IsCanDelete bool       `json:"isCanDelete"`
	Owner       string    `json:"owner"`
	OrgName     string    `json:"orgName"`
	DocType     DocType   `json:"docType"`
	UpdatedAt   string `json:"updatedAt"`
	CreatedAt   string `json:"createdAt"`
}

type FilterGetAllGap struct {
	Skip         int      `json:"skip"`
	Limit        int      `json:"limit"`
	CertID       *string  `json:"certId"`
	FarmerID       *string  `json:"farmerId"`
	AreaCode     *string  `json:"areaCode"`
	District     *string  `json:"district"`
	Province     *string  `json:"province"`
	AreaRaiFrom  *float32 `json:"areaRaiFrom"`
	AreaRaiTo    *float32 `json:"areaRaiTo"`
	AvailableGap *string  `json:"availableGap"`
	CreatedAtFrom  *string `json:"createdAtFrom"`
	CreatedAtTo *string `json:"createdAtTo"`
	IssueDateTo    *string `json:"issueDateTo"`
	ExpireDateFrom  *string `json:"expireDateFrom"`
	ExpireDateTo    *string `json:"expireDateTo"`
	Gaps           []string `json:"gaps"`
}

type GapTransactionResponse struct {
	Id          			 string    `json:"id"`
	DisplayCertID      string    `json:"displayCertId"`
	CertID      string    `json:"certId"`
	AreaCode    string    `json:"areaCode"`
	AreaRai     float32   `json:"areaRai"`
	AreaStatus  string    `json:"areaStatus"`
	OldAreaCode string    `json:"oldAreaCode"`
	IssueDate   string    `json:"issueDate"`
	ExpireDate  string    `json:"expireDate"`
	District    string    `json:"district"`
	Province    string    `json:"province"`
	UpdatedDate string    `json:"updatedDate"`
	Source      string    `json:"source"`
	FarmerID    string    `json:"farmerId"`
	UpdatedAt   string `json:"updatedAt"`
	CreatedAt   string `json:"createdAt"`
	TotalSold   float32   `json:"totalSold"`
	IsCanDelete bool       `json:"isCanDelete"`
}

type GetAllGapResponse struct {
	Data  string                `json:"data"`
	Obj   []*GapTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

type GetGapByCertIdResponse struct {
	Data string              `json:"data"`
	Obj  *GapTransactionResponse `json:"obj"`
}