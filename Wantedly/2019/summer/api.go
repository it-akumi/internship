package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type simpleResponse struct {
	Message string `json:"message"`
}

var (
	port string
	db   *sql.DB
)

const (
	localPort           string = "8080"
	localDataSourceName string = "postgres://wantedly@db/webapp?sslmode=disable"
)

// init set global variables to different values between local and Heroku
func init() {
	port = os.Getenv("PORT") // Use $PORT in Heroku
	if port == "" {
		port = localPort
	}

	dsn := os.Getenv("DATABASE_URL") // Use $DATABASE_URL in Heroku
	if dsn == "" {
		dsn = localDataSourceName
	}
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(simpleResponse{Message: "Hello World!!"})
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func usersGetHandler(w http.ResponseWriter, r *http.Request) {
	users, err := getUsers(db)
	if users == nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(users)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func usersPostHandler(w http.ResponseWriter, r *http.Request) {
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := make([]byte, length)
	length, _ = r.Body.Read(body) // ignore EOF
	u := new(user)
	if err := json.Unmarshal(body[:length], u); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	createdUser, err := createUser(db, u)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(createdUser)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(b))
}

func userGetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user, err := getUser(db, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func userPutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := make([]byte, length)
	length, _ = r.Body.Read(body) // ignore EOF
	u := new(user)
	if err := json.Unmarshal(body[:length], u); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updatedUser, err := updateUser(db, vars["id"], u)
	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(updatedUser)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(b))
}

func userDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := deleteUser(db, vars["id"])
	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		return
	case err != nil:
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/users", usersGetHandler).Methods("GET")
	r.HandleFunc("/users", usersPostHandler).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", userGetHandler).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", userPutHandler).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", userDeleteHandler).Methods("Delete")

	http.Handle("/", r)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
