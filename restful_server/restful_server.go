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

func getUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("GET /users")
	log.Println(db)
	rows, err := db.Query(`SELECT * FROM users`)
	log.Println("OK1")
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
		//log.Printf("\n id : %d\n name : %s\n email : %s\n created_at : %s\n updated_at : %s\n", user.ID, user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
		//log.Println(user)
		users = append(users, user)
	}
	data, jsonErr := json.Marshal(users)
	if jsonErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	log.Println("OK2")
}

func getUserByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("GET /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)
}

func addUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("POST /users")

	bufbody := new(bytes.Buffer)
	bufbody.ReadFrom(r.Body)
	body := bufbody.String()
	log.Println(body)

	var user User
	if err := json.Unmarshal([]byte(body), &user); err != nil {
		log.Fatal(err)
	}
	log.Println(user)

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

	data, jsonErr := json.Marshal(user)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func updateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("PUT /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)
}

func deleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Println("DELETE /users/:id")
	userID := chi.URLParam(r, "userID")
	log.Println(userID)
	db.Exec(`DELETE FROM users WHERE id=$1`, userID)
}

func main() {
	// host=[docker image name] port=[docker port]
	db, dbOpenErr := sql.Open("postgres", "host=db port=5432 user=pq_user password=password dbname=app_db sslmode=disable")
	if dbOpenErr != nil {
		log.Fatal(dbOpenErr)
	}

	_, dbExecErr := db.Exec(`INSERT INTO users (name, email, created_at, updated_at) VALUES ($1, $2, $3, $4);`, "testname", "hoge@example.com", "2019-05-01T02:34:56.789012345+09:00", "2019-05-01T02:34:56.789012345+09:00")
	if dbExecErr != nil {
		log.Fatal(dbExecErr)
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
	err := http.ListenAndServe(":8080", r) // set port number
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
