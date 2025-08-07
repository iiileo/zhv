# ZHV 构建配置
BINARY_NAME=zhv
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT=$(shell git rev-parse --short HEAD)

# Go 构建参数
GO_BUILD_FLAGS=-ldflags="-s -w -X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"
GO_BUILD_ENV=CGO_ENABLED=0

# 支持的平台
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	freebsd/amd64

.PHONY: all build build-all clean test deps help

# 默认目标
all: build

# 构建当前平台的二进制文件
build:
	@echo "构建 $(BINARY_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	$(GO_BUILD_ENV) go build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) ./

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 安装依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f *.tar.gz *.zip

# 构建所有平台的二进制文件
build-all: clean
	@echo "构建所有平台的二进制文件..."
	@$(foreach platform,$(PLATFORMS),\
		$(call build_platform,$(platform)))

# 构建指定平台的函数
define build_platform
	$(eval GOOS := $(word 1,$(subst /, ,$(1))))
	$(eval GOARCH := $(word 2,$(subst /, ,$(1))))
	$(eval SUFFIX := $(if $(filter windows,$(GOOS)),.exe,))
	$(eval ARCHIVE_EXT := $(if $(filter windows,$(GOOS)),.zip,.tar.gz))
	@echo "构建 $(GOOS)/$(GOARCH)..."
	@GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO_BUILD_ENV) go build $(GO_BUILD_FLAGS) -o $(BINARY_NAME)-$(GOOS)-$(GOARCH)$(SUFFIX) ./
	@if [ "$(GOOS)" = "windows" ]; then \
		zip $(BINARY_NAME)-$(GOOS)-$(GOARCH).zip $(BINARY_NAME)-$(GOOS)-$(GOARCH)$(SUFFIX) README.md; \
	else \
		tar -czf $(BINARY_NAME)-$(GOOS)-$(GOARCH).tar.gz $(BINARY_NAME)-$(GOOS)-$(GOARCH)$(SUFFIX) README.md; \
	fi
endef

# 安装到本地
install: build
	@echo "安装到本地..."
	sudo cp $(BINARY_NAME) /usr/local/bin/

# 卸载
uninstall:
	@echo "从本地卸载..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# 显示版本信息
version:
	@echo "版本: $(VERSION)"
	@echo "构建时间: $(BUILD_TIME)"
	@echo "Git提交: $(GIT_COMMIT)"

# 运行应用（示例）
run: build
	./$(BINARY_NAME) --help

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@which golangci-lint >/dev/null 2>&1 || { echo "请先安装 golangci-lint"; exit 1; }
	golangci-lint run

# 显示帮助信息
help:
	@echo "可用的 make 目标:"
	@echo "  build      - 构建当前平台的二进制文件"
	@echo "  build-all  - 构建所有支持平台的二进制文件"
	@echo "  test       - 运行测试"
	@echo "  deps       - 下载和整理依赖"
	@echo "  clean      - 清理构建文件"
	@echo "  install    - 安装到本地 (/usr/local/bin)"
	@echo "  uninstall  - 从本地卸载"
	@echo "  version    - 显示版本信息"
	@echo "  run        - 构建并运行（显示帮助）"
	@echo "  fmt        - 格式化代码"
	@echo "  lint       - 运行代码检查"
	@echo "  help       - 显示此帮助信息"
	@echo ""
	@echo "支持的平台:"
	@$(foreach platform,$(PLATFORMS),echo "  $(platform)";)
