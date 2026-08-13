package dto

type DomainResponse struct {
	ID         int    `json:"id"`
	DomainName string `json:"domain_name"`
	Registrar  string `json:"registrar"`
	ExpiryDate string `json:"expiry_date"`
	Status     string `json:"status"`
}
