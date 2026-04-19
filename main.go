package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Response struct {
	Secret string `json:"secret"`
}

var rateLimit = make(map[string]time.Time)
var mu sync.Mutex

func generateSecret(w http.ResponseWriter, r *http.Request) {
	// Rate limiting: 10 requests per minute per IP
	ip := r.RemoteAddr
	mu.Lock()
	last := rateLimit[ip]
	if time.Since(last) < time.Minute && rateLimit[ip] != (time.Time{}) {
		// Check if more than 10, but simple count not implemented, just time based
		if time.Since(last) < time.Minute/10 { // rough
			mu.Unlock()
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
	}
	rateLimit[ip] = time.Now()
	mu.Unlock()

	// Generate 32 bytes random key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	secret := base64.StdEncoding.EncodeToString(key)

	resp := Response{Secret: secret}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/generate-secret", generateSecret)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
