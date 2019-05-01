package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// SimpleMessage return message sent as json
type SimpleMessage struct {
	Message string `json:"message"`
}

func respond(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // parse options; nothing by default
	fmt.Println(r.Form)
	fmt.Println("path: \"", r.URL.Path, "\"")
	fmt.Println("scheme: \"", r.URL.Scheme, "\"")
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
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

func main() {
	http.HandleFunc("/", respond)            // set routing
	err := http.ListenAndServe(":8080", nil) // set port number
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
