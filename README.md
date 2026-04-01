# dev-template

> 🚀 高质量项目模板生成器 - 为有经验的工程师打造

一个**非脚手架式**的项目模板生成工具，专为有经验的工程师设计。快速创建生产级 Go、Python、Java 服务模板，内置最佳实践：配置管理、日志、测试、CI/CD、Docker、Kubernetes 等。

非常适合：
- 🎯 想从零搭建干净项目的开发者
- 🏢 小团队快速启动新服务
- 💡 Side Project 快速开发
- 📚 学习最佳实践的工程师

## ✨ 特性

### 三套完整模板

#### 1️⃣ **Go 微服务模板** (`go-service`)
- ✅ Gin Web 框架
- ✅ 结构化日志 (zap)
- ✅ 配置管理 (viper)
- ✅ 健康检查端点
- ✅ Prometheus metrics
- ✅ 优雅关闭
- ✅ 完整的单元测试
- ✅ Docker + CI/CD

#### 2️⃣ **Python API 服务模板** (`python-api`)
- ✅ FastAPI 异步框架
- ✅ Pydantic v2 数据验证
- ✅ 完整类型注解 (mypy strict)
- ✅ pytest + 测试覆盖率
- ✅ pre-commit hooks
- ✅ Black + Ruff 格式化
- ✅ Prometheus metrics
- ✅ Docker + CI/CD

#### 3️⃣ **Java Spring Boot 服务模板** (`java-service`)
- ✅ Spring Boot 3.2
- ✅ Spring Actuator
- ✅ JUnit 5 测试
- ✅ Checkstyle 代码检查
- ✅ JaCoCo 覆盖率
- ✅ Logback 日志
- ✅ Prometheus metrics
- ✅ Kubernetes 配置
- ✅ Docker + CI/CD

## 🚀 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/ariesxiao/dev-template.git
cd dev-template

# 构建
make build

# 或直接安装到系统
make install
```

### 使用

```bash
# 创建 Go 微服务
dev-template new go-service my-orders-service

# 创建 Python API 服务
dev-template new python-api user-api

# 创建 Java Spring Boot 服务
dev-template new java-service payment-service

# 列出所有可用模板
dev-template list
```

## 📖 使用示例

### 创建 Go 服务

```bash
# 生成项目
dev-template new go-service my-service

# 进入项目
cd my-service

# 安装依赖
make deps

# 运行测试
make test

# 启动服务
make run
```

服务启动后访问：
- http://localhost:8080/health - 健康检查
- http://localhost:8080/metrics - Prometheus metrics
- http://localhost:8080/api/v1/hello - 示例 API

### 创建 Python 服务

```bash
# 生成项目
dev-template new python-api my-api

# 进入项目
cd my-api

# 创建虚拟环境
python -m venv venv
source venv/bin/activate

# 安装依赖
make install-dev

# 运行测试
make test

# 启动服务
make run-dev
```

服务启动后访问：
- http://localhost:8000/docs - Swagger UI 文档
- http://localhost:8000/api/v1/health - 健康检查
- http://localhost:8000/metrics - Prometheus metrics

### 创建 Java 服务

```bash
# 生成项目
dev-template new java-service my-service

# 进入项目
cd my-service

# 构建项目
./mvnw clean install

# 运行测试
./mvnw test

# 启动服务
./mvnw spring-boot:run
```

服务启动后访问：
- http://localhost:8080/actuator/health - 健康检查
- http://localhost:8080/actuator/prometheus - Prometheus metrics
- http://localhost:8080/api/v1/hello - 示例 API

## 🎯 模板详情

### Go 微服务模板

生成的项目结构：
```
my-service/
├── main.go              # 应用入口
├── internal/
│   ├── config/         # 配置管理 (viper)
│   ├── handler/        # HTTP 处理器
│   ├── logger/         # 日志 (zap)
│   └── server/         # HTTP 服务器 (Gin)
├── configs/            # 配置文件
├── Dockerfile          # Docker 配置
├── Makefile           # 构建脚本
└── .github/workflows/ # CI/CD 配置
```

**特色功能：**
- 🔥 热重载开发
- 📊 内置 Prometheus metrics
- 🏥 完整的健康检查
- ⚡ 优雅关闭
- 🧪 单元测试示例

### Python API 服务模板

生成的项目结构：
```
my-api/
├── app/
│   ├── main.py              # 应用入口
│   ├── core/
│   │   ├── config.py        # 配置 (pydantic-settings)
│   │   └── logging.py       # 结构化日志
│   └── api/
│       ├── router.py        # 路由
│       └── endpoints/       # API 端点
├── tests/                   # pytest 测试
├── pyproject.toml          # Poetry 配置
├── .pre-commit-config.yaml # Pre-commit hooks
├── Dockerfile              # Docker 配置
└── Makefile               # 构建脚本
```

**特色功能：**
- 🐍 完整类型注解 (mypy strict mode)
- ✨ pre-commit hooks (black, ruff, mypy)
- 📝 自动生成 OpenAPI 文档
- 🧪 pytest + 覆盖率报告
- ⚡ 异步 FastAPI

### Java Spring Boot 服务模板

生成的项目结构：
```
my-service/
├── src/
│   ├── main/
│   │   ├── java/com/example/
│   │   │   ├── Application.java  # 主类
│   │   │   ├── config/          # 配置
│   │   │   ├── controller/      # REST 控制器
│   │   │   └── dto/             # 数据传输对象
│   │   └── resources/
│   │       ├── application.yml  # 配置文件
│   │       └── logback-spring.xml # 日志配置
│   └── test/                    # JUnit 测试
├── k8s/                         # Kubernetes 配置
├── Dockerfile                   # Docker 配置
├── pom.xml                      # Maven 配置
├── checkstyle.xml              # Checkstyle 配置
└── Makefile                    # 构建脚本
```

**特色功能：**
- ☕ Spring Boot 3.2 + Java 17
- 🏥 Spring Actuator 健康检查
- 📊 Micrometer + Prometheus
- 🧪 JUnit 5 + JaCoCo 覆盖率
- ☸️ 生产级 Kubernetes 配置
- 🔍 Checkstyle 代码规范

## 🛠️ 开发

### 构建项目

```bash
make build
```

### 运行测试

```bash
make test
```

### 演示命令

```bash
# 演示生成 Go 项目
make demo-go

# 演示生成 Python 项目
make demo-python

# 演示生成 Java 项目
make demo-java

# 列出所有模板
make list
```

## 📦 项目结构

```
dev-template/
├── main.go                      # CLI 入口
├── cmd/dev-template/commands/   # CLI 命令
│   ├── root.go                 # 根命令
│   ├── new.go                  # new 命令
│   └── list.go                 # list 命令
├── internal/generator/          # 模板生成器
│   ├── generator.go            # 核心生成逻辑
│   ├── template_go.go          # Go 模板
│   ├── template_python.go      # Python 模板
│   └── template_java.go        # Java 模板
├── Makefile                    # 构建脚本
└── README.md                   # 文档
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

### 添加新模板

1. 在 `internal/generator/` 中创建新的模板文件
2. 实现 `GetXXXTemplate()` 函数
3. 在 `generator.go` 的 `initTemplates()` 中注册
4. 更新文档

### 改进现有模板

1. 编辑对应的模板文件（`template_*.go`）
2. 测试生成的项目
3. 提交 PR

## 📝 许可证

MIT License

## 🙏 致谢

感谢所有开源项目：
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [FastAPI](https://github.com/tiangolo/fastapi) - Python Web 框架
- [Spring Boot](https://spring.io/projects/spring-boot) - Java 框架

---

**享受编码！** 🎉
