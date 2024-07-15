package models

type Shipper struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Province   string `json:"province"`
	District   string `json:"district"`
	SubDistrict string `json:"subDistrict"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type Receiver struct {
	Name       string `json:"name"`
	Address    string `json:"address"`
	Country    string `json:"country"`
	PostalCode string `json:"postalCode"`
}

type ProductAndPackaging struct {
	HsCode           string `json:"hscode"`
	ProductType      string `json:"productType"`
	ProductGrade     string `json:"productGrade"`
	ProductCategory  string `json:"productCategory"`
	PalletNumber     string `json:"palletNumber"`
	PalletIdentifier string `json:"palletIdentifier"`
	PalletSize       string `json:"palletSize"`
	PalletQuantity   string `json:"palletQuantity"`
	ContainerNumber  string `json:"containerNumber"`
}

type Invoice struct {
	InvoiceNumber     string                `json:"invoiceNumber"`
	InvoiceDate       string                `json:"invoiceDate"`
	TotalWeight       string                `json:"totalWeight"`
	ExportNumber      string                `json:"exportNumber"`
	LotNumber         string                `json:"lotNumber"`
	ProductAndPackaging []ProductAndPackaging `json:"productAndPackaging"`
}

type TransactionFormE struct {
	Id                      string  `json:"id"`
	ReferenceNo       		string  `json:"referenceNo"`
	CountryOfIssuance       string  `json:"countryOfIssuance"`
	RequestType             string  `json:"requestType"`
	Status                  string  `json:"status"`
	CarrierName             string  `json:"carrierName"`
	DepartureCity           string  `json:"departureCity"`
	DestinationCity         string  `json:"destinationCity"`
	CountryOfProduction     string  `json:"countryOfProduction"`
	CountryOfImport         string  `json:"countryOfImport"`
	ExportDate              string  `json:"exportDate"`
	ImportSpecialCondition  string  `json:"importSpecialCondition"`
	PreviousReferenceNumber string  `json:"previousReferenceNumber"`
	CreatedAt               string  `json:"createdAt"`
	UpdatedAt               string  `json:"updatedAt"`
	Owner                   string  `json:"owner"`
	OrgName                 string  `json:"orgName"`
	CancelReason            string  `json:"cancelReason"`
	CreatedById             string  `json:"createdById"`
	Shipper                 Shipper `json:"shipper"`
	Receiver                Receiver `json:"receiver"`
	Invoice                 Invoice `json:"invoice"`
    DocType                 DocType  `json:"docType"`  
}

type FormEFilterParams struct {
	CreatedById                string   `json:"createdById"`
	ReferenceNo       		   string    `json:"referenceNo"`
	ExportNumber     		   string    `json:"exportNumber"`
	RequestType             string  `json:"requestType"`
	LotNumber         		   string    `json:"lotNumber"`
	StartDate   			   string 	 `json:"startDate"`
	EndDate                    string 	 `json:"endDate"`
	Status                     string  `json:"status"`
	Skip                       int    	 `json:"skip"`
	Limit                      int    	 `json:"limit"`
}

type TransactionFormEResponse struct {
	Data   []*TransactionFormE `json:"obj"`
	Total int                   `json:"total"`
}

type FormETransactionHistory struct {
	TxId      string                `json:"tx_id"`
	IsDelete  bool                  `json:"isDelete"`
	Value     []*TransactionFormE `json:"value"`
	Timestamp string                `json:"timestamp"`
}