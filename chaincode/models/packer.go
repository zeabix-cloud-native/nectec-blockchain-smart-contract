package models

type TransactionPacker struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UserId    string    `json:"userId"`
	Owner     string    `json:"owner"`
	OrgName   string    `json:"orgName"`
	PackingHouseName           string    `json:"packingHouseName"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	IsCanExport  bool  `json:"isCanExport"`
	IsCanDelete  bool  `json:"isCanDelete"`
	PackerGmp *PackerGmp `json:"packerGmp"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
	DocType	  DocType	`json:"docType"`
}

type FilterGetAllPacker struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Search                     *string `json:"search"`
	PackerGmp string `json:"packerGmp"`
	PackingHouseName           string    `json:"packingHouseName"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
}

type PackerGmp struct {
	Id                         string    `json:"id"`
	PackerId 				   string    `json:"packerId"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	IsCanExport   			   bool  `json:"isCanExport"`
	PackingHouseName           string    `json:"packingHouseName"`
	UpdatedDate                string    `json:"updatedDate"`
	Source                     string    `json:"source"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  string `json:"updatedAt"`
	CreatedAt                  string `json:"createdAt"`
}

type PackerTransactionResponse struct {
	Id        string    `json:"id"`
	CertId    string    `json:"certId"`
	UserId    string    `json:"userId"`
	IsCanExport   			   bool  `json:"isCanExport"`
	IsCanDelete   			   bool  `json:"isCanDelete"`
	PackerGmp PackerGmp `json:"packerGmp"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
	PackingHouseName           string    `json:"packingHouseName"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
}

type PackerGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*PackerTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

type PackerByIdResponse struct {
	Data string              `json:"data"`
	Obj  *PackerTransactionResponse `json:"obj"`
}
