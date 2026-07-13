package domain

type Company struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Opportunity struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"companyId"`
}
