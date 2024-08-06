package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func NewResourceEndpoint(route string, response string) {
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

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	})
}

func NewResourceArrayEndpoint(route string, response []string) {
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

func SerializeMapToJSONArray[T any](data map[string]T) []string {
	result := make([]string, 0)

	for _, v := range data {
		jsonBytes, err := json.Marshal(v)
		if err == nil {
			result = append(result, string(jsonBytes))
		}
	}

	return result
}

func SerializeStruct[T any](data T) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}
