package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"time"
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

	logrus.Tracef("Checking if user is in DB:\n")

	query := "select CASE  when EXISTS(SELECT 1 FROM users u WHERE u.username ILIKE '%' || $1 || '%') then true else false end;"

	row := r.db.QueryRowContext(context.Background(), query, name)

	if row == nil {

		log.Print("Error getting query row context\n")
		return false, nil
	}

	var isUserIn bool
	err := row.Scan(&isUserIn)

	if err != nil {
		logrus.Tracef("Error scanning row, %v", err)
		return false, err
	}

	logrus.Tracef("Successfully checked if user is in Db: %v", isUserIn)

	return isUserIn, nil

}

func (r *Repository) GetAllUsers() ([]*GistOwner, error) {

	logrus.Tracef("Getting all tracked users\n")

	query := "select u.username, u.user_id from users u"

	rows, err := r.db.QueryContext(context.Background(), query)

	if rows == nil {

		return nil, fmt.Errorf("No users yet\n")

	}

	var users []*GistOwner

	for rows.Next() {

		g := &GistOwner{}
		err := rows.Scan(&g.Login, &g.Id)

		if err != nil {
			return nil, fmt.Errorf("Error scanning users query response %v", err)
		}
		users = append(users, g)

	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Error closing pgsql rows, %v", err)
	}

	logrus.Tracef("Succesfully got all users")

	return users, nil

}

func (r *Repository) InsertUser(user *GistOwner) error {

	logrus.Tracef("Inserting user: %v", user.Login)

	query := "INSERT INTO users (user_id, username) VALUES ($1,$2);"

	row := r.db.QueryRowContext(context.Background(), query, user.Id, user.Login)

	if row.Err() != nil {

		return fmt.Errorf("error inserting user %v", row.Err().Error())
	}

	logrus.Tracef("Succesfully inserted user: %v", user.Login)

	return nil

}

func (r *Repository) InsertGistPgAndPipe(gist *Gist) error {

	logrus.Tracef("Inserting gist: %v", gist.GistUUID)

	isUserInDb, err := r.IsUserBeingTracked(gist.GistOwner.Login)

	if err != nil {

		return fmt.Errorf("error getting user %v", err)

	}

	if isUserInDb == false {

		logrus.Tracef("User not in the database, adding it")

		err := r.InsertUser(&gist.GistOwner)

		if err != nil {

			return fmt.Errorf("Error creating user, aborting gist insertion %v", err)

		}

	}

	logrus.Tracef("Inserting all gist files")

	for idx, val := range gist.GistFiles {

		query := "INSERT INTO gists (gist_id, gist_file_title, raw_url_link, username) VALUES ($1,$2,$3,$4);"

		row := r.db.QueryRowContext(context.Background(), query, gist.GistUUID, idx, val.RawUrl, gist.GistOwner.Login)

		if row.Err() != nil {

			return fmt.Errorf("error inserting gist %v", row.Err().Error())
		}

		logrus.Tracef("Succesfully inserted file %s out of: %d", idx, len(gist.GistFiles))

		logrus.Tracef("Attempting to create an activity in Pipedrive")

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

		logrus.Tracef("Succesfully created pipedrive activity with subject: %v", val.Filename)

	}

	return nil

}

func (r *Repository) NewRoutine() (*Routine, error) {

	query := "INSERT INTO routine (routine_id) VALUES (DEFAULT) RETURNING routine_id, created_at"

	row := r.db.QueryRowContext(context.Background(), query)

	if row.Err() != nil {

		return nil, fmt.Errorf("error inserting new routine %v", row.Err().Error())
	}

	logrus.Tracef("Succesfully inserted new routine\n")

	var routine Routine

	_ = row.Scan(&routine.Id, &routine.CreatedAt)

	return &routine, nil

}

func (r *Repository) NewSession(username string) (*Session, error) {

	query := "INSERT INTO session (session_id, user_id) VALUES (DEFAULT, (SELECT user_id FROM users u WHERE u.username ILIKE '%' || $1 || '%')) RETURNING session_id, user_id, created_at"

	row := r.db.QueryRowContext(context.Background(), query, username)

	if row.Err() != nil {

		return nil, fmt.Errorf("error inserting new routine %v", row.Err().Error())
	}

	logrus.Tracef("Succesfully inserted new session for user %s\n", username)

	var session Session

	_ = row.Scan(&session.SessionId, &session.UserId, &session.CreatedAt)

	return &session, nil

}

