package repository

import (
	models "khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type AiringTimeRepoInterface interface {
	GetAiringTime() ([]models.AiringTime, error)
}

type RepoAiringTime struct {
	*sqlx.DB
}

func NewAiringTimeRepository(db *sqlx.DB) *RepoAiringTime {
	return &RepoAiringTime{db}
}

//mthods get airing_time from database

func (r *RepoAiringTime) GetAiringTime() ([]models.AiringTime, error) {
	query := `SELECT * FROM airing_time`

	var airingTimes []models.AiringTime
	err := r.DB.Select(&airingTimes, query)

	return airingTimes, err
}
