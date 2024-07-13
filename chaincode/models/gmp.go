package models

type TransactionGmp struct {
	Id                         string    `json:"id"`
	PackerId 									 string    `json:"packerId"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	PackingHouseName           string    `json:"packingHouseName"`
	UpdatedDate                string    `json:"updatedDate"`
	Source                     string    `json:"source"`
	DocType                    DocType   `json:"docType"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  string `json:"updatedAt"`
	CreatedAt                  string `json:"createdAt"`
	IsCanDelete 			   bool       `json:"isCanDelete"`
}

type FilterGetAllGmp struct {
	Skip                       int     `json:"skip"`
	PackerId 									 string  `json:"packerId"`
	Limit                      int     `json:"limit"`
	PackingHouseRegisterNumber *string `json:"packingHouseRegisterNumber"`
	PackingHouseName 		   *string `json:"packingHouseName"`
	Address                    *string `json:"address"`
	Search                     *string `json:"search"`
	AvailableGmp               *string `json:"availableGmp"`
}

type GmpTransactionResponse struct {
	Id                         string    `json:"id"`
	PackerId 									 string    `json:"packerId"`
	PackingHouseRegisterNumber string    `json:"packingHouseRegisterNumber"`
	Address                    string    `json:"address"`
	PackingHouseName           string    `json:"packingHouseName"`
	UpdatedDate                string    `json:"updatedDate"`
	IsCanDelete 			   bool       `json:"isCanDelete"`
	Source                     string    `json:"source"`
	UpdatedAt                  string `json:"updatedAt"`
	CreatedAt                  string `json:"createdAt"`
}

type GmpGetAllResponse struct {
	Data  string                `json:"data"`
	Obj   []*GmpTransactionResponse `json:"obj"`
	Total int                   `json:"total"`
}

type GetByRegisterNumberResponse struct {
	Data string              `json:"data"`
	Obj  *GmpTransactionResponse `json:"obj"`
}
