package generator

// GetPythonAPITemplate 返回 Python API 服务模板
func GetPythonAPITemplate() *Template {
	return &Template{
		Name:        "python-api",
		Description: "Python API 服务模板 - FastAPI 生产级应用",
		Features:    "FastAPI, pytest, pre-commit, mypy, black, ruff, Docker, CI/CD",
		Files: []FileTemplate{
			// pyproject.toml
			{
				Path: "pyproject.toml",
				Content: `[tool.poetry]
name = "{{.ProjectName}}"
version = "1.0.0"
description = "FastAPI 服务"
authors = ["Your Name <you@example.com>"]
readme = "README.md"
packages = [{include = "app"}]

[tool.poetry.dependencies]
python = "^3.11"
fastapi = "^0.109.0"
uvicorn = {extras = ["standard"], version = "^0.27.0"}
pydantic = "^2.5.3"
pydantic-settings = "^2.1.0"
python-json-logger = "^2.0.7"
prometheus-client = "^0.19.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.4.4"
pytest-cov = "^4.1.0"
pytest-asyncio = "^0.23.3"
httpx = "^0.26.0"
black = "^23.12.1"
ruff = "^0.1.11"
mypy = "^1.8.0"
pre-commit = "^3.6.0"

[tool.black]
line-length = 100
target-version = ['py311']

[tool.ruff]
line-length = 100
target-version = "py311"

[tool.mypy]
python_version = "3.11"
strict = true
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true

[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
python_classes = ["Test*"]
python_functions = ["test_*"]
asyncio_mode = "auto"

[tool.coverage.run]
source = ["app"]
omit = ["*/tests/*", "*/test_*.py"]

[tool.coverage.report]
exclude_lines = [
    "pragma: no cover",
    "def __repr__",
    "raise AssertionError",
    "raise NotImplementedError",
    "if __name__ == .__main__.:",
]

[build-system]
requires = ["poetry-core"]
build-backend = "poetry.core.masonry.api"
`,
			},

			// main.py
			{
				Path: "app/main.py",
				Content: `"""
{{.ProjectName}} FastAPI 应用主入口
"""
import logging
from contextlib import asynccontextmanager
from typing import AsyncGenerator

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from prometheus_client import make_asgi_app

from app.api.router import api_router
from app.core.config import settings
from app.core.logging import setup_logging


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncGenerator[None, None]:
    """应用生命周期管理"""
    # 启动
    logging.info("🚀 应用启动中...")
    logging.info(f"环境: {settings.ENVIRONMENT}")
    logging.info(f"服务名: {settings.SERVICE_NAME}")
    logging.info(f"版本: {settings.VERSION}")

    yield

    # 关闭
    logging.info("👋 应用关闭中...")


def create_app() -> FastAPI:
    """创建 FastAPI 应用"""
    # 设置日志
    setup_logging(settings.LOG_LEVEL)

    app = FastAPI(
        title=settings.SERVICE_NAME,
        description="FastAPI 服务",
        version=settings.VERSION,
        lifespan=lifespan,
    )

    # CORS 中间件
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.ALLOWED_ORIGINS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # 路由
    app.include_router(api_router, prefix=settings.API_V1_PREFIX)

    # Prometheus metrics
    metrics_app = make_asgi_app()
    app.mount("/metrics", metrics_app)

    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.PORT,
        reload=settings.ENVIRONMENT == "development",
        log_level=settings.LOG_LEVEL.lower(),
    )
`,
			},

			// 配置
			{
				Path: "app/core/config.py",
				Content: `"""
应用配置
"""
from typing import List

from pydantic import Field
from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """应用设置"""

    # 基础配置
    SERVICE_NAME: str = Field(default="{{.ProjectName}}", description="服务名称")
    VERSION: str = Field(default="1.0.0", description="版本")
    ENVIRONMENT: str = Field(default="development", description="环境")

    # 服务器配置
    HOST: str = Field(default="0.0.0.0", description="主机地址")
    PORT: int = Field(default=8000, description="端口")

    # API 配置
    API_V1_PREFIX: str = Field(default="/api/v1", description="API v1 前缀")

    # 日志配置
    LOG_LEVEL: str = Field(default="INFO", description="日志级别")

    # CORS 配置
    ALLOWED_ORIGINS: List[str] = Field(
        default=["http://localhost:3000"],
        description="允许的跨域源"
    )

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=True,
    )


settings = Settings()
`,
			},

			// 日志配置
			{
				Path: "app/core/logging.py",
				Content: `"""
日志配置
"""
import logging
import sys
from typing import Any

from pythonjsonlogger import jsonlogger


class CustomJsonFormatter(jsonlogger.JsonFormatter):
    """自定义 JSON 日志格式"""

    def add_fields(
        self,
        log_record: dict[str, Any],
        record: logging.LogRecord,
        message_dict: dict[str, Any],
    ) -> None:
        super().add_fields(log_record, record, message_dict)
        log_record["level"] = record.levelname
        log_record["logger"] = record.name


def setup_logging(level: str = "INFO") -> None:
    """设置日志"""
    # 根日志器
    root_logger = logging.getLogger()
    root_logger.setLevel(level)

    # 清除现有处理器
    root_logger.handlers.clear()

    # 控制台处理器
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(level)

    # JSON 格式
    formatter = CustomJsonFormatter(
        "%(timestamp)s %(level)s %(name)s %(message)s",
        rename_fields={"asctime": "timestamp"},
    )
    console_handler.setFormatter(formatter)

    root_logger.addHandler(console_handler)
`,
			},

			// API 路由
			{
				Path: "app/api/router.py",
				Content: `"""
API 路由
"""
from fastapi import APIRouter

from app.api.endpoints import health, hello

api_router = APIRouter()

# 包含各个端点路由
api_router.include_router(health.router, tags=["health"])
api_router.include_router(hello.router, tags=["hello"])
`,
			},

			// 健康检查端点
			{
				Path: "app/api/endpoints/health.py",
				Content: `"""
健康检查端点
"""
import time
from typing import Dict, Any

from fastapi import APIRouter

router = APIRouter()

# 记录启动时间
START_TIME = time.time()


@router.get("/health", response_model=Dict[str, Any])
async def health_check() -> Dict[str, Any]:
    """健康检查"""
    return {
        "status": "ok",
        "uptime": time.time() - START_TIME,
    }


@router.get("/ready", response_model=Dict[str, Any])
async def readiness_check() -> Dict[str, Any]:
    """就绪检查"""
    # 这里可以检查依赖服务（数据库、缓存等）
    return {
        "status": "ready",
    }
`,
			},

			// Hello 端点
			{
				Path: "app/api/endpoints/hello.py",
				Content: `"""
Hello 示例端点
"""
import logging
from datetime import datetime
from typing import Optional

from fastapi import APIRouter, Query
from pydantic import BaseModel, Field
from prometheus_client import Counter, Histogram

router = APIRouter()
logger = logging.getLogger(__name__)

# Prometheus metrics
hello_requests_total = Counter(
    "hello_requests_total",
    "Hello 请求总数",
)

hello_request_duration = Histogram(
    "hello_request_duration_seconds",
    "Hello 请求延迟",
)


class HelloResponse(BaseModel):
    """Hello 响应"""
    message: str = Field(..., description="问候消息")
    timestamp: str = Field(..., description="时间戳")


class EchoRequest(BaseModel):
    """Echo 请求"""
    message: str = Field(..., description="要回显的消息", min_length=1, max_length=1000)


class EchoResponse(BaseModel):
    """Echo 响应"""
    echo: str = Field(..., description="回显的消息")
    timestamp: str = Field(..., description="时间戳")


@router.get("/hello", response_model=HelloResponse)
@hello_request_duration.time()
async def hello(name: Optional[str] = Query(default="世界", description="名称")) -> HelloResponse:
    """Hello 端点"""
    hello_requests_total.inc()
    logger.info(f"处理 Hello 请求: name={name}")

    return HelloResponse(
        message=f"你好, {name}!",
        timestamp=datetime.utcnow().isoformat(),
    )


@router.post("/echo", response_model=EchoResponse)
async def echo(request: EchoRequest) -> EchoResponse:
    """Echo 端点"""
    logger.info(f"处理 Echo 请求: message={request.message}")

    return EchoResponse(
        echo=request.message,
        timestamp=datetime.utcnow().isoformat(),
    )
`,
			},

			// __init__.py 文件
			{Path: "app/__init__.py", Content: `"""{{.ProjectName}} 应用包"""`},
			{Path: "app/core/__init__.py", Content: `"""核心模块"""`},
			{Path: "app/api/__init__.py", Content: `"""API 模块"""`},
			{Path: "app/api/endpoints/__init__.py", Content: `"""API 端点"""`},

			// 测试
			{
				Path: "tests/conftest.py",
				Content: `"""
pytest 配置
"""
import pytest
from fastapi.testclient import TestClient

from app.main import create_app


@pytest.fixture
def client() -> TestClient:
    """测试客户端"""
    app = create_app()
    return TestClient(app)
`,
			},

			{
				Path: "tests/test_health.py",
				Content: `"""
健康检查测试
"""
from fastapi.testclient import TestClient


def test_health_check(client: TestClient) -> None:
    """测试健康检查"""
    response = client.get("/api/v1/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ok"
    assert "uptime" in data


def test_readiness_check(client: TestClient) -> None:
    """测试就绪检查"""
    response = client.get("/api/v1/ready")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "ready"
`,
			},

			{
				Path: "tests/test_hello.py",
				Content: `"""
Hello 端点测试
"""
from fastapi.testclient import TestClient


def test_hello_default(client: TestClient) -> None:
    """测试默认 Hello"""
    response = client.get("/api/v1/hello")
    assert response.status_code == 200
    data = response.json()
    assert data["message"] == "你好, 世界!"
    assert "timestamp" in data


def test_hello_with_name(client: TestClient) -> None:
    """测试带名称的 Hello"""
    response = client.get("/api/v1/hello?name=Python")
    assert response.status_code == 200
    data = response.json()
    assert data["message"] == "你好, Python!"


def test_echo(client: TestClient) -> None:
    """测试 Echo"""
    payload = {"message": "测试消息"}
    response = client.post("/api/v1/echo", json=payload)
    assert response.status_code == 200
    data = response.json()
    assert data["echo"] == payload["message"]
    assert "timestamp" in data


def test_echo_validation(client: TestClient) -> None:
    """测试 Echo 验证"""
    # 空消息
    response = client.post("/api/v1/echo", json={"message": ""})
    assert response.status_code == 422
`,
			},

			{Path: "tests/__init__.py", Content: ``},

			// .env.example
			{
				Path: ".env.example",
				Content: `# 服务配置
SERVICE_NAME={{.ProjectName}}
VERSION=1.0.0
ENVIRONMENT=development

# 服务器配置
HOST=0.0.0.0
PORT=8000

# 日志配置
LOG_LEVEL=INFO

# CORS 配置
ALLOWED_ORIGINS=["http://localhost:3000"]
`,
			},

			// requirements.txt (用于不使用 poetry 的情况)
			{
				Path: "requirements.txt",
				Content: `fastapi>=0.109.0
uvicorn[standard]>=0.27.0
pydantic>=2.5.3
pydantic-settings>=2.1.0
python-json-logger>=2.0.7
prometheus-client>=0.19.0
`,
			},

			{
				Path: "requirements-dev.txt",
				Content: `pytest>=7.4.4
pytest-cov>=4.1.0
pytest-asyncio>=0.23.3
httpx>=0.26.0
black>=23.12.1
ruff>=0.1.11
mypy>=1.8.0
pre-commit>=3.6.0
`,
			},

			// Makefile
			{
				Path: "Makefile",
				Content: `.PHONY: help install install-dev test lint format run docker clean

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## 安装生产依赖
	pip install -r requirements.txt

install-dev: ## 安装开发依赖
	pip install -r requirements.txt -r requirements-dev.txt
	pre-commit install

install-poetry: ## 使用 Poetry 安装依赖
	poetry install

test: ## 运行测试
	pytest -v --cov=app --cov-report=html --cov-report=term

test-watch: ## 监视模式运行测试
	pytest-watch

lint: ## 代码检查
	ruff check app tests
	mypy app

format: ## 格式化代码
	black app tests
	ruff check --fix app tests

run: ## 运行服务
	python -m app.main

run-dev: ## 开发模式运行（热重载）
	uvicorn app.main:app --reload --host 0.0.0.0 --port 8000

docker-build: ## 构建 Docker 镜像
	docker build -t {{.ProjectName}}:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8000:8000 --env-file .env {{.ProjectName}}:latest

clean: ## 清理缓存文件
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	rm -rf .pytest_cache .coverage htmlcov .mypy_cache .ruff_cache

.DEFAULT_GOAL := help
`,
			},

			// Dockerfile
			{
				Path: "Dockerfile",
				Content: `# Python 基础镜像
FROM python:3.11-slim as base

# 设置工作目录
WORKDIR /app

# 安装系统依赖
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# 复制依赖文件
COPY requirements.txt .

# 安装 Python 依赖
RUN pip install --no-cache-dir -r requirements.txt

# 复制应用代码
COPY app ./app

# 创建非 root 用户
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

# 暴露端口
EXPOSE 8000

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD python -c "import requests; requests.get('http://localhost:8000/api/v1/health')"

# 运行应用
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
`,
			},

			// .pre-commit-config.yaml
			{
				Path: ".pre-commit-config.yaml",
				Content: `repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
      - id: check-json
      - id: check-toml
      - id: check-merge-conflict
      - id: detect-private-key

  - repo: https://github.com/psf/black
    rev: 23.12.1
    hooks:
      - id: black
        language_version: python3.11

  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.1.11
    hooks:
      - id: ruff
        args: [--fix, --exit-non-zero-on-fix]

  - repo: https://github.com/pre-commit/mirrors-mypy
    rev: v1.8.0
    hooks:
      - id: mypy
        additional_dependencies: [types-all]
        args: [--strict]
`,
			},

			// GitHub Actions
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
    strategy:
      matrix:
        python-version: ["3.11", "3.12"]

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 Python
        uses: actions/setup-python@v4
        with:
          python-version: ${{ matrix.python-version }}

      - name: 安装依赖
        run: |
          pip install -r requirements.txt -r requirements-dev.txt

      - name: 运行测试
        run: |
          pytest -v --cov=app --cov-report=xml

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.xml

  lint:
    name: 代码检查
    runs-on: ubuntu-latest

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 Python
        uses: actions/setup-python@v4
        with:
          python-version: "3.11"

      - name: 安装依赖
        run: |
          pip install black ruff mypy

      - name: Black 格式检查
        run: black --check app tests

      - name: Ruff 检查
        run: ruff check app tests

      - name: Mypy 类型检查
        run: mypy app

  build:
    name: 构建 Docker 镜像
    runs-on: ubuntu-latest
    needs: [test, lint]

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 构建镜像
        run: docker build -t {{.ProjectName}}:latest .
`,
			},

			// .gitignore
			{
				Path: ".gitignore",
				Content: `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
build/
develop-eggs/
dist/
downloads/
eggs/
.eggs/
lib/
lib64/
parts/
sdist/
var/
wheels/
*.egg-info/
.installed.cfg
*.egg

# 虚拟环境
venv/
env/
ENV/
.venv

# 测试
.pytest_cache/
.coverage
htmlcov/
.tox/
.mypy_cache/
.ruff_cache/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# 环境变量
.env
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

Python FastAPI 服务 - 使用 dev-template 生成

## 特性

- ✅ FastAPI 异步框架
- ✅ Pydantic v2 数据验证
- ✅ 结构化 JSON 日志
- ✅ Prometheus metrics
- ✅ 完整的类型注解 (mypy strict)
- ✅ pytest 单元测试
- ✅ pre-commit hooks
- ✅ Black + Ruff 代码格式化
- ✅ Docker 支持
- ✅ GitHub Actions CI/CD

## 快速开始

### 创建虚拟环境

` + "```bash" + `
python -m venv venv
source venv/bin/activate  # Linux/Mac
# 或
.\venv\Scripts\activate  # Windows
` + "```" + `

### 安装依赖

` + "```bash" + `
make install-dev
` + "```" + `

或使用 Poetry:

` + "```bash" + `
poetry install
` + "```" + `

### 运行服务

` + "```bash" + `
make run-dev
` + "```" + `

服务将在 http://localhost:8000 启动

API 文档: http://localhost:8000/docs

### 运行测试

` + "```bash" + `
make test
` + "```" + `

### 代码检查和格式化

` + "```bash" + `
make lint    # 代码检查
make format  # 格式化代码
` + "```" + `

## API 端点

- ` + "`GET /api/v1/health`" + ` - 健康检查
- ` + "`GET /api/v1/ready`" + ` - 就绪检查
- ` + "`GET /metrics`" + ` - Prometheus metrics
- ` + "`GET /api/v1/hello?name=xxx`" + ` - 示例 API
- ` + "`POST /api/v1/echo`" + ` - Echo API
- ` + "`GET /docs`" + ` - Swagger UI 文档
- ` + "`GET /redoc`" + ` - ReDoc 文档

## 环境变量

复制 ` + "`.env.example`" + ` 到 ` + "`.env`" + ` 并修改配置：

` + "```bash" + `
cp .env.example .env
` + "```" + `

## Docker

### 构建镜像

` + "```bash" + `
make docker-build
` + "```" + `

### 运行容器

` + "```bash" + `
make docker-run
` + "```" + `

## 开发工具

### Pre-commit

安装 pre-commit hooks:

` + "```bash" + `
pre-commit install
` + "```" + `

手动运行所有 hooks:

` + "```bash" + `
pre-commit run --all-files
` + "```" + `

## 项目结构

` + "```" + `
.
├── app/
│   ├── main.py              # 应用入口
│   ├── core/
│   │   ├── config.py        # 配置
│   │   └── logging.py       # 日志
│   └── api/
│       ├── router.py        # 路由
│       └── endpoints/       # 端点
├── tests/                   # 测试
├── Dockerfile              # Docker 配置
├── pyproject.toml          # Poetry 配置
├── Makefile               # 构建脚本
└── .pre-commit-config.yaml # Pre-commit 配置
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
