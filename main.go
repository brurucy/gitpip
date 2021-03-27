package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gitpip/pkg"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	logLevel        = "debug"
	httpServicePort = 8080
)

func main() {

	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))

	if err != nil {

		fmt.Errorf("%v", err)

	}

	defer dbConn.Close()

	repository := pkg.NewRepository(dbConn)
	handler := pkg.NewHandler(repository)
	router := mux.NewRouter()

	handler.RegisterRoutes(router)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpServicePort),
		Handler: router,
	}
	log.Println("Serving HTTP on port: ", httpServicePort)

	go func() {
		err = httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("Could not start http server")
		}

	}()

	ticker := time.NewTicker(180 * time.Minute)
	done := make(chan bool)

	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:

			repository.Routine()
		}
	}

}
