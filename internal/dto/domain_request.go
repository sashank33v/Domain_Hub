package dto

type CreateDomainRequest struct {
	DomainName string `json:"domain_name"`
	Registrar  string `json:"registrar"`
	ExpiryDate string `json:"expiry_date"`
	Status     string `json:"status"`
}

type UpdateDomainRequest struct {
	DomainName string `json:"domain_name"`
	Registrar  string `json:"registrar"`
	ExpiryDate string `json:"expiry_date"`
	Status     string `json:"status"`
}
