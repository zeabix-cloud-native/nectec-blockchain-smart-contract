package models

type PlantTypeModel struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Province    string      `json:"province"`
	District    string      `json:"district"`
	PostCode    string      `json:"postCode"`
	Email    	string      `json:"email"`
	IssueDate   string      `json:"issueDate"`
	ExpiredDate string      `json:"expiredDate"`
	PlantType   string      `json:"plantType"`
	ExporterId  string      `json:"exporterId"`
    Owner       string      `json:"owner"`
	OrgName     string      `json:"orgName"`
	DocType     DocType     `json:"docType"`
	CreatedAt   string      `json:"createdAt"`
    UpdatedAt   string      `json:"updatedAt"`
}