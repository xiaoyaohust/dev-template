.PHONY: help build install test clean run demo-go demo-python demo-java demo-kotlin demo-scala list

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## 安装依赖
	go mod download
	go mod tidy

build: ## 构建二进制文件
	go build -o bin/dev-template main.go

install: ## 安装到系统
	go install

test: ## 运行测试
	go test -v -race -cover ./...

clean: ## 清理构建文件
	rm -rf bin/
	go clean

run: ## 运行 CLI (示例)
	go run main.go

# 示例命令
demo-go: build ## 演示创建 Go 项目
	./bin/dev-template new go-service demo-go-service

demo-python: build ## 演示创建 Python 项目
	./bin/dev-template new python-api demo-python-api

demo-java: build ## 演示创建 Java 项目
	./bin/dev-template new java-service demo-java-service

demo-kotlin: build ## 演示创建 Kotlin 项目
	./bin/dev-template new kotlin-service demo-kotlin-service

demo-scala: build ## 演示创建 Scala 项目
	./bin/dev-template new scala-service demo-scala-service

list: build ## 列出所有可用模板
	./bin/dev-template list

.DEFAULT_GOAL := help
