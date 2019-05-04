package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

// User return user data
type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// return user data specified by userID
func findUserByID(userID string, db *sql.DB) bool {
	if rows, err := db.Query(`SELECT * FROM users where id=$1;`, userID); err != nil {
		log.Fatal(err)
	} else {
		if !rows.Next() {
			return false
		}
		return true
	}
	return false
}

// return `User` including parameters given by request body
func parseBody(r *http.Request) User {
	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	log.Println("body : ", body)

	var user User
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		log.Fatal(err)
	}
	return user
}

func getUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("GET /users")

	rows, err := db.Query(`SELECT * FROM users;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if len(users) != 0 {
		data, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("data : ", string(data))
		w.Write(data)
	} else {
		var emptyArray [0]User
		data, err := json.Marshal(emptyArray)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("data : ", string(data))
		w.Write(data)
	}
}

func getUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("GET /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)

	var user User
	if rows, err := db.Query(`SELECT * FROM users where id=$1;`, userID); err != nil {
		log.Fatal(err)
	} else {
		if !rows.Next() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func addUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("POST /users")

	user := parseBody(r)
	if user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	timestamp := time.Now().Format(time.RFC3339Nano)
	if _, err := db.Exec(`INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4);`, user.Name, user.Email, timestamp, timestamp); err != nil {
		log.Fatal(err)
	}

	var userID int
	if rows, err := db.Query(`SELECT currval('users_id_seq');`); err != nil {
		log.Fatal(err)
	} else {
		rows.Next()
		rows.Scan(&userID)
	}

	if rows, err := db.Query(`SELECT * FROM users WHERE id=$1;`, userID); err != nil {
		log.Fatal(err)
	} else {
		rows.Next()
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func updateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("PUT /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)

	if !findUserByID(userID, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := parseBody(r)
	if user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timestamp := time.Now().Format(time.RFC3339Nano)

	if _, err := db.Exec(`UPDATE users SET name=$1, email=$2, updated_at=$3 WHERE id=$4;`, user.Name, user.Email, timestamp, userID); err != nil {
		log.Fatal(err)
	}

	if rows, err := db.Query(`SELECT * FROM users where id=$1;`, userID); err != nil {
		log.Fatal(err)
	} else {
		if !rows.Next() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func deleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("DELETE /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)

	if !findUserByID(userID, db) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := db.Exec(`DELETE FROM users WHERE id=$1;`, userID); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	// host=[docker image name] port=[docker port]
	db, err := sql.Open("postgres", "host=db port=5432 user=pq_user password=password dbname=app_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// set routing
	r := chi.NewRouter()
	r.Route("/users", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			getUsers(w, r, db)
		}) // GET /users
		r.Get("/{userID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			getUserByID(w, r, db)
		}) // GET /users/:id
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			addUser(w, r, db)
		}) // POST /users
		r.Put("/{userID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			updateUser(w, r, db)
		}) // PUT /users/:id
		r.Delete("/{userID:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			deleteUser(w, r, db)
		}) // DELETE /users/:id
	})
	r.HandleFunc("/", respond)
	log.Println("init complete")
	err = http.ListenAndServe(":8080", r) // set port number
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// SimpleMessage return message sent as json
type SimpleMessage struct {
	Message string `json:"message"`
}

func respond(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse options; nothing by default
	log.Println(r.Method)

	log.Println("url: \"", r.URL, "\"")
	log.Println("path: \"", r.URL.Path, "\"")
	log.Println("scheme: \"", r.URL.Scheme, "\"")
	log.Println(r.Form)
	for k, v := range r.Form {
		log.Println("key:", k)
		log.Println("val:", strings.Join(v, ""))
	}

	sm := SimpleMessage{"Hello World!!"}
	data, err := json.Marshal(sm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
