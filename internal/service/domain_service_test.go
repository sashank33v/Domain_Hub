package service

import (
	"domainhub/internal/models"
	"testing"
	"time"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  models.Domain
		wantErr bool
	}{
		{
			name: "valid active domain",
			domain: models.Domain{
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "ACTIVE",
			},
			wantErr: false,
		}, {
			name: "valid expired domain",
			domain: models.Domain{
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "EXPIRED",
			},
			wantErr: false,
		}, {
			name: "missing domain name",
			domain: models.Domain{
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "missing registrar",
			domain: models.Domain{
				DomainName: "example.com",
				ExpiryDate: time.Now(),
				Status:     "ACTIVE",
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			domain: models.Domain{
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "INVALID",
			},
			wantErr: true,
		},
		{
			name: "missing expiry date",
			domain: models.Domain{
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				Status:     "ACTIVE",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := validateDomain(&tt.domain)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateDomain() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}
