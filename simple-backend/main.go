package main

import (
	"fmt"
	"net/http"
	"log"
	"strconv"
	"sync"
)

var payloads map[int][]byte
var once sync.Once

func main() {
	generatePayloads()

	http.HandleFunc("/data/v1/nocachingheaders/", func(w http.ResponseWriter, r *http.Request) {
		PublicDataHandler(w, r)
	})


	http.HandleFunc("/data/v1/info", func(w http.ResponseWriter, r *http.Request) {
		InfoHandler(w, r)
	})

	http.HandleFunc("/data/v1/queryresource", func(w http.ResponseWriter, r *http.Request) {
		QueryResourceHandler(w, r)
	})


	// Start server
	fmt.Println("Server is listening on port 8084...")
	http.ListenAndServe(":8084", nil)
}

func QueryResourceHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Query handler invoked")

	queryParams := r.URL.Query()
	param1 := queryParams.Get("param1")
	param2 := queryParams.Get("param2")
	param3 := queryParams.Get("param3")

	if param1 == "" || param2 == "" || param3 == "" {
		http.Error(w, "Missing query parameters", http.StatusBadRequest)
		return
	}

	response := fmt.Sprintf("param1: %s, param2: %s, param3: %s", param1, param2, param3)

	// Write response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func PublicDataHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("DataHandler invoked")

	id, err := strconv.Atoi(r.URL.Path[len("/data/v1/nocachingheaders/"):])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	payload, ok := payloads[id]
	if !ok {
		http.Error(w, "ID not found", http.StatusNotFound)
		return
	}

	w.Write(payload)
	additionalContent := []byte(" public data for ID " + strconv.Itoa(id))
	w.Write(additionalContent)
}


func InfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("InfoHandler invoked")
	w.Write([]byte("This is general information"))
}

func generatePayloads() {
	once.Do(func() {
		payloads = make(map[int][]byte)
		for i := 1; i <= 100; i++ {
			payload := make([]byte, 1024)
			for j := range payload {
				payload[j] = 'x'
			}
			payloads[i] = payload
		}
	})
}
