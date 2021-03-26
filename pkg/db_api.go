package src

import (
	"database/sql"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {

	return &Repository{
		db: db,
	}

}

func (r *Repository) IsUserInDb(name string) bool {

	log.Printf("Checking if user is in DB:\n")

	//query := ""

	return false

}
