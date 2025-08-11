package handlers

import (
    "encoding/json"
    "net/http"
    "time"
)

type HealthResponse struct {
    Status    string `json:"status"`
    Version   string `json:"version"`
    Timestamp string `json:"timestamp"`
    Message   string `json:"message"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    response := HealthResponse{
        Status:    "healthy",
        Version:   "1.0.0",
        Timestamp: time.Now().Format(time.RFC3339),
        Message:   "🍪 Agent Cookie is watching your firewall!",
    }
    
    json.NewEncoder(w).Encode(response)
}