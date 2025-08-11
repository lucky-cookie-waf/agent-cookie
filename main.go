package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    
    "github.com/gorilla/mux"
    "github.com/lucky-cookie-waf/agent-cookie/config"
    "github.com/lucky-cookie-waf/agent-cookie/handlers"
)

var version = "1.0.0"

func main() {
    // 커맨드라인 플래그
    configPath := flag.String("config", "config.yaml", "Configuration file path")
    showVersion := flag.Bool("version", false, "Show version")
    flag.Parse()
    
    if *showVersion {
        log.Printf("🍪 Agent Cookie v%s", version)
        os.Exit(0)
    }
    
    log.Println("🍪 Agent Cookie is starting...")
    
    // 설정 로드
    cfg, err := config.LoadFromFile(*configPath)
    if err != nil {
        log.Printf("⚠️  Config file not found, using defaults: %v", err)
        cfg = config.GetDefault()
    }
    
    if cfg.Debug {
        log.Println("🐛 Debug mode enabled")
    }
    
    // 라우터 설정
    r := mux.NewRouter()
    
    // API 엔드포인트
    r.HandleFunc("/health", handlers.HealthHandler).Methods("GET")
    r.HandleFunc("/logs", handlers.LogsHandler).Methods("GET")
    r.HandleFunc("/rules", handlers.RulesHandler).Methods("GET", "POST", "DELETE")
    
    // 404 핸들러
    r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"error": "🍪 Agent Cookie: Endpoint not found"}`))
    })
    
    log.Printf("🕵️ Agent Cookie is watching on %s", cfg.ListenAddr)
    log.Printf("🌐 Central server: %s", cfg.CentralServer)
    log.Fatal(http.ListenAndServe(cfg.ListenAddr, r))
}