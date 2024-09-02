package repository

import (
	"fmt"
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type GenreMovieRepoInterface interface {
	InsertGenreMovie(tx *sqlx.Tx, data *models.GenreMovie) (*models.GenreMovie, error)
	UpdateGenreMovie(tx *sqlx.Tx, id string, genreIDs []int) error
	DeleteGenreMovie(id string) (string, error)
}

type RepoGenreMovie struct {
	*sqlx.DB
}

func NewGenreMovieRepository(db *sqlx.DB) *RepoGenreMovie {
	return &RepoGenreMovie{db}
}

func (r *RepoGenreMovie) InsertGenreMovie(tx *sqlx.Tx, data *models.GenreMovie) (*models.GenreMovie, error) {
	if data.Movie_id == "" || data.Genre_id == 0 {
		return nil, fmt.Errorf("genre_id and movie_id cannot be empty or zero")
	}

	query := `
    INSERT INTO public.genre_movies (
      genre_id,
      movie_id
    ) VALUES (
      :genre_id,
      :movie_id
    ) RETURNING id, genre_id, movie_id, created_at, updated_at;
  `
	var results models.GenreMovie
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
	} else {
		return nil, fmt.Errorf("no rows returned after insert")
	}

	return &results, nil
}

func (r *RepoGenreMovie) UpdateGenreMovie(tx *sqlx.Tx, id string, genreIDs []int) error {
	if id == "" {
		return fmt.Errorf("movie_id cannot be empty")
	}

	if len(genreIDs) < 1 {
		return fmt.Errorf("genre_ids cannot be empty or null")
	}

	// Delete existing genres
	queryDelete := `DELETE FROM public.genre_movies WHERE movie_id = $1`
	if _, err := tx.Exec(queryDelete, id); err != nil {
		return err
	}

	// Insert new genres
	queryInsert := `INSERT INTO public.genre_movies (movie_id, genre_id) VALUES ($1, $2)`

	for _, genreID := range genreIDs {
		if _, err := tx.Exec(queryInsert, id, genreID); err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoGenreMovie) DeleteGenreMovie(id string) (string, error) {
	query := `
    DELETE FROM public.genre_movies
    WHERE id = $1
    RETURNING *;
  `

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
