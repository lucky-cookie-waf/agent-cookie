package handlers

import (
    "encoding/json"
    "net/http"
    "log"
    
    "github.com/lucky-cookie-waf/agent-cookie/services"
)

func LogsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        getLogsHandler(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getLogsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    log.Println("🕵️ Collecting new ModSecurity logs...")
    
    // 새로운 로그만 수집 (raw strings)
    newLogs, err := services.CollectNewLogs()
    if err != nil {
        log.Printf("❌ Failed to collect logs: %v", err)
        http.Error(w, "Failed to collect logs", http.StatusInternalServerError)
        return
    }
    
    log.Printf("✅ Found %d new log entries", len(newLogs))
    
    // 중앙 서버로 raw 로그 전송
    response := map[string]interface{}{
        "logs":      newLogs,
        "count":     len(newLogs),
        "timestamp": time.Now().Format(time.RFC3339),
    }
    
    json.NewEncoder(w).Encode(response)
}