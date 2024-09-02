package repository

import (
	"khalifgfrz/tickitz-be/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderDetailsRepositoryInterface interface {
	CreateOrderDetails(order_id string, orders []models.OrderDetails) (string, error)
	GetDetailOrder(order_id string) ([]models.GetOrderDetails, error)
}

type OrderDetailsRepository struct {
	*sqlx.DB
}

func NewOrderDetailsRepository(db *sqlx.DB) *OrderDetailsRepository {
	return &OrderDetailsRepository{db}
}

func (r *OrderDetailsRepository) CreateOrderDetails(order_id string, orders []models.OrderDetails) (string, error) {
	query := `INSERT INTO public.order_details
	(order_id, seat_id)
	VALUES(:order_id, :seat_id);`

	for _, order := range orders {
		data := map[string]interface{}{
			"order_id": order_id,
			"seat_id":  order.Seat_id,
		}
		_, err := r.NamedExec(query, data)
		if err != nil {
			return "", err
		}
	}
	return order_id, nil
}

func (r *OrderDetailsRepository) GetDetailOrder(order_id string) ([]models.GetOrderDetails, error) {
	query := `SELECT 
        od.order_id, 
        array_agg(od.seat_id) AS seat_id, 
        array_agg(s.name) AS seat_names
    FROM 
        order_details od
    JOIN 
        seats s 
    ON 
        od.seat_id = s.id
    WHERE 
        od.order_id = $1
    GROUP BY 
        od.order_id;`

	var data []models.GetOrderDetails

	err := r.Select(&data, query, order_id)
	if err != nil {
		return nil, err
	}

	return data, nil
}
