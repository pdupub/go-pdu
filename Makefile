# 项目名称
APP_NAME := pdu

# Go 编译器
GO := go

# 构建目录
BUILD_DIR := build

# 源文件目录
SRC_DIR := ./cmd/pdu

# 默认目标
all: build

# 构建目标
build: $(SRC_DIR)/main.go
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

# 清理目标
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)

# 安装依赖
deps:
	@echo "Installing dependencies..."
	@$(GO) mod tidy

# 运行目标
run: build
	@echo "Running $(APP_NAME) with arguments $(ARGS)..."
	@./$(BUILD_DIR)/$(APP_NAME) $(ARGS)

# 测试目标
test:
	@echo "Running tests..."
	@$(GO) test ./...

# 帮助信息
help:
	@echo "Usage:"
	@echo "  make         Build the project"
	@echo "  make build   Build the project"
	@echo "  make clean   Clean the build directory"
	@echo "  make deps    Install dependencies"
	@echo "  make run     Build and run the project with arguments"
	@echo "  make test    Run tests"
	@echo "  make help    Show this help message"

.PHONY: all build clean deps run test help
