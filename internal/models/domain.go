package models

import "time"

type Domain struct {
	ID         int       `json:"id"`
	DomainName string    `json:"domain_name"`
	Registrar  string    `json:"registrar"`
	ExpiryDate time.Time `json:"expiry_date"`
	Status     string    `json:"status"`
}
