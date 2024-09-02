package repository

import (
	"khalifgfrz/tickitz-be/internal/models/moviesAdd"

	"github.com/jmoiron/sqlx"
)

type SalesRepoInterface interface {
	GetAllSales() ([]moviesAdd.GetSales, error)
}

type RepoSales struct {
	*sqlx.DB
}

func NewSalesRepository(db *sqlx.DB) *RepoSales {
	return &RepoSales{db}
}

func (r *RepoSales) GetAllSales() ([]moviesAdd.GetSales, error) {
	query := `SELECT date, sales FROM dummy_sales`

	var data []moviesAdd.GetSales
	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}

	return data, nil
}
