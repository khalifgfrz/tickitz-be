package repository

import (
	"khalifgfrz/tickitz-be/internal/models"

	"github.com/jmoiron/sqlx"
)

type PaymentsRepoInterface interface {
	GetAllPayments() ([]models.Payments, error)
}

type PaymentsRepo struct {
	*sqlx.DB
}

func NewPaymentsRepository(db *sqlx.DB) *PaymentsRepo {
	return &PaymentsRepo{db}
}

func (r *PaymentsRepo) GetAllPayments() ([]models.Payments, error) {
	query := `SELECT id, name FROM payment_methods`

	var data []models.Payments
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}
