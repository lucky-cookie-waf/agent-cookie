package services

import (
    "bufio"
    "crypto/md5"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
    "time"
)

type LogState struct {
    LastPosition int64  `json:"last_position"`
    LastChecksum string `json:"last_checksum"`
    LastModTime  int64  `json:"last_mod_time"`
}

// 새로운 로그만 수집
func CollectNewLogs() ([]string, error) {
    logPaths := []string{
        "/var/log/modsecurity/audit.log",
        "/var/log/httpd/modsecurity_audit.log", 
        "/var/log/apache2/modsecurity_audit.log",
    }
    
    var newLogs []string
    
    for _, logPath := range logPaths {
        if _, err := os.Stat(logPath); err == nil {
            logs, err := readNewLogsFromFile(logPath)
            if err != nil {
                continue
            }
            newLogs = append(newLogs, logs...)
        }
    }
    
    return newLogs, nil
}

func readNewLogsFromFile(logPath string) ([]string, error) {
    // 상태 파일 경로
    stateFile := getStateFilePath(logPath)
    
    // 이전 상태 로드
    lastState, err := loadState(stateFile)
    if err != nil {
        // 상태 파일이 없으면 파일 끝부터 시작
        lastState = &LogState{LastPosition: getFileSize(logPath)}
    }
    
    // 파일 열기
    file, err := os.Open(logPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    // 파일 정보 확인
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }
    
    currentSize := fileInfo.Size()
    
    // 파일이 작아졌으면 로테이션된 것 (처음부터 읽기)
    if currentSize < lastState.LastPosition {
        lastState.LastPosition = 0
    }
    
    // 마지막 위치부터 읽기
    _, err = file.Seek(lastState.LastPosition, 0)
    if err != nil {
        return nil, err
    }
    
    var newLines []string
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.TrimSpace(line) != "" {
            newLines = append(newLines, line)
        }
    }
    
    // 새로운 상태 저장
    newState := &LogState{
        LastPosition: currentSize,
        LastChecksum: calculateChecksum(newLines),
        LastModTime:  fileInfo.ModTime().Unix(),
    }
    
    saveState(stateFile, newState)
    
    return newLines, scanner.Err()
}

// 상태 파일 경로 생성
func getStateFilePath(logPath string) string {
    hasher := md5.New()
    hasher.Write([]byte(logPath))
    hash := fmt.Sprintf("%x", hasher.Sum(nil))
    return filepath.Join("/tmp", fmt.Sprintf("agent-cookie-state-%s.json", hash[:8]))
}

// 파일 크기 가져오기
func getFileSize(filePath string) int64 {
    if info, err := os.Stat(filePath); err == nil {
        return info.Size()
    }
    return 0
}

// 상태 로드
func loadState(stateFile string) (*LogState, error) {
    data, err := ioutil.ReadFile(stateFile)
    if err != nil {
        return nil, err
    }
    
    // 간단한 형태로 파싱 (실제로는 JSON 사용)
    lines := strings.Split(string(data), "\n")
    if len(lines) >= 3 {
        position := parseInt64(lines[0])
        checksum := lines[1]
        modTime := parseInt64(lines[2])
        
        return &LogState{
            LastPosition: position,
            LastChecksum: checksum,
            LastModTime:  modTime,
        }, nil
    }
    
    return nil, fmt.Errorf("invalid state file")
}

// 상태 저장  
func saveState(stateFile string, state *LogState) error {
    content := fmt.Sprintf("%d\n%s\n%d\n", 
        state.LastPosition, 
        state.LastChecksum, 
        state.LastModTime)
    
    return ioutil.WriteFile(stateFile, []byte(content), 0644)
}

// 체크섬 계산
func calculateChecksum(lines []string) string {
    hasher := md5.New()
    for _, line := range lines {
        hasher.Write([]byte(line))
    }
    return fmt.Sprintf("%x", hasher.Sum(nil))
}

// 문자열을 int64로 변환
func parseInt64(s string) int64 {
    if val, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64); err == nil {
        return val
    }
    return 0
}