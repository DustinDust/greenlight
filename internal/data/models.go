package data

import (
	"database/sql"
	"errors"
)

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MoviesRepository
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MoviesModel{DB: db},
	}
}

// Demonstrate creating a mock model using repository interface
// Used for testing
type MockMovieModel struct {
	DB *sql.DB
}

func (mm MockMovieModel) FindAll(title string, geners []string, filter Filter) ([]*Movie, Metadata, error) {
	return make([]*Movie, 0), Metadata{}, nil
}

func (mm MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (mm MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (mm MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (mm MockMovieModel) Delete(id int64) error {
	return nil
}

func NewMockModels(db *sql.DB) Models {
	return Models{
		Movies: MockMovieModel{DB: db},
	}
}
