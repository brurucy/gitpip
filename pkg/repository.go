package pkg

import (
	"context"
	"database/sql"
	"fmt"
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

func (r *Repository) IsUserBeingTracked(name string) (bool, error) {

	log.Printf("Checking if user is in DB:\n")

	query := "select CASE  when EXISTS(SELECT 1 FROM users u WHERE u.username ILIKE '%' || $1 || '%') then true else false end;"

	row := r.db.QueryRowContext(context.Background(), query, name)

	if row == nil {

		log.Print("Error getting query row context\n")
		return false, nil
	}

	var isUserIn bool
	err := row.Scan(&isUserIn)

	if err != nil {
		log.Printf("Error scanning row, %v", err)
		return false, err
	}

	log.Printf("Successfully checked if user is in Db: %v", isUserIn)

	return isUserIn, nil

}

func (r *Repository) InsertUser(user *GistOwner) error {

	log.Printf("Inserting user: %v", user.Login)

	query := "INSERT INTO users (user_id, username) VALUES ($1,$2);"

	row := r.db.QueryRowContext(context.Background(), query, user.Id, user.Login)

	if row.Err() != nil {

		return fmt.Errorf("error inserting user %v", row.Err().Error())
	}

	log.Printf("Succesfully inserted user: %v", user.Login)

	return nil

}

func (r *Repository) InsertGistPgAndPipe(gist *Gist) error {

	log.Printf("Inserting gist: %v", gist.GistUUID)

	isUserInDb, err := r.IsUserBeingTracked(gist.GistOwner.Login)

	if err != nil {

		return fmt.Errorf("error getting user %v", err)

	}

	if isUserInDb == false {

		log.Printf("User not in the database, adding it")

		err := r.InsertUser(&gist.GistOwner)

		if err != nil {

			return fmt.Errorf("Error creating user, aborting gist insertion %v", err)

		}

	}

	log.Printf("Inserting all gist files")

	for idx, val := range gist.GistFiles {

		query := "INSERT INTO gists (gist_id, gist_file_title, raw_url_link, user_id) VALUES ($1,$2,$3,$4);"

		row := r.db.QueryRowContext(context.Background(), query, gist.GistUUID, idx, val.RawUrl, gist.GistOwner.Id)

		if row.Err() != nil {

			return fmt.Errorf("error inserting gist %v", row.Err().Error())
		}

		log.Printf("Succesfully inserted file %s out of: %d", idx, len(gist.GistFiles))

		log.Printf("Attempting to create an activity in Pipedrive")

		rawFileText, err := GistTextDownloader(val.RawUrl)

		if err != nil {

			return fmt.Errorf("error downloading gist from rawUrl %v", err)

		}

		pipedriveRequest := &PipeAddActivityRequest{

			Subject: val.Filename,
			Note:    rawFileText,
			Done:    false,
		}

		err = NewPipedriveActivity(pipedriveRequest)

		if err != nil {

			return fmt.Errorf("failed to set pipedrive activity")

		}

		log.Printf("Succesfully created pipedrive activity with subject: %v", val.Filename)

	}

	return nil

}

func (r *Repository) IsGistInDb(id string) (bool, error) {

	log.Printf("Checking if gist is in DB:\n")

	query := "select CASE when EXISTS( SELECT 1 FROM gists g WHERE g.gist_id ILIKE '$' || $1 || '$' ) then true else false end;"

	row := r.db.QueryRowContext(context.Background(), query, id)

	if row == nil {

		log.Printf("Error getting query row context\n")
		return false, nil
	}

	var isGistIn bool
	err := row.Scan(&isGistIn)

	if err != nil {
		log.Printf("Error scanning row, %v", err)
		return false, err
	}

	log.Printf("Successfully checked if gist is in Db")

	return isGistIn, nil

}
