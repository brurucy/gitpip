package pkg

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler struct {
	repository *Repository
}

func NewHandler(r *Repository) *Handler {
	return &Handler{
		repository: r,
	}
}
func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/users", h.GetAllUsers).Methods(http.MethodGet)
	router.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/users/{username}/{id}", h.TrackUser).Methods(http.MethodPost)
	router.HandleFunc("/latestgists/{username}", h.LatestGists).Methods(http.MethodPost)
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, _ *http.Request) {

	logrus.Trace("Got a request to get all tracked users")

	users, err := h.repository.GetAllUsers()

	if err != nil {

		logrus.Errorf("Error fetching all users: %v", err)
		w.WriteHeader(400)
		return

	}

	if len(users) == 0 {

		users = append(users, &GistOwner{})

	}

	json.NewEncoder(w).Encode(users)
	w.WriteHeader(200)

	return
}

func (h *Handler) TrackUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]
	id := vars["id"]

	logrus.Tracef("Got a request to track user %s", username)

	i, err := strconv.ParseInt(id, 10, 32)

	newOwner := &GistOwner{Login: username,
		Id: int32(i)}

	err = h.repository.InsertUser(newOwner)

	if err != nil {

		logrus.Errorf("Bad request %v", err)
		w.WriteHeader(400)

		return

	}

	json.NewEncoder(w).Encode(newOwner)
	w.WriteHeader(200)

	return
}

func (h *Handler) LatestGists(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	logrus.Trace("Got a request to get all latest gists")

	gists, err := h.repository.LatestGists(username)

	if err != nil {

		logrus.Errorf("Error %v", err)
		w.WriteHeader(400)
		return
	}

	_, err = h.repository.NewSession(username)

	if err != nil {

		logrus.Errorf("Error %v", err)
		w.WriteHeader(400)
		return
	}

	if len(gists) == 0 {

		gists = append(gists, &GistSummary{})

	}

	json.NewEncoder(w).Encode(gists)
	w.WriteHeader(200)

	return
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {

	logrus.Trace("Got a health check")

	w.Write([]byte("Alive"))
	w.WriteHeader(200)

}
