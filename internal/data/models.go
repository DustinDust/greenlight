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
