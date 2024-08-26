package models

import "time"

var schemaPayments = `
CREATE TABLE public.payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

type Payments struct {
	ID         int        `db:"id" json:"id"`
	Name       string     `db:"name" json:"name"`
	Image      string     `db:"image" json:"image"`
	Created_at *time.Time `db:"created_at" json:"created_at"`
	Updated_at *time.Time `db:"updated_at" json:"updated_at"`
}
