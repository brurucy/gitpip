package main

import (
	"database/sql"
	"fmt"
	"gitpip/pkg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	logLevel        = "debug"
	httpServicePort = 8080
)

func main() {

	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.Info("Starting program")

	logrus.Trace("Attempting to connect to postgres")
	dbConn, err := sql.Open("postgres", os.Getenv("POSTGRES_CONNECTION_STRING"))
	logrus.Info(os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {

		logrus.Fatalf("%v", err)

	}

	defer dbConn.Close()

	repository := pkg.NewRepository(dbConn)
	handler := pkg.NewHandler(repository)
	router := mux.NewRouter()

	logrus.Trace("Registering routes")
	handler.RegisterRoutes(router)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpServicePort),
		Handler: router,
	}
	logrus.Infof("Serving HTTP on port: %v", httpServicePort)

	go func() {
		err = httpSrv.ListenAndServe()
		if err != nil {
			log.Fatalf("Could not start http server")
		}

	}()

	ticker := time.NewTicker(5 * time.Minute)
	done := make(chan bool)

	logrus.Info("Starting routine")
	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:

			repository.Routine()
		}
	}

}
