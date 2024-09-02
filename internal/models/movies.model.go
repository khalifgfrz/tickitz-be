package models

import (
	"encoding/json"
	"time"
)

var schemaMovies = `
CREATE TABLE public.movies (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	title varchar(255) NULL,
	image text NULL,
	director varchar(50) NULL,
	casts text NULL,
	duration varchar(20) NULL,
	release_date date NULL,
	synopsis text NULL,
	created_at timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp NULL,
	is_deleted bool NULL DEFAULT false,
	banner text NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id),
	CONSTRAINT unique_title UNIQUE (title)
);
`

type Movies struct {
	Id           string     `db:"id" json:"id"`
	Title        string     `db:"title" json:"title" form:"title" valid:"stringlength(2|100)~Nama Movies minimal 2 dan maksimal 100"`
	Image        string     `db:"image" json:"image" valid:"-"`
	Banner       *string    `db:"banner" json:"banner,omitempty"`
	Genres       *string    `db:"genres" json:"genres" form:"genres" valid:"-"`
	Director     string     `db:"director" json:"director" form:"director" valid:"stringlength(5|100)~Nama Product minimal 5 dan maksimal 100"`
	Casts        string     `db:"casts" json:"casts" form:"casts" valid:"-"`
	Duration     string     `db:"duration" json:"duration" form:"duration" valid:"-"`
	Release_Date time.Time  `db:"release_date" json:"release_date" form:"release_date" valid:"date"`
	Synopsis     string     `db:"synopsis" json:"synopsis" form:"synopsis" valid:"-"`
	Is_deleted   bool       `db:"is_deleted" json:"is_deleted"`
	Created_at   *time.Time `db:"created_at" json:"created_at"`
	Updated_at   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type MovieDetails struct {
	Movies
	Genres      string `json:"genres"`       // Comma-separated list of genres
	AiringDates string `json:"airing_dates"` // Comma-separated list of airing date ranges (start_date - end_date)
	AiringTimes string `json:"airing_times"` // Comma-separated list of airing times
	Locations   string `json:"locations"`    // Comma-separated list of locations
}

type AllMovies struct {
	Id           string     `db:"id" json:"id"`
	Title        string     `db:"title" json:"title"`
	Image        string     `db:"image" json:"image"`
	Banner       *string    `db:"banner" json:"banner,omitempty"`
	Genres       string     `db:"genres" json:"genres"`
	Director     string     `db:"director" json:"director"`
	Casts        string     `db:"casts" json:"casts"`
	Duration     string     `db:"duration" json:"duration"`
	Release_Date *time.Time `db:"release_date" json:"release_date"`
	Synopsis     string     `db:"synopsis" json:"synopsis"`
	Is_deleted   bool       `db:"is_deleted" json:"is_deleted"`
	Created_at   *time.Time `db:"created_at" json:"created_at"`
	Updated_at   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

func (u AllMovies) MarshalJSON() ([]byte, error) {
	type Alias AllMovies
	return json.Marshal(&struct {
		*Alias
		Release_Date string `json:"release_date"`
	}{
		Release_Date: formatDate(u.Release_Date),
		Alias:        (*Alias)(&u),
	})
}

func formatDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

type MovieResponse []AllMovies

type MoviesQuery struct {
	Page   int     `form:"page"`
	Limit  int     `form:"limit"`
	Search *string `form:"name"`
	Filter *string `form:"genres"`
}

type MoviesBanner struct {
	Id     string `db:"id" json:"id"`
	Banner string `db:"banner" json:"banner"`
}

type MoviesBody struct {
	Title        *string `db:"title" json:"title" form:"title" valid:"stringlength(2|100)~Nama Movies minimal 2 dan maksimal 100"`
	Image        *string `db:"image" json:"image" valid:"-"`
	Genres       *string `json:"genres" form:"genres"`
	Director     *string `db:"director" json:"director" form:"director" valid:"stringlength(5|100)~Nama Product minimal 5 dan maksimal 100"`
	Casts        *string `db:"casts" json:"casts" form:"casts" valid:"-"`
	Duration     *string `db:"duration" json:"duration" form:"duration" valid:"-"`
	Release_Date *string `db:"release_date" json:"release_date" form:"release_date" valid:"-"`
	Synopsis     *string `db:"synopsis" json:"synopsis" form:"synopsis" valid:"-"`
	AiringDate   *string `json:"airing_date" form:"airing_date" valid:"-"`
	AiringTime   *string `json:"airing_time" form:"airing_time"`
	Locations    *string `json:"locations" form:"locations"`
}
