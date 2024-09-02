// create repository package
package repository

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type LocationRepoInterface interface {
	CreateLocation(data *models.Locations) (string, error)
	GetAllLocations() ([]models.Locations, error)
	UpdateLocation(id int, data *models.Locations) (*models.Locations, error)
	DeleteLocation(id int) (string, error)
}

type LocationRepo struct {
	*sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) *LocationRepo {
	return &LocationRepo{db}
}

func (r *LocationRepo) CreateLocation(data *models.Locations) (string, error) {
	query := `INSERT INTO public.location ("name") VALUES (:name)`
	var result models.Locations
	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return "", err
		}
	}

	return "Data created", nil
}

func (r *LocationRepo) GetAllLocations() ([]models.Locations, error) {
	query := `SELECT id, name FROM locations`

	var data []models.Locations
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *LocationRepo) UpdateLocation(id int, data *models.Locations) (*models.Locations, error) {
	query := `UPDATE public.location SET "name" = :name WHERE id = :id RETURNING *`

	data.ID = id

	var result models.Locations
	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (r *LocationRepo) DeleteLocation(id int) (string, error) {
	query := `DELETE FROM public.location WHERE id = $1`

	var results models.Locations
	rows, err := r.DB.Queryx(query, id)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&results)
		if err != nil {
			return "", err
		}
	}

	return "Data deleted", nil
}
