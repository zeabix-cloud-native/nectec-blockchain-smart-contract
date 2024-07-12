package models

type TransactionPackaging struct {
	Id                         string    `json:"id"`
	ContainerId 			   string    `json:"containerId"`
	ExportId 				   string    `json:"exportId"`
	LotNumber 				   string    `json:"lotNumber"`
	BoxId 				       string    `json:"boxId"`
	Gap 				       string    `json:"gap"`
	Gmp 				       string    `json:"gmp"`
	GradeName 				   string    `json:"gradeName"`
	Gtin14 				       string    `json:"gtin14"`
	Gtin13 				       string    `json:"gtin13"`
	VarietyName 			   string    `json:"varietyName"`
	DocType                    DocType   `json:"docType"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  string `json:"updatedAt"`
	CreatedAt                  string `json:"createdAt"`
}

type ReadPackaging struct {
	Id                         string    `json:"id"`
	ContainerId 			   string    `json:"containerId"`
	ExportId 				   string    `json:"exportId"`
	LotNumber 				   string    `json:"lotNumber"`
	BoxId 				       string    `json:"boxId"`
	Gap 				       string    `json:"gap"`
	Gmp 				       string    `json:"gmp"`
	GradeName 				   string    `json:"gradeName"`
	Gtin14 				       string    `json:"gtin14"`
	Gtin13 				       string    `json:"gtin13"`
	VarietyName 			   string    `json:"varietyName"`
	DocType                    DocType   `json:"docType"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  string `json:"updatedAt"`
	CreatedAt                  string `json:"createdAt"`
	GmpDetail                  TransactionGmp `json:"gmpDetail"` 
}

type PackagingFilterParams struct {
	ContainerId 			   string    `json:"containerId"`
	BoxId 				       string    `json:"boxId"`
	Gtin13 				       string    `json:"gtin13"`
	Gap 				       string    `json:"gap"`
	StartDate   			   string 	 `json:"startDate"`
	EndDate                    string 	 `json:"endDate"`
	Skip                       int    	 `json:"skip"`
	Limit                      int    	 `json:"limit"`
}

type TransactionPackagingResponse struct {
	Data   []*TransactionPackaging `json:"obj"`
	Total int                   `json:"total"`
}