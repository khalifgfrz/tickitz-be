package repository

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type LocationMovieTimeRepoInterface interface {
	CreateLocationMovie(tx *sqlx.Tx, data *models.LocationMovieTime) (string, error)
	GetMovieLocTimeById(id string) (*models.LocationMovieTime, error)
	DeleteLocationMovie(id string) (string, error)
	UpdateMovieLocations(id string, locationIDs []int) error
}

type RepoLocationMovie struct {
	*sqlx.DB
}

func NewLocationMovieRepository(db *sqlx.DB) *RepoLocationMovie {
	return &RepoLocationMovie{db}
}

func (r *RepoLocationMovie) CreateLocationMovie(tx *sqlx.Tx, data *models.LocationMovieTime) (string, error) {
	query := `
    INSERT INTO public.location_movie_time (
      "location_id",
			"movie_time_id"
		) VALUES (
			:location_id,
			:movie_time_id)`

	var results models.LocationMovieTime
	rows, err := tx.NamedQuery(query, data)
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

	return "Data created", nil
}

func (r *RepoLocationMovie) GetMovieLocTimeById(id string) (*models.LocationMovieTime, error) {
	query := `
    SELECT * FROM public.location_movie_time
    WHERE "movie_time_id" = $1`

	var data models.LocationMovieTime
	if err := r.DB.QueryRowx(query, id).StructScan(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *RepoLocationMovie) DeleteLocationMovie(id string) (string, error) {
	query := `DELETE FROM public.location_movie_time WHERE "movie_time_id" = :movie_time_id`

	var results models.GenreMovie
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

func (r *RepoLocationMovie) UpdateMovieLocations(id string, locationIDs []int) error {
	queryDelete := `DELETE FROM public.location_movie_time WHERE movie_time_id IN (SELECT id FROM public.movies_time WHERE movie_id = $1)`
	if _, err := r.DB.Exec(queryDelete, id); err != nil {
		return err
	}

	for _, locationID := range locationIDs {
		queryInsert := `INSERT INTO public.location_movie_time (location_id, movie_time_id) VALUES ($1, (SELECT id FROM public.movies_time WHERE movie_id = $2))`
		if _, err := r.DB.Exec(queryInsert, locationID, id); err != nil {
			return err
		}
	}

	return nil
}
