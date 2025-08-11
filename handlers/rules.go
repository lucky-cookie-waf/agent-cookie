package handlers

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    
    "github.com/lucky-cookie-waf/agent-cookie/services"
)

func RulesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        getCurrentRulesHandler(w, r)
    case "POST":
        addRuleHandler(w, r)
    case "DELETE":
        deleteRuleHandler(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

// 현재 커스텀 룰들 조회
func getCurrentRulesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    log.Println("📋 Reading current custom rules...")
    
    rules, err := services.GetCustomRules()
    if err != nil {
        log.Printf("❌ Failed to read rules: %v", err)
        http.Error(w, "Failed to read rules", http.StatusInternalServerError)
        return
    }
    
    response := map[string]interface{}{
        "rules":   rules,
        "message": "🍪 Custom rules retrieved by Agent Cookie!",
    }
    
    json.NewEncoder(w).Encode(response)
}

// 룰 추가
func addRuleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }
    
    ruleContent := string(body)
    log.Printf("🍪 Adding new rule... (%d bytes)", len(ruleContent))
    
    err = services.AddRule(ruleContent)
    if err != nil {
        log.Printf("❌ Failed to add rule: %v", err)
        http.Error(w, "Failed to add rule", http.StatusInternalServerError)
        return
    }
    
    log.Println("✅ Rule added successfully!")
    
    response := map[string]string{
        "status":  "success",
        "message": "🍪 Rule added by Agent Cookie!",
    }
    json.NewEncoder(w).Encode(response)
}

// 룰 삭제
func deleteRuleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    ruleID := r.URL.Query().Get("id")
    if ruleID == "" {
        http.Error(w, "Rule ID is required", http.StatusBadRequest)
        return
    }
    
    log.Printf("🗑️ Removing rule: %s", ruleID)
    
    err := services.RemoveRule(ruleID)
    if err != nil {
        log.Printf("❌ Failed to remove rule: %v", err)
        http.Error(w, "Failed to remove rule", http.StatusInternalServerError)
        return
    }
    
    log.Println("✅ Rule removed successfully!")
    
    response := map[string]string{
        "status":  "success", 
        "message": "🍪 Rule removed by Agent Cookie!",
    }
    json.NewEncoder(w).Encode(response)
}