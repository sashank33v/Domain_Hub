package handlers

import (
	"database/sql"

	"domainhub/internal/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type testDomainService struct {
	createErr  error
	getAll     []models.Domain
	total      int
	getByID    *models.Domain
	getbyIDErr error
	updateErr  error
	deleteErr  error
}

func (t *testDomainService) Create(domains *models.Domain) error {
	return t.createErr
}
func (t *testDomainService) GetAll(limit int, offset int, status string, registrar string) ([]models.Domain, int, error) {
	return t.getAll, t.total, nil
}

func (t *testDomainService) GetByID(id int) (*models.Domain, error) {
	return t.getByID, t.getbyIDErr
}

func (t *testDomainService) Update(id int, domain *models.Domain) error {
	return t.updateErr
}

func (t *testDomainService) Delete(id int) error {
	return t.deleteErr
}

func TestCreateHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		testservice := &testDomainService{}
		handler := NewDomainHandler(testservice)
		body := `{
		"domain_name":"test.com",
		"registrar":"GoDaddy",
		"expiry_date":"2030-01-01T00:00:00Z",
		"status":"ACTIVE"
		}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/domains",
			strings.NewReader(body),
		)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusCreated,
				rec.Code,
			)
		}

	})
	t.Run("invalid JSON", func(t *testing.T) {
		handler := NewDomainHandler(&testDomainService{})

		req := httptest.NewRequest(
			http.MethodPost,
			"/domains",
			strings.NewReader(`invalid Json`),
		)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("Invalid Domain", func(t *testing.T) {
		testService := &testDomainService{
			createErr: errors.New("Invalid Domain"),
		}

		handler := NewDomainHandler(testService)

		body := `{
		"domain_name":"test.com",
		"registrar":"GoDaddy",
		"expiry_date":"2030-01-01T00:00:00Z",
		"status":"ACTIVE"
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/domains",
			strings.NewReader(body),
		)

		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusInternalServerError,
				rec.Code,
			)
		}
	})

}
func TestGetAllHandler(t *testing.T) {
	testservice := testDomainService{
		getAll: []models.Domain{
			{
				ID:         1,
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "ACTIVE",
			},
		},
		total: 1,
	}
	handler := NewDomainHandler(&testservice)

	req := httptest.NewRequest(
		http.MethodGet,
		"/domains?page=1&limit=10",
		nil,
	)
	rec := httptest.NewRecorder()
	handler.GetAll(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}

func TestGetByIDHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testService := &testDomainService{
			getByID: &models.Domain{
				ID:         1,
				DomainName: "example.com",
				Registrar:  "GoDaddy",
				ExpiryDate: time.Now(),
				Status:     "ACTIVE",
			},
		}
		handler := NewDomainHandler(testService)

		req := httptest.NewRequest(
			http.MethodGet,
			"/domains/1",
			nil,
		)

		r := chi.NewRouter()
		r.Get("/domains/{id}", handler.GetByID)

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}

	})
	t.Run("Invalid Id", func(t *testing.T) {
		handler := NewDomainHandler(&testDomainService{})

		req := httptest.NewRequest(
			http.MethodGet,
			"/domains/abc",
			nil,
		)

		r := chi.NewRouter()
		r.Get("/domains/{id}", handler.GetByID)

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusBadRequest,
				rec.Code,
			)
		}
	})

	t.Run("not found", func(t *testing.T) {
		testservice := &testDomainService{
			getbyIDErr: sql.ErrNoRows,
		}
		handler := NewDomainHandler(testservice)

		req := httptest.NewRequest(
			http.MethodGet,
			"/domains/999",
			nil,
		)
		r := chi.NewRouter()
		r.Get("/domains/{id}", handler.GetByID)

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf(
				"expected %d, got %d",
				http.StatusNotFound,
				rec.Code,
			)
		}

	},
	)
}

func TestUpdatehandler(t *testing.T) {
	testservice := &testDomainService{}

	handler := NewDomainHandler(testservice)

	body := `{
"domain_name":"updated.com",
"registrar":"Cloudfare",
"expiry_date":"2031-01-01T00:00:00Z",
"status":"ACTIVE"
}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/domains/1",
		strings.NewReader(body),
	)
	r := chi.NewRouter()
	r.Put("/domains/{id}", handler.Update)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
}
func TestDeleteHandler(t *testing.T) {

	testService := &testDomainService{}

	handler := NewDomainHandler(testService)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/domains/1",
		nil,
	)

	r := chi.NewRouter()
	r.Delete("/domains/{id}", handler.Delete)

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

}
