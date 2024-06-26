package models

import "time"

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
	UpdatedDate string    `json:"updatedDate"`
	Source      string    `json:"source"`
	FarmerID    string    `json:"farmerId"`
	Owner       string    `json:"owner"`
	OrgName     string    `json:"orgName"`
	DocType     DocType   `json:"docType"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
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
	IssueDate    *string  `json:"issueDate"`
	ExpireDate   *string  `json:"expireDate"`
	AvailableGap *string  `json:"availableGap"`
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
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
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