package data

import (
	"database/sql"
	"errors"
	"greenlight/internal/validator"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"` // Used for optimistic locking - preventing race condition when updating
}

func ValidateMovieInput(v *validator.Validator, movie Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
}

type MoviesModel struct {
	DB *sql.DB
}

func (m MoviesModel) Insert(movie *Movie) error {
	statement := "INSERT INTO MOVIES (title, year, runtime, genres) VALUES ($1, $2, $3, $4) RETURNING id, version, created_at"
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	row := m.DB.QueryRow(statement, args...)
	return row.Scan(&movie.ID, &movie.Version, &movie.CreatedAt)
}

func (m MoviesModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrorRecordNotFound
	}
	statement := "SELECT id, created_at, title, year, runtime, genres,  version FROM MOVIES where id=$1"
	row := m.DB.QueryRow(statement, id)
	movie := Movie{}
	if err := row.Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m MoviesModel) Update(movie *Movie) error {
	/*
		'version' query condition is used for optimistic locking,
		preventing updates to happens concurrently between 2 client API calls,
		instead return an error
	*/
	statement := "UPDATE MOVIES SET title=$1, year=$2, runtime=$3, genres=$4, version=version + 1 WHERE id=$5 AND version=$6 RETURNING version"
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	// Get new version into arguments YO
	err := m.DB.QueryRow(statement, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrorEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m MoviesModel) Delete(id int64) error {
	if id < 1 {
		return ErrorRecordNotFound
	}
	statement := "DELETE FROM MOVIES WHERE id=$1"
	result, err := m.DB.Exec(statement, id)
	if err != nil {
		return err
	}
	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

type MoviesRepository interface {
	Insert(movie *Movie) error
	Get(id int64) (*Movie, error)
	Update(movie *Movie) error
	Delete(id int64) error
}
