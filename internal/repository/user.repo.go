package repository

import (
	"khalifgfrz/tickitz-be/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepositoryInterface interface {
	CreateData(body *models.Auth) (string, error)
	UpdateData(data *models.User, id string) (string, error)
	// GetAllData()(*models.Users , error)
	GetDetailData(id string) (*models.UserDetails, error)
	// DeleteData(id string) (string, error)
}

type UserRepository struct {
	*sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) CreateData(body *models.Auth) (string, error) {
	query := `INSERT INTO public.users (email, password, role, is_deleted ) VALUES ( :email, :password, 'user', FALSE)`

	_, err := r.NamedExec(query, body)
	if err != nil {
		return "", err
	}

	return "Create data success", nil
}

func (r *UserRepository) UpdateData(data *models.User, id string) (string, error) {
	query := `UPDATE public.users SET
        first_name = COALESCE(NULLIF($1, ''), first_name),
		last_name = COALESCE(NULLIF($2, ''), last_name),
		email = COALESCE(NULLIF($3, ''), email),
		password = COALESCE(NULLIF($4, ''), password),
		image = COALESCE(NULLIF($5, ''), image),
		phone_number = COALESCE(NULLIF($6, ''), phone_number),
		point = COALESCE(NULLIF($7, ''), point),
		updated_at = now()
	WHERE id = $8`

	params := []interface{}{
		data.First_name,
		data.Last_name,
		data.Email,
		data.Password,
		data.Image,
		data.Phone_number,
		data.Point,
		id,
	}

	_, err := r.Exec(query, params...)
	if err != nil {
		return "", err
	}
	return "Update data success", nil
}

// func (r *UserRepository) GetAllData()(*models.Users, error){
// 	query := `select id, first_name, last_name, email  FROM users where is_deleted = FALSE`
// 	data := models.Users{}

// 	err := r.Select(&data, query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &data, nil
// }

func (r *UserRepository) GetDetailData(id string) (*models.UserDetails, error) {
	query := `select id, first_name, last_name, email, image,phone_number, point FROM public.users WHERE id = $1 AND is_deleted = FALSE `
	data := models.UserDetails{}
	err := r.Get(&data, query, id)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// func (r *UserRepository) DeleteData(id string) (string, error) {
// 	query := `UPDATE public.users SET is_deleted = true WHERE id = $1`
// 	_, err := r.Exec(query, id)
// 	if err != nil {
// 		return "", err
// 	}

// 	return "Delete success", nil
// }
