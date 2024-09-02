package repository

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type GenreRepoInterface interface {
	CreateGenres(data *models.Genres) (string, error)
	GetAllGenres() ([]models.Genres, error)
}

type RepoGenres struct {
	*sqlx.DB
}

func NewGenresRepository(db *sqlx.DB) *RepoGenres {
	return &RepoGenres{db}
}

func (r *RepoGenres) CreateGenres(data *models.Genres) (string, error) {
	query := `INSERT INTO genres (name) VALUES (:name) RETURNING *`

	var result models.Genres
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

func (r *RepoGenres) GetAllGenres() ([]models.Genres, error) {
	query := `SELECT id, name FROM genres`

	var data []models.Genres
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}
