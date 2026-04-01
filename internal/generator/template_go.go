package generator

// GetGoServiceTemplate 返回 Go 微服务模板
func GetGoServiceTemplate() *Template {
	return &Template{
		Name:        "go-service",
		Description: "Go 微服务模板 - 生产级 REST API 服务",
		Features:    "Gin/Fiber, 配置管理, 结构化日志, 健康检查, Prometheus metrics, Docker, CI/CD",
		Files: []FileTemplate{
			// go.mod
			{
				Path: "go.mod",
				Content: `module {{.ModulePath}}

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/prometheus/client_golang v1.18.0
	github.com/spf13/viper v1.18.2
	go.uber.org/zap v1.26.0
)
`,
			},

			// main.go
			{
				Path: "main.go",
				Content: `package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/handler"
	"{{.ModulePath}}/internal/logger"
	"{{.ModulePath}}/internal/server"
)

func main() {
	// 初始化配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log := logger.New(cfg.Log.Level)
	defer log.Sync()

	log.Info("启动服务",
		"service", cfg.Service.Name,
		"version", cfg.Service.Version,
		"port", cfg.Server.Port,
	)

	// 创建 HTTP 服务器
	h := handler.New(log)
	srv := server.New(cfg, h, log)

	// 启动服务器
	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal("服务器启动失败", "error", err)
		}
	}()

	log.Info("服务器已启动", "address", fmt.Sprintf(":%d", cfg.Server.Port))

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器强制关闭", "error", err)
	}

	log.Info("服务器已关闭")
}
`,
			},

			// 配置管理
			{
				Path: "internal/config/config.go",
				Content: `package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Service ServiceConfig ` + "`mapstructure:\"service\"`" + `
	Server  ServerConfig  ` + "`mapstructure:\"server\"`" + `
	Log     LogConfig     ` + "`mapstructure:\"log\"`" + `
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Name    string ` + "`mapstructure:\"name\"`" + `
	Version string ` + "`mapstructure:\"version\"`" + `
	Env     string ` + "`mapstructure:\"env\"`" + `
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int  ` + "`mapstructure:\"port\"`" + `
	ReadTimeout  int  ` + "`mapstructure:\"read_timeout\"`" + `
	WriteTimeout int  ` + "`mapstructure:\"write_timeout\"`" + `
	EnablePprof  bool ` + "`mapstructure:\"enable_pprof\"`" + `
}

// LogConfig 日志配置
type LogConfig struct {
	Level string ` + "`mapstructure:\"level\"`" + `
}

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 从环境变量读取
	v.AutomaticEnv()

	// 从配置文件读取
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	// 配置文件可选
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("service.name", "{{.ProjectName}}")
	v.SetDefault("service.version", "1.0.0")
	v.SetDefault("service.env", "development")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 10)
	v.SetDefault("server.write_timeout", 10)
	v.SetDefault("server.enable_pprof", false)

	v.SetDefault("log.level", "info")
}
`,
			},

			// 日志
			{
				Path: "internal/logger/logger.go",
				Content: `package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New 创建新的日志实例
func New(level string) *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 设置日志级别
	zapLevel := zapcore.InfoLevel
	_ = zapLevel.UnmarshalText([]byte(level))
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, _ := config.Build()
	return logger.Sugar()
}
`,
			},

			// HTTP 服务器
			{
				Path: "internal/server/server.go",
				Content: `package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"{{.ModulePath}}/internal/config"
	"{{.ModulePath}}/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Server HTTP 服务器
type Server struct {
	config  *config.Config
	handler *handler.Handler
	logger  *zap.SugaredLogger
	engine  *gin.Engine
	server  *http.Server
}

// New 创建新的服务器实例
func New(cfg *config.Config, h *handler.Handler, log *zap.SugaredLogger) *Server {
	// 设置 Gin 模式
	if cfg.Service.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	s := &Server{
		config:  cfg,
		handler: h,
		logger:  log,
		engine:  engine,
	}

	s.setupRoutes()

	return s
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 健康检查
	s.engine.GET("/health", s.handler.Health)
	s.engine.GET("/ready", s.handler.Ready)

	// Metrics
	s.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API 路由
	api := s.engine.Group("/api/v1")
	{
		api.GET("/hello", s.handler.Hello)
		api.POST("/echo", s.handler.Echo)
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.Port),
		Handler:      s.engine,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}
`,
			},

			// HTTP 处理器
			{
				Path: "internal/handler/handler.go",
				Content: `package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// HTTP 请求计数器
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "HTTP 请求总数",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTP 请求延迟
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP 请求延迟",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

// Handler HTTP 处理器
type Handler struct {
	logger    *zap.SugaredLogger
	startTime time.Time
}

// New 创建新的处理器
func New(log *zap.SugaredLogger) *Handler {
	return &Handler{
		logger:    log,
		startTime: time.Now(),
	}
}

// Health 健康检查
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"uptime": time.Since(h.startTime).String(),
	})
}

// Ready 就绪检查
func (h *Handler) Ready(c *gin.Context) {
	// 这里可以检查依赖服务（数据库、缓存等）是否就绪
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}

// Hello 示例处理器
func (h *Handler) Hello(c *gin.Context) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), "200").Inc()
	}()

	name := c.DefaultQuery("name", "世界")
	h.logger.Infow("处理 Hello 请求", "name", name)

	c.JSON(http.StatusOK, gin.H{
		"message": "你好, " + name + "!",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// EchoRequest Echo 请求体
type EchoRequest struct {
	Message string ` + "`json:\"message\" binding:\"required\"`" + `
}

// Echo 回显处理器
func (h *Handler) Echo(c *gin.Context) {
	var req EchoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw("无效的请求", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求: " + err.Error(),
		})
		return
	}

	h.logger.Infow("处理 Echo 请求", "message", req.Message)

	c.JSON(http.StatusOK, gin.H{
		"echo": req.Message,
		"time": time.Now().Format(time.RFC3339),
	})
}
`,
			},

			// 测试文件
			{
				Path: "internal/handler/handler_test.go",
				Content: `package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTestHandler() *Handler {
	logger := zap.NewNop().Sugar()
	return New(logger)
}

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/health", nil)

	handler.Health(c)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if status, ok := resp["status"].(string); !ok || status != "ok" {
		t.Errorf("期望 status=ok, 得到 %v", resp["status"])
	}
}

func TestHello(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler()

	tests := []struct {
		name     string
		query    string
		expected string
	}{
		{"默认名称", "", "你好, 世界!"},
		{"自定义名称", "?name=Go", "你好, Go!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/v1/hello"+tt.query, nil)

			handler.Hello(c)

			if w.Code != http.StatusOK {
				t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
			}

			var resp map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("解析响应失败: %v", err)
			}

			if msg, ok := resp["message"].(string); !ok || msg != tt.expected {
				t.Errorf("期望 message=%s, 得到 %v", tt.expected, resp["message"])
			}
		})
	}
}

func TestEcho(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := setupTestHandler()

	reqBody := EchoRequest{Message: "测试消息"}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/api/v1/echo", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Echo(c)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 得到 %d", http.StatusOK, w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if echo, ok := resp["echo"].(string); !ok || echo != reqBody.Message {
		t.Errorf("期望 echo=%s, 得到 %v", reqBody.Message, resp["echo"])
	}
}
`,
			},

			// 配置文件
			{
				Path: "configs/config.yaml",
				Content: `service:
  name: {{.ProjectName}}
  version: "1.0.0"
  env: development

server:
  port: 8080
  read_timeout: 10
  write_timeout: 10
  enable_pprof: false

log:
  level: info
`,
			},

			// Dockerfile
			{
				Path: "Dockerfile",
				Content: `# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /build

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建二进制文件
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /build/app .
COPY --from=builder /build/configs ./configs

# 暴露端口
EXPOSE 8080

# 运行
CMD ["./app"]
`,
			},

			// Makefile
			{
				Path: "Makefile",
				Content: `.PHONY: help deps test build run docker clean

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go mod tidy

test: ## 运行测试
	go test -v -race -cover ./...

test-coverage: ## 生成测试覆盖率报告
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build: ## 构建二进制文件
	go build -o bin/{{.ProjectName}} .

run: ## 运行服务
	go run main.go

docker-build: ## 构建 Docker 镜像
	docker build -t {{.ProjectName}}:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 --name {{.ProjectName}} {{.ProjectName}}:latest

lint: ## 代码检查
	golangci-lint run

fmt: ## 格式化代码
	go fmt ./...
	goimports -w .

clean: ## 清理构建文件
	rm -rf bin/
	rm -f coverage.out coverage.html

.DEFAULT_GOAL := help
`,
			},

			// GitHub Actions CI
			{
				Path: ".github/workflows/ci.yml",
				Content: `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    name: 测试
    runs-on: ubuntu-latest
    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: 安装依赖
        run: go mod download

      - name: 运行测试
        run: go test -v -race -coverprofile=coverage.out ./...

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out

  lint:
    name: 代码检查
    runs-on: ubuntu-latest
    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest

  build:
    name: 构建
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: 构建
        run: go build -v .
`,
			},

			// .gitignore
			{
				Path: ".gitignore",
				Content: `# 二进制文件
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# 测试
*.test
*.out
coverage.out
coverage.html

# Go
vendor/

# IDE
.idea/
.vscode/
*.swp
*.swo
*~

# 配置
*.env
.env.local

# 日志
*.log

# OS
.DS_Store
Thumbs.db
`,
			},

			// README
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

Go 微服务项目 - 使用 dev-template 生成

## 特性

- ✅ Gin Web 框架
- ✅ 结构化日志 (zap)
- ✅ 配置管理 (viper)
- ✅ 健康检查端点
- ✅ Prometheus metrics
- ✅ 优雅关闭
- ✅ Docker 支持
- ✅ GitHub Actions CI/CD
- ✅ 单元测试

## 快速开始

### 安装依赖

` + "```bash" + `
make deps
` + "```" + `

### 运行服务

` + "```bash" + `
make run
` + "```" + `

服务将在 http://localhost:8080 启动

### 运行测试

` + "```bash" + `
make test
` + "```" + `

### 构建

` + "```bash" + `
make build
` + "```" + `

## API 端点

- ` + "`GET /health`" + ` - 健康检查
- ` + "`GET /ready`" + ` - 就绪检查
- ` + "`GET /metrics`" + ` - Prometheus metrics
- ` + "`GET /api/v1/hello?name=xxx`" + ` - 示例 API
- ` + "`POST /api/v1/echo`" + ` - Echo API

## Docker

### 构建镜像

` + "```bash" + `
make docker-build
` + "```" + `

### 运行容器

` + "```bash" + `
make docker-run
` + "```" + `

## 配置

配置文件位于 ` + "`configs/config.yaml`" + `，也可以通过环境变量覆盖：

- ` + "`SERVICE_NAME`" + `
- ` + "`SERVER_PORT`" + `
- ` + "`LOG_LEVEL`" + `

## 项目结构

` + "```" + `
.
├── main.go              # 入口文件
├── internal/
│   ├── config/         # 配置管理
│   ├── handler/        # HTTP 处理器
│   ├── logger/         # 日志
│   └── server/         # HTTP 服务器
├── configs/            # 配置文件
├── Dockerfile          # Docker 配置
├── Makefile           # 构建脚本
└── .github/workflows/ # CI/CD 配置
` + "```" + `

## License

MIT
`,
			},
		},
	}
}
`,
			},