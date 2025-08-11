#!/bin/bash
set -e

# Agent Cookie 설치 스크립트 🍪
VERSION="latest"
REPO="lucky-cookie-waf/agent-cookie"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.agent-cookie"

echo "🍪 Installing Agent Cookie..."

# 시스템 정보 감지
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case $OS in
        linux*) OS="linux" ;;
        darwin*) OS="darwin" ;;
        mingw*|msys*|cygwin*) OS="windows" ;;
        *) echo "❌ Unsupported OS: $OS" && exit 1 ;;
    esac
    
    case $ARCH in
        x86_64|amd64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
    esac
    
    echo "📋 Detected platform: $OS-$ARCH"
}

# 바이너리 다운로드
download_binary() {
    BINARY_NAME="agent-cookie"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="agent-cookie.exe"
    fi
    
    DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/agent-cookie-$OS-$ARCH"
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_URL="$DOWNLOAD_URL.exe"
    fi
    
    echo "📥 Downloading from: $DOWNLOAD_URL"
    
    # 임시 파일로 다운로드
    TMP_FILE="/tmp/agent-cookie-$OS-$ARCH"
    curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"
    
    # 실행 권한 부여
    chmod +x "$TMP_FILE"
    
    # 설치 디렉토리로 이동
    if command -v sudo >/dev/null 2>&1; then
        sudo mv "$TMP_FILE" "$INSTALL_DIR/agent-cookie"
        echo "✅ Installed to $INSTALL_DIR/agent-cookie"
    else
        # sudo 없으면 로컬 디렉토리에 설치
        mkdir -p "$HOME/.local/bin"
        mv "$TMP_FILE" "$HOME/.local/bin/agent-cookie"
        echo "✅ Installed to $HOME/.local/bin/agent-cookie"
        
        # PATH에 추가 안내
        if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
            echo "⚠️  Add $HOME/.local/bin to your PATH:"
            echo "   echo 'export PATH=\$PATH:$HOME/.local/bin' >> ~/.bashrc"
            echo "   source ~/.bashrc"
        fi
    fi
}

# 설정 파일 생성
# ❗️ 웹서버 도메인 수정 필요
create_config() {
    echo "⚙️  Creating configuration..."
    mkdir -p "$CONFIG_DIR"
    
    cat > "$CONFIG_DIR/config.yaml" <<EOF
# Agent Cookie Configuration 🍪
listen_addr: "127.0.0.1:8080"
central_server: "${CENTRAL_SERVER:-https://api.cookie-security.com}"
log_path: "${LOG_PATH:-/var/log/modsecurity}"
interval: "${INTERVAL:-30s}"
debug: false
EOF
    
    echo "📁 Config created at: $CONFIG_DIR/config.yaml"
}

# 에이전트 시작
start_agent() {
    echo "🚀 Starting Agent Cookie..."
    
    # 백그라운드로 실행
    if command -v agent-cookie >/dev/null 2>&1; then
        AGENT_CMD="agent-cookie"
    elif [ -x "$HOME/.local/bin/agent-cookie" ]; then
        AGENT_CMD="$HOME/.local/bin/agent-cookie"
    else
        AGENT_CMD="$INSTALL_DIR/agent-cookie"
    fi
    
    # 이미 실행 중인지 확인
    if pgrep -f "agent-cookie" > /dev/null; then
        echo "⚠️  Agent Cookie is already running"
        return 0
    fi
    
    # 백그라운드 실행
    nohup "$AGENT_CMD" --config "$CONFIG_DIR/config.yaml" > "$CONFIG_DIR/agent.log" 2>&1 &
    AGENT_PID=$!
    echo $AGENT_PID > "$CONFIG_DIR/agent.pid"
    
    # 잠시 대기 후 상태 확인
    sleep 2
    if kill -0 $AGENT_PID 2>/dev/null; then
        echo "✅ Agent Cookie started successfully! (PID: $AGENT_PID)"
        
        # 헬스체크
        if curl -s http://127.0.0.1:8080/health > /dev/null 2>&1; then
            echo "🍪 Agent Cookie is healthy and ready!"
        else
            echo "⚠️  Agent started but health check failed. Check logs: $CONFIG_DIR/agent.log"
        fi
    else
        echo "❌ Failed to start Agent Cookie. Check logs: $CONFIG_DIR/agent.log"
        exit 1
    fi
}

# 설치 완료 메시지
show_completion() {
    echo ""
    echo "🎉 Agent Cookie installation complete!"
    echo ""
    echo "📊 Check status:    curl http://127.0.0.1:8080/health"
    echo "📋 View logs:       tail -f $CONFIG_DIR/agent.log"
    echo "⚙️  Edit config:     nano $CONFIG_DIR/config.yaml"
    echo "🛑 Stop agent:      kill \$(cat $CONFIG_DIR/agent.pid)"
    echo ""
    echo "🔗 Documentation: https://github.com/$REPO"
}

# 메인 실행
main() {
    detect_platform
    download_binary
    create_config
    start_agent
    show_completion
}

# 에러 핸들링
trap 'echo "❌ Installation failed. Check the error above."' ERR

# 실행!
main "$@"