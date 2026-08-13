package dto

import "domainhub/internal/models"

func ToDomainResponse(domain *models.Domain) DomainResponse {
	return DomainResponse{
		ID:         domain.ID,
		DomainName: domain.DomainName,
		Registrar:  domain.Registrar,
		ExpiryDate: domain.ExpiryDate.Format("2006-01-02T15:04:05Z07:00"),
		Status:     domain.Status,
	}
}
