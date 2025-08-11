# Agent Cookie Build Configuration 🍪

# 변수 설정
BINARY_NAME=agent-cookie
VERSION=1.0.0
BUILD_DIR=build

# 기본 타겟
.PHONY: build clean test build-all help

# 로컬 빌드
build:
	@echo "🍪 Building Agent Cookie..."
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME)

# 모든 플랫폼 빌드
build-all: clean
	@echo "🍪 Building Agent Cookie for all platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64
	
	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	
	@echo "✅ Build complete! Check $(BUILD_DIR)/ directory"

# 테스트
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# 정리
clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)

# 개발용 실행
dev:
	@echo "🚀 Running in development mode..."
	go run main.go --config config.yaml

# 도움말
help:
	@echo "🍪 Agent Cookie Build Commands:"
	@echo "  make build     - Build for current platform"
	@echo "  make build-all - Build for all platforms"
	@echo "  make test      - Run tests"
	@echo "  make clean     - Clean build files"
	@echo "  make dev       - Run in development mode"
	@echo "  make help      - Show this help"