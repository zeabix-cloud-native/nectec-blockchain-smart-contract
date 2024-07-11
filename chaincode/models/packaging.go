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