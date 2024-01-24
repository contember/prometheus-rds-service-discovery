package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Target struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

var config ScrapingConfig

// @see https://prometheus.io/docs/prometheus/latest/http_sd/#writing-http-service-discovery
func getClusters(w http.ResponseWriter, r *http.Request) {
	var targets = scrapeClusters(config)
	jsonResponse, err := json.Marshal(targets)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(jsonResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello, service running!")
}

func main() {
	var envConfig, err = buildFromEnvironments()

	config = envConfig

	if err != nil {
		fmt.Println("Error during loading config from environment variables", err)
		return
	}

	fmt.Println("RDS Instance discovery running")
	fmt.Println("Flowing cluster ready for scraping:", config.Targets)

	// Launch HTTP server
	http.HandleFunc("/discovery", getClusters)
	http.HandleFunc("/.health", getHealth)

	// todo @sitole: make possible to run with $PORT
	error := http.ListenAndServe(":3333", nil)

	if error != nil {
		log.Fatal(error)
	}
}
