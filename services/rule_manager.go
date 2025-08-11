package services

import (
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

// ModSecurity 설정 경로들
var modSecurityPaths = []string{
    "/etc/modsecurity/rules/",
    "/etc/apache2/modsecurity/rules/", 
    "/usr/local/modsecurity/rules/",
    "/etc/httpd/modsecurity/rules/",
}

// 커스텀 룰 파일명
const customRulesFile = "agent-cookie-custom.conf"

// 룰 추가
func AddRule(ruleContent string) error {
    ruleDir, err := findModSecurityRulesDir()
    if err != nil {
        return fmt.Errorf("🚫 ModSecurity rules directory not found: %v", err)
    }
    
    filePath := filepath.Join(ruleDir, customRulesFile)
    
    // 기존 파일이 있으면 append, 없으면 새로 생성
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("🚫 Failed to open rules file: %v", err)
    }
    defer file.Close()
    
    // 룰 추가 (줄바꿈 포함)
    _, err = file.WriteString(fmt.Sprintf("\n# Added by Agent Cookie\n%s\n", ruleContent))
    if err != nil {
        return fmt.Errorf("🚫 Failed to write rule: %v", err)
    }
    
    return nil
}

// 룰 삭제 (ID로)
func RemoveRule(ruleID string) error {
    ruleDir, err := findModSecurityRulesDir()
    if err != nil {
        return fmt.Errorf("🚫 ModSecurity rules directory not found: %v", err)
    }
    
    filePath := filepath.Join(ruleDir, customRulesFile)
    
    // 파일 읽기
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("🚫 Failed to read rules file: %v", err)
    }
    
    // 룰 제거
    newContent := removeRuleFromContent(string(content), ruleID)
    
    // 파일 다시 쓰기
    err = ioutil.WriteFile(filePath, []byte(newContent), 0644)
    if err != nil {
        return fmt.Errorf("🚫 Failed to update rules file: %v", err)
    }
    
    return nil
}

// 모든 커스텀 룰 삭제
func ClearCustomRules() error {
    ruleDir, err := findModSecurityRulesDir()
    if err != nil {
        return fmt.Errorf("🚫 ModSecurity rules directory not found: %v", err)
    }
    
    filePath := filepath.Join(ruleDir, customRulesFile)
    
    // 파일 삭제
    if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("🚫 Failed to remove rules file: %v", err)
    }
    
    return nil
}

// 현재 커스텀 룰들 조회
func GetCustomRules() (string, error) {
    ruleDir, err := findModSecurityRulesDir()
    if err != nil {
        return "", fmt.Errorf("🚫 ModSecurity rules directory not found: %v", err)
    }
    
    filePath := filepath.Join(ruleDir, customRulesFile)
    
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return "", nil // 파일이 없으면 빈 문자열 반환
        }
        return "", fmt.Errorf("🚫 Failed to read rules file: %v", err)
    }
    
    return string(content), nil
}

// ModSecurity 룰 디렉토리 찾기
func findModSecurityRulesDir() (string, error) {
    for _, path := range modSecurityPaths {
        if info, err := os.Stat(path); err == nil && info.IsDir() {
            return path, nil
        }
    }
    return "", fmt.Errorf("no ModSecurity rules directory found")
}

// 룰 ID로 해당 룰 제거
func removeRuleFromContent(content, ruleID string) string {
    lines := strings.Split(content, "\n")
    var newLines []string
    skipUntilNext := false
    
    for _, line := range lines {
        // 룰 ID가 포함된 라인 찾기
        if strings.Contains(line, fmt.Sprintf(`id "%s"`, ruleID)) {
            skipUntilNext = true
            continue
        }
        
        // 다음 룰이 시작되면 스킵 해제
        if skipUntilNext && (strings.HasPrefix(strings.TrimSpace(line), "SecRule") || 
                           strings.HasPrefix(strings.TrimSpace(line), "#")) {
            if !strings.Contains(line, ruleID) {
                skipUntilNext = false
                newLines = append(newLines, line)
            }
        } else if !skipUntilNext {
            newLines = append(newLines, line)
        }
    }
    
    return strings.Join(newLines, "\n")
}