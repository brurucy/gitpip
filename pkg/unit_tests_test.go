package pkg

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

func TestApiUnmarshal(t *testing.T) {

	resp, err := NewGistApiRequest("brurucy")

	if err != nil {

		t.Fatalf("%v", err)

	}

	fmt.Println(resp)

}

func TestPipeApi(t *testing.T) {

	req := &PipeAddActivityRequest{
		Subject: "Testing-again",
		Note:    "Gist",
		Done:    false,
	}

	err := NewPipedriveActivity(req)

	if err != nil {

		t.Fatalf("%v", err)

	}

}

func TestGistDownloader(t *testing.T) {

	gistUrl := "https://gist.githubusercontent.com/robotsnowfall/e823db4efe48513088bc74f08670d78c/raw/614f4fe4713a7fdef59b469508e6960286b94d92/gistfile1.txt"

	_, err := GistTextDownloader(gistUrl)

	if err != nil {

		t.Fatalf("%v", err)

	}

}

func TestIsUserBeingTracked(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	vals, err := repository.IsUserBeingTracked("brurucy")

	if err != nil {
		t.Error(err)
	}

	if vals != false {
		t.Errorf("Unexpected query result: %v", vals)
	}

}

func TestIfIsGistAdded(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	vals, err := repository.IsGistInDb("51b13376431d67d20548d9e008c465f3")

	if err != nil {
		t.Error(err)
	}

	if vals != false {
		t.Errorf("Unexpected query result: %v", vals)
	}

}

func TestInsertUser(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	newUser := &GistOwner{
		Login: "brurucys",
		Id:    929292,
	}

	err = repository.InsertUser(newUser)

	if err != nil {
		t.Error(err)
	}

}

func TestInsertGist(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	gist, err := NewGistApiRequest("noah")

	if err != nil {
		t.Error(err)
	}

	err = repository.InsertGistPgAndPipe(&(*gist)[0])

	if err != nil {
		t.Error(err)
	}

}

func TestGetAll(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	users, err := repository.GetAllUsers()

	for _, val := range users {
		fmt.Println(*val)
	}

}

func TestNewRoutine(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	sess, err := repository.NewRoutine()

	if err != nil {

		t.Error(err)

	}

	fmt.Println(sess)

}

func TestRoutine(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	me := &GistOwner{
		Login: "brurucy",
		Id:    4424,
	}

	randomPerson := &GistOwner{
		Login: "lewisgaul",
		Id:    16408073,
	}

	repository.InsertUser(me)
	repository.InsertUser(randomPerson)

	err = repository.Routine()

	fmt.Println(err)

}

func TestSession(t *testing.T) {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		t.Error(err)

	}

	defer dbConn.Close()

	repository := NewRepository(dbConn)

	_, err = repository.NewSession("brurucy")

	if err != nil {

		t.Fatal()

	}

}
