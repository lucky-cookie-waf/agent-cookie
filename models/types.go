package models

import "time"

// 로그 엔트리 (단순화 - raw 로그만 전송)
type LogEntry struct {
    Content   string    `json:"content"`
    Timestamp time.Time `json:"timestamp"`
}

// 로그 응답
type LogResponse struct {
    Logs      []string `json:"logs"`
    Count     int      `json:"count"`
    Timestamp string   `json:"timestamp"`
}

// 룰 요청/응답
type RuleRequest struct {
    Content string `json:"content"`
    ID      string `json:"id,omitempty"`
}

type RuleResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

// 헬스체크 응답
type HealthResponse struct {
    Status    string `json:"status"`
    Version   string `json:"version"`
    Timestamp string `json:"timestamp"`
    Message   string `json:"message"`
}

// 에이전트 설정
type Config struct {
    ListenAddr    string `yaml:"listen_addr"`
    CentralServer string `yaml:"central_server"`
    LogPath       string `yaml:"log_path"`
    Interval      string `yaml:"interval"`
    Debug         bool   `yaml:"debug"`
}