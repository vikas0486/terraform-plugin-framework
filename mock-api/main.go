package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type Keystore struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var db = map[string]Keystore{}

func createKeystore(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", 405)
		return
	}

	k := Keystore{
		ID:   "ks-001",
		Name: "payment-keystore",
	}

	db[k.ID] = k

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(k)
}

func getKeystore(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/keystores/")

	k, ok := db[id]

	if !ok {

		http.Error(w, "Not Found", 404)

		return
	}

	json.NewEncoder(w).Encode(k)
}

func main() {

	http.HandleFunc("/keystores", createKeystore)

	http.HandleFunc("/keystores/", getKeystore)

	log.Println("Mock API running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
