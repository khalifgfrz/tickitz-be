package repository

import (
	"fmt"
	"khalifgfrz/tickitz-be/internal/models/moviesAdd"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AiringTimeDateRepoInterface interface {
	GetAiringTimeDate() ([]moviesAdd.AiringTimeDate, error)
	UpdateMovieAiringDetails(tx *sqlx.Tx, movieID string, airingDates []string, airingTimeIDs []int) error
	InsertAiringTimeDate(tx *sqlx.Tx, airingTimeDate *moviesAdd.AiringTimeDate) (*moviesAdd.AiringTimeDate, error)
	DeleteAiringTimeDatesByMovieID(tx *sqlx.Tx, movieID string) error
}

type RepoAiringTimeDate struct {
	tx *sqlx.Tx
	db *sqlx.DB
	*RepoAiringDate
}

// Initialize RepoAiringTimeDate properly
func NewAiringTimeDateRepository(db *sqlx.DB, tx *sqlx.Tx) *RepoAiringTimeDate {
	return &RepoAiringTimeDate{
		tx:             tx,
		db:             db,
		RepoAiringDate: &RepoAiringDate{}, // Ensure RepoAiringDate is properly initialized if needed
	}
}

func (r *RepoAiringTimeDate) GetAiringTimeDate() ([]moviesAdd.AiringTimeDate, error) {
	query := `SELECT * FROM airing_time_date`

	var airingTimeDates []moviesAdd.AiringTimeDate
	err := r.DB.Select(&airingTimeDates, query)

	return airingTimeDates, err
}

func (r *RepoAiringTimeDate) InsertAiringTimeDate(tx *sqlx.Tx, airingTimeDate *moviesAdd.AiringTimeDate) (*moviesAdd.AiringTimeDate, error) {
	query := `INSERT INTO airing_time_date (airing_time_id, date_id) VALUES (:airing_time_id,:date_id) RETURNING *`

	var results moviesAdd.AiringTimeDate
	rows, err := tx.NamedQuery(query, airingTimeDate)
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

func (r *RepoAiringDate) DeleteAiringTimeDatesByMovieID(tx *sqlx.Tx, movieID string) error {
	query := `DELETE FROM airing_time_date
						WHERE id IN (
								SELECT airing_time_date_id
								FROM movies_time
								WHERE movie_id = $1
						);`
	_, err := tx.Exec(query, movieID)
	return err
}

func (r *RepoAiringDate) UpdateMovieAiringDetails(tx *sqlx.Tx, movieID string, airingDates []string, airingTimeIDs []int) error {
	// Delete existing airing time dates
	queryDelete := `
		DELETE FROM public.airing_time_date
		WHERE date_id IN (
			SELECT mt.airing_time_date_id
			FROM public.movies_time mt
			WHERE mt.movie_id = $1
		);
	`
	if _, err := tx.Exec(queryDelete, movieID); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	for _, dateRange := range airingDates {
		dates := strings.Split(dateRange, ",")
		var airingDate moviesAdd.AiringDate

		if len(dates) == 2 {
			airingDate = moviesAdd.AiringDate{
				Start_date: strings.TrimSpace(dates[0]),
				End_date:   strings.TrimSpace(dates[1]),
			}
		} else if len(dates) == 1 {
			airingDate = moviesAdd.AiringDate{
				Start_date: strings.TrimSpace(dates[0]),
				End_date:   strings.TrimSpace(dates[0]),
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("invalid date range: %s", dateRange)
		}

		existingAiringDate, err := r.GetAiringDateByInput(tx, &airingDate)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error getting airing date: %w", err)
		}

		if existingAiringDate != nil {
			airingDate.Id = existingAiringDate.Id
		} else {
			insertedAiringDates, err := r.CreateAiringDate(tx, &airingDate)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error creating airing date: %w", err)
			}
			if len(insertedAiringDates) > 0 {
				airingDate.Id = insertedAiringDates[0].Id
			} else {
				tx.Rollback()
				return fmt.Errorf("no airing date was inserted")
			}
		}

		for _, airingTimeId := range airingTimeIDs {
			newAiringTimeDate := moviesAdd.AiringTimeDate{
				Airing_time_id: airingTimeId,
				Date_id:        airingDate.Id,
			}
			if _, err := r.InsertAiringTimeDate(tx, &newAiringTimeDate); err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting airing time date: %w", err)
			}
		}
	}

	// Commit transaction if successful
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
