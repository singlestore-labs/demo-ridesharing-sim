package api

import (
	"fmt"
	"net/http"
	"strings"
)

func NewResourceEndpoint(route string, response []string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		response := fmt.Sprintf("[%s]", strings.Join(response, ", "))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	})
}

func Serve() {
	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
