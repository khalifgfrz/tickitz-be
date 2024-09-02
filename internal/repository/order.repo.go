package repository

import (
	"khalifgfrz/tickitz-be/internal/models"

	"github.com/jmoiron/sqlx"
)

type OrderRepositoryInterface interface {
	CreateData(body *models.Order) (string, error)
	GetAllData() (*models.GetOrders, error)
	GetDetailData(id string) (*models.GetOrder, error)
	GetHistoryOrder(id string) ([]models.GetOrder, error)
}

type OrderRepository struct {
	*sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) CreateData(body *models.Order) (string, error) {
	query := `INSERT INTO public.orders
(user_id, payment_method_id, movie_id, date, time, seat_count, total, ticket_status)
VALUES($1, $2, $3, $4, $5, $6, $7, $8) returning id;`
	params := []interface{}{
		body.User_id,
		body.Payment_method_id,
		body.Movie_id,
		body.Date,
		body.Time,
		body.Seat_count,
		body.Total,
		body.Ticket_status,
	}

	var id string
	err := r.Get(&id, query, params...)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *OrderRepository) GetAllData() (*models.GetOrders, error) {
	query := `select id, order_number, user_id, payment_method_id, movie_id, date, time, seat_count, ticket_status, total from orders`
	data := models.GetOrders{}

	err := r.Select(&data, query)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *OrderRepository) GetDetailData(id string) (*models.GetOrder, error) {
	query := `SELECT 
    o.id,
    m.title AS movie_title,
    o.seat_count,
    o.ticket_status,
    o.total,
    o.date AS date,
    o.time AS time,
    COALESCE(STRING_AGG(DISTINCT g.name, ', '), '') AS genres
	FROM 
    public.orders o
	JOIN 
    public.movies m ON o.movie_id = m.id
	LEFT JOIN 
    public.genre_movies gm ON m.id = gm.movie_id
	LEFT JOIN 
    public.genres g ON gm.genre_id = g.id
	LEFT JOIN 
    public.order_details od ON od.order_id = o.id
	WHERE 
    o.id = $1
	GROUP BY 
    o.id, m.title, o.seat_count, o.ticket_status, o.total, o.date, o.time
	ORDER BY 
    o.date ASC;
`
	data := models.GetOrder{}

	err := r.Get(&data, query, id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *OrderRepository) GetHistoryOrder(id string) ([]models.GetOrder, error) {
	query := `
	SELECT 
    o.id,
		o.order_number,
    m.title AS movie_title,
    o.seat_count,
    o.ticket_status,
    o.total,
    o.date AS date,
    o.time AS time,
    COALESCE(STRING_AGG(DISTINCT g.name, ', '), '') AS genres
	FROM 
    public.orders o
	JOIN 
    public.movies m ON o.movie_id = m.id
	LEFT JOIN 
    public.genre_movies gm ON gm.movie_id = m.id
	LEFT JOIN 
    public.genres g ON gm.genre_id = g.id
	WHERE 
    o.user_id = $1
	GROUP BY 
    o.id, o.order_number, m.title, o.seat_count, o.ticket_status, o.total, o.date, o.time
	ORDER BY 
    o.date ASC;

	`
	var data []models.GetOrder

	err := r.Select(&data, query, id)
	if err != nil {
		return nil, err
	}
	return data, nil
}
