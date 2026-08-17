package handlers

import (
	"database/sql"
	"domainhub/internal/dto"
	"domainhub/internal/models"
	"domainhub/internal/service"
	"domainhub/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type DomainServiceInterface interface {
	Create(*models.Domain) error
	GetAll(int, int, string, string) ([]models.Domain, int, error)
	GetByID(int) (*models.Domain, error)
	Update(int, *models.Domain) error
	Delete(int) error
}

type DomainHandler struct {
	service DomainServiceInterface
}

func NewDomainHandler(service DomainServiceInterface) *DomainHandler {
	return &DomainHandler{
		service: service,
	}

}

// @Summary Create a domain
// @Description Creates a new domain
// @Tags domains
// @Accept json
// @Produce json
// @Param domain body dto.CreateDomainRequest true "Domain"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /domains [post]
func (h *DomainHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDomainRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid request body",
		)
		return
	}
	expiryDate, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid expiry date",
		)
		return
	}

	domain := models.Domain{
		DomainName: req.DomainName,
		Registrar:  req.Registrar,
		ExpiryDate: expiryDate,
		Status:     req.Status,
	}
	err = h.service.Create(&domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDomain) {
			utils.SendErrorResponse(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}

		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to create domain",
		)
		return
	}
	response := dto.ToDomainResponse(&domain)
	utils.SendResponse(
		w,
		http.StatusCreated,
		"Domain created successfully",
		response,
		nil,
	)

}

// GetAll godoc
// @Summary Get all domains
// @Description Returns domains with pagination and optional status/registrar filtering
// @Tags domains
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of records per page" default(10)
// @Param status query string false "Filter by status" Enums(ACTIVE,EXPIRED)
// @Param registrar query string false "Search by registrar"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /domains [get]
func (h *DomainHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	status := r.URL.Query().Get("status")
	registrar := r.URL.Query().Get("registrar")
	pagestr := r.URL.Query().Get("page")
	limitstr := r.URL.Query().Get("limit")

	if pagestr != "" {
		var err error

		page, err = strconv.Atoi(pagestr)
		if err != nil || page < 1 {
			utils.SendErrorResponse(
				w,
				http.StatusBadRequest,
				"Invalid Page",
			)
			return
		}
	}
	if limitstr != "" {
		var err error

		limit, err = strconv.Atoi(limitstr)
		if err != nil || limit < 1 {
			utils.SendErrorResponse(
				w,
				http.StatusBadRequest,
				"Invalid limit",
			)
			return
		}
	}

	if status != "" && status != "ACTIVE" && status != "EXPIRED" {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid status",
		)
		return
	}

	offset := (page - 1) * limit

	domains, total, err := h.service.GetAll(limit, offset, status, registrar)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to fetch domains",
		)
		return
	}
	responses := make([]dto.DomainResponse, 0, len(domains))

	for _, domain := range domains {
		responses = append(responses, dto.ToDomainResponse(&domain))
	}
	totalPages := (total + limit - 1) / limit

	pagination := models.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	utils.SendResponse(
		w,
		http.StatusOK,
		"Domains fetched successfully",
		responses,
		&pagination,
	)
}

// GetByID godoc
// @Summary Get domain by ID
// @Description Returns a single domain using its ID
// @Tags domains
// @Produce json
// @Param id path int true "Domain ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /domains/{id} [get]
func (h *DomainHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}
	domain, err := h.service.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.SendErrorResponse(
				w,
				http.StatusNotFound,
				"Domain not found",
			)
			return
		}
		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}
	response := dto.ToDomainResponse(domain)
	utils.SendResponse(
		w,
		http.StatusOK,
		"Domain fetched successfully",
		response,
		nil,
	)
}

// Update godoc
// @Summary Update a domain
// @Description Updates an existing domain using its ID
// @Tags domains
// @Accept json
// @Produce json
// @Param id path int true "Domain ID"
// @Param domain body dto.UpdateDomainRequest true "Updated domain data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /domains/{id} [put]
func (h *DomainHandler) Update(w http.ResponseWriter, r *http.Request) {
	idstr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}
	var req dto.UpdateDomainRequest

	err = json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid Request Body",
		)
		return
	}

	expiryDate, err := time.Parse(time.RFC3339, req.ExpiryDate)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid expiry date",
		)
		return
	}

	domain := models.Domain{
		DomainName: req.DomainName,
		Registrar:  req.Registrar,
		ExpiryDate: expiryDate,
		Status:     req.Status,
	}
	err = h.service.Update(id, &domain)
	if err != nil {
		if errors.Is(err, service.ErrInvalidDomain) {
			utils.SendErrorResponse(
				w,
				http.StatusBadRequest,
				err.Error(),
			)
			return
		}
		if errors.Is(err, sql.ErrNoRows) {
			utils.SendErrorResponse(
				w,
				http.StatusNotFound,
				"Domain not found",
			)
			return
		}

		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to update domain",
		)
		return
	}
	response := dto.ToDomainResponse(&domain)
	utils.SendResponse(
		w,
		http.StatusOK,
		"Domain updated successfully",
		response,
		nil,
	)
}

// Delete godoc
// @Summary Delete a domain
// @Description Deletes an existing domain using its ID
// @Tags domains
// @Produce json
// @Param id path int true "Domain ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /domains/{id} [delete]
func (h *DomainHandler) Delete(w http.ResponseWriter, r *http.Request) {

	idstr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		utils.SendErrorResponse(
			w,
			http.StatusBadRequest,
			"Invalid ID",
		)
		return
	}

	err = h.service.Delete(id)
	if err != nil {

		if err == sql.ErrNoRows {
			utils.SendErrorResponse(
				w,
				http.StatusBadRequest,
				"Domain not found",
			)
			return
		}
		utils.SendErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to delete domain",
		)
		return
	}
	utils.SendResponse(
		w,
		http.StatusOK,
		"Domain deleted successfully",
		nil,
		nil,
	)
}
