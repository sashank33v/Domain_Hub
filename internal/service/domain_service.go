package service

import (
	"context"
	"domainhub/internal/models"
	"domainhub/internal/repository"
	"errors"
	"fmt"
	"time"
)

var ErrInvalidDomain = errors.New("Invalid domain")

type DomainService struct {
	repo     *repository.DomainRepository
	producer DomainEventPublisher
}
type DomainEventPublisher interface {
	PublishDomainCreated(ctx context.Context, domain *models.Domain) error
}

func NewDomainService(repo *repository.DomainRepository, producer DomainEventPublisher) *DomainService {
	return &DomainService{
		repo:     repo,
		producer: producer,
	}

}

func (s *DomainService) Create(domain *models.Domain) error {

	if err := validateDomain(domain); err != nil {
		return err
	}
	if err := s.repo.Create(domain); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := s.producer.PublishDomainCreated(ctx, domain); err != nil {
		return fmt.Errorf("domain created, but event publish failed: %w", err)
	}
	return nil

}
func (s *DomainService) GetAll(limit int, offset int, status string, registrar string) ([]models.Domain, int, error) {
	return s.repo.GetAll(limit, offset, status, registrar)
}
func (s *DomainService) GetByID(id int) (*models.Domain, error) {
	return s.repo.GetByID(id)
}
func (s *DomainService) Update(id int, domain *models.Domain) error {
	if err := validateDomain(domain); err != nil {
		return err
	}

	return s.repo.Update(id, domain)
}
func (s *DomainService) Delete(id int) error {
	return s.repo.Delete(id)
}

func validateDomain(domain *models.Domain) error {
	if domain.DomainName == "" {
		return fmt.Errorf("%w: domain name is required", ErrInvalidDomain)
	}

	if domain.Registrar == "" {
		return fmt.Errorf("%w: registrar is required", ErrInvalidDomain)
	}

	if domain.Status != "ACTIVE" && domain.Status != "EXPIRED" {
		return fmt.Errorf("%w: status must be active or expired", ErrInvalidDomain)
	}
	if domain.ExpiryDate.IsZero() {
		return fmt.Errorf("%w: expiry date is required", ErrInvalidDomain)
	}
	return nil

}
