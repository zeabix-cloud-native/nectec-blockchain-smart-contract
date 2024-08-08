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
	IsCanDelete bool       `json:"isCanDelete"`
	CreatedAt   string      `json:"createdAt"`
    UpdatedAt   string      `json:"updatedAt"`
}

type PlanTypeFilterParams struct {
	AvailablePlanType	string    	 `json:"availablePlantType"`
	PlantType			string    	 `json:"plantType"`
	Skip				int    	 	 `json:"skip"`
	Limit   			int    	 	 `json:"limit"`
	Search             *string  `json:"search"`
	Province             *string  `json:"province"`
	District             *string  `json:"district"`
	CreatedAtFrom  *string `json:"createdAtFrom"`
	CreatedAtTo *string `json:"createdAtTo"`
	ExpireDateFrom     *string    `json:"ExpireDateFrom"`
	ExpireDateTo     *string    `json:"ExpireDateTo"`
}

type PlantTypeResponse struct {
	Data   []*PlantTypeModel `json:"obj"`
	Total int                   `json:"total"`
}