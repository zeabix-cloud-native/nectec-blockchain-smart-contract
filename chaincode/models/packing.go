package models

import "time"

type TransactionPacking struct {
	Id             string    `json:"id"`
	OrderID        string    `json:"orderId"`
	FarmerID       string    `json:"farmerId"`
	ForecastWeight float32   `json:"forecastWeight"`
	ActualWeight   float32   `json:"actualWeight"`
	SavedTime      string    `json:"savedTime"`
	ApprovedDate   string    `json:"approvedDate"`
	ApprovedType   string    `json:"approvedType"`
	FinalWeight    float32   `json:"finalWeight"`
	Remark         string    `json:"remark"`
	CancelReason   string    `json:"cancelReason"`
	PackerId       string    `json:"packerId"`
	Gmp            string    `json:"gmp"`
	PackingHouseName            string    `json:"packingHouseName"`
	Gap            string    `json:"gap"` // รหัสซื้อขาย
	ProcessStatus  int       `json:"processStatus"`
	SellingStep    int    `json:"sellingStep"`
	Owner          string    `json:"owner"`
	OrgName        string    `json:"orgName"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	DocType        DocType   `json:"docType"`
}

type FilterGetAllPacking struct {
	Skip               int      `json:"skip"`
	Limit              int      `json:"limit"`
	Search             *string  `json:"search"`
	PackerId           *string  `json:"packerId"`
	FarmerID					 *string  `json:"farmerId"` 
	CertID					 	 *string  `json:"certId"` 
	Gap                *string  `json:"gap"`
	StartDate          *string  `json:"startDate"`
	EndDate            *string  `json:"endDate"`
	PackingHouseName            string    `json:"packingHouseName"`
	ForecastWeightFrom *float32 `json:"forecastWeightFrom"`
	ForecastWeightTo   *float32 `json:"forecastWeightTo"`
	ProcessStatus      *int     `json:"processStatus"`
}

type PackingTransactionResponse struct {
	Id             string  `json:"id"`
	OrderID        string  `json:"orderId"`
	FarmerID       string  `json:"farmerId"`
	ForecastWeight float32 `json:"forecastWeight"`
	ActualWeight   float32 `json:"actualWeight"`
	// IsPackerSaved  bool      `json:"isPackerSaved"`
	SavedTime string `json:"savedTime"`
	// IsApproved     bool      `json:"isApproved"`
	ApprovedDate  string    `json:"approvedDate"`
	ApprovedType  string    `json:"approvedType"`
	FinalWeight   float32   `json:"finalWeight"`
	Remark        string    `json:"remark"`
	PackerId      string    `json:"packerId"`
	Gmp           						 string    `json:"gmp"`
	PackingHouseName           string    `json:"packingHouseName"`
	Gap           string    `json:"gap"`
	CancelReason           string    `json:"cancelReason"`
	ProcessStatus int       `json:"processStatus"`
	TotalSold        float32   `json:"totalSold"`
	SellingStep       int    `json:"sellingStep"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

type PackingGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*PackingTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

type PackingTransactionHistory struct {
	TxId      string                `json:"tx_id"`
	IsDelete  bool                  `json:"isDelete"`
	Value     []*PackingTransactionResponse `json:"value"`
	Timestamp string                `json:"timestamp"`
}