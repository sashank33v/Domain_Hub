package repository

import (
	"database/sql"
	"domainhub/internal/models"
)

type DomainRepository struct {
	db *sql.DB
}

func NewDomainRepository(db *sql.DB) *DomainRepository {
	return &DomainRepository{
		db: db,
	}
}
func (r *DomainRepository) Create(domain *models.Domain) error {

	query := `
	INSERT INTO domains (domain_name, registrar, expiry_date, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id;
	`

	err := r.db.QueryRow(
		query,
		domain.DomainName,
		domain.Registrar,
		domain.ExpiryDate,
		domain.Status,
	).Scan(&domain.ID)

	if err != nil {
		return err
	}
	return nil

}

func (r *DomainRepository) GetAll(limit int, offset int, status string, registrar string) ([]models.Domain, int, error) {
	var rows *sql.Rows
	var count int
	var err error
	if status == "" && registrar == "" {

		countQuery := `
		Select Count(*)
		from domains;
		`
		err = r.db.QueryRow(countQuery).Scan(&count)
		if err != nil {
			return nil, 0, err
		}
		query := `
	Select id, domain_name, registrar, expiry_date, status 
	From domains 
	order by id
	Limit $1 Offset $2;
	`
		rows, err = r.db.Query(query, limit, offset)
	} else if status != "" && registrar == "" {
		countQuery := `
		Select Count(*)
		from domains
		where status = $1;
		`
		err = r.db.QueryRow(countQuery, status).Scan(&count)
		if err != nil {
			return nil, 0, err
		}
		query := `
		Select id, domain_name, registrar, expiry_date, status
		from domains
		where status = $1
		order by id
		limit $2 offset $3;
		`
		rows, err = r.db.Query(query, status, limit, offset)
	} else if status == "" && registrar != "" {
		countQuery := `
		Select Count(*)
		from domains
		WHERE registrar ILIKE '%' || $1 || '%';
		`
		err = r.db.QueryRow(countQuery, registrar).Scan(&count)
		if err != nil {
			return nil, 0, err
		}

		query := `
		Select id, domain_name, registrar, expiry_date, status
		From domains
		where registrar ILIKE '%' || $1 || '%'
		order by id
		Limit $2 offset $3;
		`
		rows, err = r.db.Query(query, registrar, limit, offset)

	} else {

		countQuery := `
		Select Count(*)
		from domains
		WHERE status = $1 AND
		 registrar ILIKE '%' || $2 || '%';
		`
		err = r.db.QueryRow(countQuery, status, registrar).Scan(&count)
		if err != nil {
			return nil, 0, err
		}
		query := `
		select id, domain_name, registrar, expiry_date, status
		From domains
		where status = $1 and registrar ILIKE '%' || $2 || '%'
		order by id
		Limit $3 offset $4;
		`
		rows, err = r.db.Query(query, status, registrar, limit, offset)
	}
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var domains []models.Domain

	for rows.Next() {
		var domain models.Domain

		err := rows.Scan(
			&domain.ID,
			&domain.DomainName,
			&domain.Registrar,
			&domain.ExpiryDate,
			&domain.Status,
		)

		if err != nil {
			return nil, 0, err
		}
		domains = append(domains, domain)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return domains, count, nil
}

func (r *DomainRepository) GetByID(id int) (*models.Domain, error) {
	query := `
	Select id, domain_name, registrar, expiry_date, status
	From domains
	Where id = $1;
	`
	var domain models.Domain
	err := r.db.QueryRow(query, id).Scan(
		&domain.ID,
		&domain.DomainName,
		&domain.Registrar,
		&domain.ExpiryDate,
		&domain.Status,
	)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}
func (r *DomainRepository) Update(id int, domain *models.Domain) error {
	query := `
	Update domains
	Set
	domain_name= $1,
	registrar= $2,
	expiry_date = $3,
	status = $4
	where id = $5;
	`
	result, err := r.db.Exec(
		query,
		domain.DomainName,
		domain.Registrar,
		domain.ExpiryDate,
		domain.Status,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *DomainRepository) Delete(id int) error {
	query := `
	Delete from domains
	where id = $1;
	`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
