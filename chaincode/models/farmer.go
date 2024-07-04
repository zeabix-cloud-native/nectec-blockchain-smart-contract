package models

import "time"

type TransactionFarmer struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	ProfileImg    string    `json:"profileImg"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	FarmerGaps []FarmerGap `json:"farmerGaps"`
	DocType DocType `json:"docType"`
}

type FilterGetAllFarmer struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	FarmerGap string `json:"farmerGap"`
}

type FarmerGap struct {
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
	IsCanDelete bool       `json:"isCanDelete"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FarmerTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	ProfileImg    string    `json:"profileImg"`
	FarmerGaps []FarmerGap `json:"farmerGaps"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type FarmerGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*FarmerTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

type FarmerTransactionHistory struct {
	TxId      string                `json:"tx_id"`
	IsDelete  bool                  `json:"isDelete"`
	Value     []*FarmerTransactionResponse `json:"value"`
	Timestamp string                `json:"timestamp"`
}
