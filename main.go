package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// Define a struct to parse and hold the incoming payload
type WebhookPayload struct {
	Object string        `json:"object"`
	Entry  []interface{} `json:"entry"`
}
type ResponseGet struct {
	Message string `json:"message"`
}

// Handle the GET request for token verification and POST for webhook events
func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Retrieve the secure token from the environment
	verifyToken := os.Getenv("SECURE_VERIFY_TOKEN")
	if verifyToken == "" {
		log.Fatal("SECURE_VERIFY_TOKEN not set in environment")
	}
	switch r.Method {
	case "GET":
		token := r.URL.Query().Get("verify_token")
		challenge := r.URL.Query().Get("challenge")

		if token == verifyToken {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(challenge))
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseGet{Message: "Token verified"})

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ResponseGet{Message: "Invalid verification token"})
		}
	case "POST":
		var payload WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Log or process the payload data
		log.Printf("Received event: %+v\n", payload)

		// Encode the parsed payload back as JSON and send it in the response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/webhook", verifyWebhook)
	log.Println("Starting server on port 8088...")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
