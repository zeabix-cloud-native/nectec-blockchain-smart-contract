package models

type TransactionHscode struct {
	Id                         string    `json:"id"`
	Hscode 					   string    `json:"hscode"`
	Description 			   string    `json:"description"`
	DocType                    DocType   `json:"docType"`
	Owner                      string    `json:"owner"`
	OrgName                    string    `json:"orgName"`
	UpdatedAt                  string    `json:"updatedAt"`
	CreatedAt                  string    `json:"createdAt"`
}

type HscodeTransactionResponse struct {
	Id                         string    `json:"id"`
	Hscode 					   string    `json:"hscode"`
	Description 			   string    `json:"description"`
	UpdatedAt                  string    `json:"updatedAt"`
	CreatedAt                  string    `json:"createdAt"`
}

type HscodeFilterParams struct {
	Skip                       int    	 `json:"skip"`
	Limit                      int    	 `json:"limit"`
}

type TransactionHscodeResponse struct {
	Data   []*TransactionHscode `json:"obj"`
	Total int                   `json:"total"`
}
