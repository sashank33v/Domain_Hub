package repository

import (
	"domainhub/internal/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	expirydate := time.Date(
		2030,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"domain_name",
			"registrar",
			"expiry_date",
			"status",
		},
	).AddRow(4,
		"example.com",
		"GoDaddy",
		expirydate,
		"ACTIVE",
	)
	mock.ExpectQuery(
		regexp.QuoteMeta(`
	Select id, domain_name, registrar, expiry_date, status
	From domains
	Where id = $1;
		`),
	).
		WithArgs(4).
		WillReturnRows(rows)

	domain, err := repo.GetByID(4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if domain.ID != 4 {
		t.Errorf("expected ID 4, got %d", domain.ID)
	}

	if domain.DomainName != "example.com" {
		t.Errorf(
			"expected domain example.com, got %s",
			domain.DomainName,
		)
	}

	if domain.Status != "ACTIVE" {
		t.Errorf(
			"expected status ACTIVE, got %s",
			domain.Status,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}

}

func TestCreate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	domain := &models.Domain{
		DomainName: "test.com",
		Registrar:  "GoDaddy",
		ExpiryDate: time.Date(
			2030,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		),
		Status: "ACTIVE",
	}

	mock.ExpectQuery(
		regexp.QuoteMeta(`
	INSERT INTO domains (domain_name, registrar, expiry_date, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`),
	).
		WithArgs(
			domain.DomainName,
			domain.Registrar,
			domain.ExpiryDate,
			domain.Status,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(10),
		)

	err = repo.Create(domain)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if domain.ID != 10 {
		t.Errorf("expected ID 10, got %d", domain.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}
}
func TestUpdate(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	domain := &models.Domain{
		DomainName: "updated.com",
		Registrar:  "Cloudflare",
		ExpiryDate: time.Date(
			2031,
			1,
			1,
			0,
			0,
			0,
			0,
			time.UTC,
		),
		Status: "ACTIVE",
	}

	mock.ExpectExec(
		regexp.QuoteMeta(`
	Update domains
	Set
	domain_name= $1,
	registrar= $2,
	expiry_date = $3,
	status = $4
	where id = $5;
	`),
	).
		WithArgs(
			domain.DomainName,
			domain.Registrar,
			domain.ExpiryDate,
			domain.Status,
			4,
		).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	err = repo.Update(4, domain)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}
}
func TestDelete(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	mock.ExpectExec(
		regexp.QuoteMeta(`
	Delete from domains
	where id = $1;
	`),
	).
		WithArgs(4).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)

	err = repo.Delete(4)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}
}
func TestGetAll(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := NewDomainRepository(db)

	expiryDate := time.Date(
		2030,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
		Select Count(*)
		from domains
		WHERE status = $1 AND registrar ILIKE '%' || $2 || '%';
		`),
	).
		WithArgs("ACTIVE", "GoDaddy").
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(5),
		)

	rows := sqlmock.NewRows(
		[]string{
			"id",
			"domain_name",
			"registrar",
			"expiry_date",
			"status",
		},
	).AddRow(
		4,
		"example.com",
		"GoDaddy",
		expiryDate,
		"ACTIVE",
	).AddRow(
		7,
		"test.com",
		"GoDaddy",
		expiryDate,
		"ACTIVE",
	)

	mock.ExpectQuery(
		regexp.QuoteMeta(`
		select id, domain_name, registrar, expiry_date, status
		From domains
		where status = $1 and registrar ILIKE '%' || $2 || '%'
		order by id
		Limit $3 offset $4;
		`),
	).
		WithArgs("ACTIVE", "GoDaddy", 2, 2).
		WillReturnRows(rows)

	domains, total, err := repo.GetAll(
		2,
		2,
		"ACTIVE",
		"GoDaddy",
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}

	if len(domains) != 2 {
		t.Errorf("expected 2 domains, got %d", len(domains))
	}

	if domains[0].ID != 4 {
		t.Errorf("expected first domain ID 4, got %d", domains[0].ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("mock expectations were not met: %v", err)
	}
}
