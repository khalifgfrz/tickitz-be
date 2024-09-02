package repository

import (
	"database/sql"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type AiringDateRepoInterface interface {
	CreateAiringDate(tx *sqlx.Tx, dates *models.AiringDate) ([]models.AiringDate, error)
	GetAiringDate() (*models.AiringDate, error)
	GetAiringDateByInput(tx *sqlx.Tx, input *models.AiringDate) (*models.AiringDate, error)
}

type RepoAiringDate struct {
	*sqlx.DB
	*RepoAiringTimeDate
}

func NewAiringDateRepository(db *sqlx.DB) *RepoAiringDate {
	return &RepoAiringDate{db, &RepoAiringTimeDate{}}
}

// insert dynamically
func (r *RepoAiringDate) CreateAiringDate(tx *sqlx.Tx, dates *models.AiringDate) ([]models.AiringDate, error) {
	query := `INSERT INTO airing_date (start_date, end_date) 
	VALUES (:start_date, :end_date)
	RETURNING id, start_date, end_date, created_at, updated_at`

	var results []models.AiringDate
	rows, err := tx.NamedQuery(query, *dates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var result models.AiringDate
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// methods get all airing_date
func (r *RepoAiringDate) GetAiringDate() (*models.AiringDate, error) {
	query := `SELECT * FROM airing_date`

	var airingDate models.AiringDate
	err := r.DB.Get(&airingDate, query)

	return &airingDate, err
}

func (r *RepoAiringDate) GetAiringDateByInput(tx *sqlx.Tx, input *models.AiringDate) (*models.AiringDate, error) {
	query := `SELECT id, start_date, end_date, created_at, updated_at FROM airing_date
	WHERE start_date = $1 AND end_date = $2`

	var airingDate models.AiringDate
	err := tx.Get(&airingDate, query, input.Start_date, input.End_date)

	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were found, return nil without error
			return nil, nil
		}
		// For other errors, return the error
		return nil, err
	}
	return &airingDate, nil
}
