package repository

import (
	"database/sql"
	"fmt"
	"khalifgfrz/tickitz-be/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

type MovieRepoInterface interface {
	CreateMovie(tx *sqlx.Tx, data *models.MoviesBody) (*models.Movies, error)
	GetAllMovies(query *models.MoviesQuery) (*models.MovieResponse, int, error)
	GetDetailMovie(id string) (*models.MovieDetails, error)
	UpdateMovieDetails(tx *sqlx.Tx, id string, data *models.MoviesBody) (*models.Movies, error)
	UpdateBannerMovie(id string, data *models.MoviesBanner) (string, error)
	DeleteMovie(id string) (string, error)
}

type RepoMovies struct {
	*sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) *RepoMovies {
	return &RepoMovies{db}
}

func (r *RepoMovies) CreateMovie(tx *sqlx.Tx, data *models.MoviesBody) (*models.Movies, error) {
	query := `
        INSERT INTO public.movies (
            "title",
            "image", 
            "director", 
            "casts",
            "duration",
            "release_date",
            "synopsis",
            "is_deleted"
        ) VALUES (
            :title,
            :image,
            :director,
            :casts,
            :duration,
            :release_date,
            :synopsis,
            FALSE
        ) RETURNING id, title, image, director, casts, duration, release_date, synopsis, created_at;
    `

	var result models.Movies

	// Use tx.NamedQuery to execute the query within the transaction
	rows, err := tx.NamedQuery(query, data)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}

func (r *RepoMovies) GetAllMovies(query *models.MoviesQuery) (*models.MovieResponse, int, error) {
	baseQuery := `
		SELECT 
			"m".id, 
			"m".title, 
			"m".image, 
			"m".banner,
			COALESCE(STRING_AGG("g"."name", ', '), '') AS genres,
			"m".director, 
			"m".casts, 
			"m".duration, 
			"m".release_date, 
			"m".synopsis, 
			"m".is_deleted, 
			"m".created_at, 
			"m".updated_at
		FROM public.movies "m"
		LEFT JOIN public.genre_movies gm ON "m"."id" = "gm"."movie_id"
		LEFT JOIN public.genres "g" ON "gm"."genre_id" = "g"."id"
	`

	countQuery := `
		SELECT COUNT(DISTINCT "m"."id") 
		FROM public.movies m
		LEFT JOIN public.genre_movies gm ON "m"."id" = "gm"."movie_id"
		LEFT JOIN public.genres "g" ON "gm"."genre_id" = "g"."id"
		WHERE "m"."is_deleted" = FALSE
	`

	whereClauses := []string{}
	var values []interface{}

	if query.Search != nil {
		searchTerm := "%" + *query.Search + "%"
		whereClauses = append(whereClauses, fmt.Sprintf(`"m"."title" ILIKE $%d`, len(values)+1))
		values = append(values, searchTerm)
	}

	if query.Filter != nil {
		genres := strings.Split(*query.Filter, ",")
		for _, genre := range genres {
			genre = strings.TrimSpace(genre)
			whereClauses = append(whereClauses, fmt.Sprintf(`EXISTS (
				SELECT 1 
				FROM public.genre_movies gm2
				JOIN public.genres g2 ON gm2.genre_id = g2.id
				WHERE gm2.movie_id = m.id AND g2.name ILIKE $%d
			)`, len(values)+1))
			values = append(values, genre)
		}
	}

	if len(whereClauses) > 0 {
		whereQuery := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereQuery
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += ` GROUP BY "m"."id"`

	if query.Page > 0 && query.Limit > 0 {
		limit := query.Limit
		offset := (query.Page - 1) * limit
		baseQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, len(values)+1, len(values)+2)
		values = append(values, limit, offset)
	}

	var data models.MovieResponse
	if err := r.Select(&data, baseQuery, values...); err != nil {
		return nil, 0, err
	}

	var total int
	if err := r.Get(&total, countQuery, values[:len(values)-2]...); err != nil {
		return nil, 0, err
	}

	return &data, total, nil
}

func (r *RepoMovies) GetDetailMovie(id string) (*models.MovieDetails, error) {
	query := `
    SELECT 
				m.id, 
				m.title, 
				m.image, 
				m.banner,
				COALESCE(STRING_AGG(DISTINCT g.name, ', '), '') AS genres,
				m.director, 
				m.casts, 
				m.duration, 
				m.release_date, 
				m.synopsis, 
				m.is_deleted, 
				m.created_at, 
				m.updated_at,
    COALESCE(
        STRING_AGG(
            DISTINCT 
            CONCAT(
                ad.start_date::text, ' - ', 
                ad.end_date::text
            ), 
            ', '
        ), 
        ''
    ) AS airing_dates,
    COALESCE(
        STRING_AGG(
            DISTINCT at.time::text, 
            ', '
        ), 
        ''
    ) AS airing_times,
    COALESCE(
        STRING_AGG(
            DISTINCT l.name, 
            ', '
        ), 
        ''
    ) AS locations
		FROM public.movies m
		LEFT JOIN public.genre_movies gm ON m.id = gm.movie_id
		LEFT JOIN public.genres g ON gm.genre_id = g.id
		LEFT JOIN public.movies_time mt ON m.id = mt.movie_id
		LEFT JOIN public.airing_time_date atd ON mt.airing_time_date_id = atd.id
		LEFT JOIN public.airing_date ad ON atd.date_id = ad.id
		LEFT JOIN public.airing_time at ON atd.airing_time_id = at.id
		LEFT JOIN public.location_movie_time lmt ON mt.id = lmt.movie_time_id
		LEFT JOIN public.locations l ON lmt.location_id = l.id
		WHERE m.id = $1 AND m.is_deleted = FALSE
		GROUP BY 
				m.id, 
				m.title, 
				m.image, 
				m.banner,
				m.director, 
				m.casts, 
				m.duration, 
				m.release_date, 
				m.synopsis, 
				m.is_deleted, 
				m.created_at, 
				m.updated_at;
  `

	var result models.MovieDetails

	row := r.DB.QueryRow(query, id)
	err := row.Scan(
		&result.Id,
		&result.Title,
		&result.Image,
		&result.Banner,
		&result.Genres,
		&result.Director,
		&result.Casts,
		&result.Duration,
		&result.Release_Date,
		&result.Synopsis,
		&result.Is_deleted,
		&result.Created_at,
		&result.Updated_at,
		&result.AiringDates,
		&result.AiringTimes,
		&result.Locations,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *RepoMovies) UpdateMovieDetails(tx *sqlx.Tx, id string, data *models.MoviesBody) (*models.Movies, error) {
	// Update query
	query := `
		UPDATE public.movies
		SET title = COALESCE(NULLIF($1, ''), title),
			image = COALESCE(NULLIF($2, ''), image),
			director = COALESCE(NULLIF($3, ''), director),
			casts = COALESCE(NULLIF($4, ''), casts),
			duration = COALESCE(NULLIF($5, ''), duration),
			release_date = COALESCE(NULLIF($6, release_date::date), release_date::date),
			synopsis = COALESCE(NULLIF($7, ''), synopsis)
		WHERE id = $8
		RETURNING id, title, image, director, casts, duration, release_date, synopsis, created_at, updated_at;
	`

	// Execute update query using tx
	var result models.Movies
	row := tx.QueryRow(query,
		data.Title,
		data.Image,
		data.Director,
		data.Casts,
		data.Duration,
		data.Release_Date,
		data.Synopsis,
		id,
	)

	err := row.Scan(
		&result.Id,
		&result.Title,
		&result.Image,
		&result.Director,
		&result.Casts,
		&result.Duration,
		&result.Release_Date,
		&result.Synopsis,
		&result.Created_at,
		&result.Updated_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows updated for id %s", id)
		}
		return nil, fmt.Errorf("failed to update movie details: %w", err)
	}

	// Query to get genres
	genreQuery := `
		SELECT COALESCE(STRING_AGG(g.name, ', '), '') AS genres
		FROM public.genre_movies gm
		JOIN public.genres g ON gm.genre_id = g.id
		WHERE gm.movie_id = $1
	`

	// Execute genre query using tx
	var genres string
	err = tx.QueryRow(genreQuery, id).Scan(&genres)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve genres: %w", err)
	}

	// Add genres to the result
	result.Genres = &genres

	return &result, nil
}

func (r *RepoMovies) UpdateBannerMovie(id string, data *models.MoviesBanner) (string, error) {
	query := `
    UPDATE public.movies
    SET
      "banner" = COALESCE(NULLIF(:banner, ''), "banner"),
			"updated_at" = now()
    WHERE "id" = :id
  `
	data.Id = id

	var updatedMovie models.Movies
	rows, err := r.DB.NamedQuery(query, data)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.StructScan(&updatedMovie)
		if err != nil {
			return "", err
		}
	}

	return "Data updated", nil
}

func (r *RepoMovies) DeleteMovie(id string) (string, error) {
	query := `
    UPDATE public.movies
		SET
			"is_deleted" = TRUE,
      "updated_at" = now()
    WHERE "id" = $1
    RETURNING *
  `

	var deletedMovie models.Movies
	if err := r.DB.QueryRowx(query, id).StructScan(&deletedMovie); err != nil {
		return "", err
	}

	return "Data deleted", nil
}