func (r *Repository) LastSessionDate(username string) (time.Time, error) {

	query := "SELECT MAX(s.created_at) FROM session s WHERE s.user_id = (SELECT u.user_id FROM users u WHERE u.username ILIKE '%' || $1 || '%');"

	row := r.db.QueryRowContext(context.Background(), query, username)

	var lastSessionCreatedAt time.Time

	if row == nil {

		logrus.Tracef("No last session found")
		// returns 0001-01-01 00:00:00 +0000 UTC
		return lastSessionCreatedAt, nil

	}

	if row.Err() != nil {

		return lastSessionCreatedAt, fmt.Errorf("error getting last session %v", row.Err().Error())

	}

	_ = row.Scan(&lastSessionCreatedAt)

	logrus.Tracef("Last session found: %s", lastSessionCreatedAt.String())

	return lastSessionCreatedAt, nil

}

func (r *Repository) LatestGists(username string) ([]*GistSummary, error) {

	query := "SELECT rgu.gist_id, u.username, g.gist_file_title, r.routine_id FROM routine_gist_user rgu LEFT JOIN users u ON rgu.user_id = u.user_id LEFT JOIN gists g on rgu.gist_id = g.gist_id LEFT JOIN routine r ON rgu.routine_id = r.routine_id WHERE r.created_at > $1 AND u.username ILIKE '%' || $2 || '%';"

	lastSessionDate, err := r.LastSessionDate(username)

	if err != nil {

		return nil, fmt.Errorf("%s", err)

	}

	logrus.Infof("Getting all latests gists for %s since %v", username, lastSessionDate)

	rows, err := r.db.QueryContext(context.Background(), query, lastSessionDate, username)

	if rows == nil {

		return nil, fmt.Errorf("No gists yet\n")

	}

	var gists []*GistSummary

	for rows.Next() {

		g := &GistSummary{}
		err := rows.Scan(&g.GistUUID, &g.GistOwner, &g.Filename, &g.RawUrl)

		if err != nil {
			return nil, fmt.Errorf("Error latest gists query response %v", err)
		}
		gists = append(gists, g)

	}

	err = rows.Close()

	if err != nil {
		return nil, fmt.Errorf("Error closing pgsql rows, %v", err)
	}

	logrus.Tracef("Succesfully got all: %d latest gists", len(gists))

	return gists, nil

}

func (r *Repository) Routine() error {

	allUsersBeingTracked, err := r.GetAllUsers()

	logrus.Info("Initializing Routine")

	if err != nil {

		return fmt.Errorf("failed to get all users %v", err)

	}

	newRoutine, err := r.NewRoutine()

	logrus.Tracef("Routine %d started", newRoutine.Id)

	logrus.Tracef("Iterating over tracked users")

	for _, val := range allUsersBeingTracked {

		currentUserGists, err := NewGistApiRequest(val.Login)

		if err != nil {

			return fmt.Errorf("failed to GET from github's API %v", err)

		}

		logrus.Tracef("Iterating over user %s's %d gists", val.Login, len(*currentUserGists))

		for _, gists := range *currentUserGists {

			logrus.Tracef("Testing if Gist is in DB")

			isGistAlreadyIn, err := r.IsGistInDb(gists.GistUUID)

			logrus.Tracef("Is gist already in: %v", isGistAlreadyIn)

			if err != nil {

				return fmt.Errorf("Failed to check if gist is in db %v", err)

			}

			if isGistAlreadyIn != true {

				logrus.Tracef("Gist %s not in", gists.GistUUID)

				err := r.InsertGistPgAndPipe(&gists)

				if err != nil {

					return fmt.Errorf("Failed to insert gist into pg or/and pipedrive %v", err)

				}

				logrus.Tracef("Populating the routine_gists table")

				query := "INSERT INTO routine_gist_user (routine_id, gist_id, user_id) VALUES ($1, $2, $3);"

				row := r.db.QueryRowContext(context.Background(), query, &newRoutine.Id, &gists.GistUUID, &val.Id)

				if row.Err() != nil {

					return fmt.Errorf("Failed to create new routine_gist")

				}

			} else {

				logrus.Tracef("Gist %s already in", gists.GistUUID)

			}

		}

	}

	return nil

}

func (r *Repository) IsGistInDb(id string) (bool, error) {

	log.Printf("Checking if gist is in DB:\n")

	query := "select CASE when EXISTS( SELECT 1 FROM gists g WHERE g.gist_id ILIKE '%' || $1 || '%' ) then true else false end;"

	row := r.db.QueryRowContext(context.Background(), query, id)

	if row == nil {

		logrus.Warnf("Error getting query row context\n")
		return false, nil
	}

	var isGistIn bool
	err := row.Scan(&isGistIn)

	if err != nil {
		logrus.Warnf("Error scanning row, %v", err)
		return false, err
	}

	logrus.Trace("Successfully checked if gist is in Db")

	return isGistIn, nil

}
