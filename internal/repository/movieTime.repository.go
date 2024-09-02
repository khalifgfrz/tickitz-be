package repository

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type MovieTimeRepoInteface interface {
	CreateMovieTime(tx *sqlx.Tx, data *models.MovieTime) (*models.MovieTime, error)
	GetTimeByMovieId(id string) (*models.MovieTime, error)
	UpdateMovieTime(id string, data *models.MovieTime) (string, error)
	DeleteMovieTime(id int) (string, error)
}

type MovieTimeRepository struct {
	*sqlx.DB
}

func NewMovieTimeRepository(db *sqlx.DB) *MovieTimeRepository {
	return &MovieTimeRepository{db}
}

func (r *MovieTimeRepository) CreateMovieTime(tx *sqlx.Tx, data *models.MovieTime) (*models.MovieTime, error) {
	query := `
		INSERT INTO public.movies_time (
      "movie_id",
      "airing_time_date_id"
    ) VALUES (
      :movie_id,
			:airing_time_date_id
    ) RETURNING *;
		`
	var results models.MovieTime
	rows, err := tx.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&results)
		if err != nil {
			return nil, err
		}
	}

	return &results, nil
}

func (r *MovieTimeRepository) GetTimeByMovieId(id string) (*models.MovieTime, error) {
	query := `SELECT * FROM public.movies_time WHERE "movie_id" = :movie_id;`

	data := models.MovieTime{}
	rows, err := r.DB.NamedQuery(query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		return &data, nil
	}

	return nil, nil
}

func (r *MovieTimeRepository) UpdateMovieTime(id string, data *models.MovieTime) (string, error) {
	query := `
    UPDATE public.movies_time
    SET "airing_time_date_id" = :airing_time_date_id
    WHERE "movie_id" = :movie_id
    RETURNING *;
  `

	data.Movie_id = id

	var results models.MovieTime
	rows, err := r.DB.NamedQuery(query, data)
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

	return "Data updated", nil
}

func (r *MovieTimeRepository) DeleteMovieTime(id int) (string, error) {
	query := `
    DELETE FROM public.movies_time
    WHERE "id" = :id
    RETURNING *;
  `

	var results models.MovieTime
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
