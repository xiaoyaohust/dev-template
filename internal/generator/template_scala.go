package generator

// GetScalaServiceTemplate 返回 Scala Play Framework 服务模板
func GetScalaServiceTemplate() *Template {
	return &Template{
		Name:        "scala-service",
		Description: "Scala Play Framework 服务模板 - 函数式 JVM 微服务",
		Features:    "Scala 3, Play Framework 3, sbt, ScalaTest, Prometheus, Logback, Docker, CI/CD",
		Files: []FileTemplate{
			// build.sbt
			{
				Path: "build.sbt",
				Content: `name := "{{.ProjectName}}"
version := "1.0.0"
scalaVersion := "3.3.1"

lazy val root = (project in file("."))
  .enablePlugins(PlayScala)
  .settings(
    libraryDependencies ++= Seq(
      guice,
      "com.typesafe.play" %% "play-json" % "3.0.2",
      // Prometheus metrics
      "io.prometheus" % "simpleclient"         % "0.16.0",
      "io.prometheus" % "simpleclient_hotspot" % "0.16.0",
      "io.prometheus" % "simpleclient_common"  % "0.16.0",
      // Test
      "org.scalatestplus.play" %% "scalatestplus-play" % "7.0.1" % Test,
    ),
    // Disable unused warning for generated routes
    scalacOptions ++= Seq(
      "-deprecation",
      "-feature",
      "-unchecked",
      "-Xfatal-warnings",
    ),
  )

// Assembly settings for fat JAR
assembly / assemblyMergeStrategy := {
  case PathList("META-INF", "services", _*) => MergeStrategy.concat
  case PathList("META-INF", _*)             => MergeStrategy.discard
  case PathList("reference.conf")           => MergeStrategy.concat
  case _                                    => MergeStrategy.first
}
`,
			},

			// project/build.properties
			{
				Path:    "project/build.properties",
				Content: `sbt.version=1.9.8` + "\n",
			},

			// project/plugins.sbt
			{
				Path: "project/plugins.sbt",
				Content: `addSbtPlugin("com.typesafe.play" % "sbt-plugin"    % "3.0.3")
addSbtPlugin("com.eed3si9n"       % "sbt-assembly"  % "2.1.5")
addSbtPlugin("org.scoverage"      % "sbt-scoverage" % "2.0.11")
`,
			},

			// app/controllers/HealthController.scala
			{
				Path: "app/controllers/HealthController.scala",
				Content: `package controllers

import play.api.libs.json.Json
import play.api.mvc.*

import javax.inject.*
import java.time.Instant
import java.util.concurrent.atomic.AtomicLong

@Singleton
class HealthController @Inject() (cc: ControllerComponents) extends AbstractController(cc) {

  private val startTime: Long = System.currentTimeMillis()

  def health: Action[AnyContent] = Action {
    val uptimeMs = System.currentTimeMillis() - startTime
    Ok(
      Json.obj(
        "status" -> "ok",
        "uptime" -> s"${uptimeMs}ms",
      )
    )
  }

  def ready: Action[AnyContent] = Action {
    // Add dependency health checks here (DB, cache, etc.)
    Ok(Json.obj("status" -> "ready"))
  }
}
`,
			},

			// app/controllers/HelloController.scala
			{
				Path: "app/controllers/HelloController.scala",
				Content: `package controllers

import models.{EchoRequest, EchoResponse, HelloResponse}
import play.api.libs.json.{JsError, Json}
import play.api.mvc.*

import javax.inject.*
import java.time.Instant
import io.prometheus.client.{Counter, Histogram}

@Singleton
class HelloController @Inject() (cc: ControllerComponents) extends AbstractController(cc) {

  private val helloCounter: Counter = Counter
    .build("hello_requests_total", "Total hello requests")
    .register()

  private val helloHistogram: Histogram = Histogram
    .build("hello_request_duration_seconds", "Hello request duration in seconds")
    .register()

  def hello(name: String = "World"): Action[AnyContent] = Action {
    val timer = helloHistogram.startTimer()
    try {
      helloCounter.inc()
      Ok(
        Json.toJson(
          HelloResponse(
            message = s"Hello, $name!",
            timestamp = Instant.now().toString,
          )
        )
      )
    } finally {
      timer.observeDuration()
    }
  }

  def echo: Action[EchoRequest] = Action(parse.json[EchoRequest]) { request =>
    Ok(
      Json.toJson(
        EchoResponse(
          echo = request.body.message,
          timestamp = Instant.now().toString,
        )
      )
    )
  }
}
`,
			},

			// app/controllers/MetricsController.scala
			{
				Path: "app/controllers/MetricsController.scala",
				Content: `package controllers

import io.prometheus.client.CollectorRegistry
import io.prometheus.client.exporter.common.TextFormat
import io.prometheus.client.hotspot.DefaultExports
import play.api.mvc.*

import javax.inject.*
import java.io.StringWriter

@Singleton
class MetricsController @Inject() (cc: ControllerComponents) extends AbstractController(cc) {

  // Register JVM metrics on startup
  DefaultExports.initialize()

  def metrics: Action[AnyContent] = Action {
    val writer = new StringWriter()
    TextFormat.write004(writer, CollectorRegistry.defaultRegistry.metricFamilySamples())
    Ok(writer.toString).as(TextFormat.CONTENT_TYPE_004)
  }
}
`,
			},

			// app/models/Dtos.scala
			{
				Path: "app/models/Dtos.scala",
				Content: `package models

import play.api.libs.json.*
import play.api.libs.functional.syntax.*

// ---- Hello ----

case class HelloResponse(
  message: String,
  timestamp: String,
)

object HelloResponse {
  implicit val writes: Writes[HelloResponse] = Json.writes[HelloResponse]
}

// ---- Echo ----

case class EchoRequest(
  message: String,
)

object EchoRequest {
  implicit val reads: Reads[EchoRequest] = Json.reads[EchoRequest]
}

case class EchoResponse(
  echo: String,
  timestamp: String,
)

object EchoResponse {
  implicit val writes: Writes[EchoResponse] = Json.writes[EchoResponse]
}
`,
			},

			// conf/routes
			{
				Path: "conf/routes",
				Content: `# Health
GET     /api/v1/health              controllers.HealthController.health
GET     /api/v1/ready               controllers.HealthController.ready

# Hello
GET     /api/v1/hello               controllers.HelloController.hello(name: String ?= "World")
POST    /api/v1/echo                controllers.HelloController.echo

# Metrics
GET     /metrics                    controllers.MetricsController.metrics
`,
			},

			// conf/application.conf
			{
				Path: "conf/application.conf",
				Content: `# Application secret — must be changed in production (min 32 chars)
play.http.secret.key = "changeme_in_production_use_32+_chars"

# Disable CSRF for REST API (add back if serving browser clients)
play.filters.disabled += "play.filters.csrf.CSRFFilter"

# Allow all hosts in development; restrict in production
play.filters.hosts {
  allowed = ["localhost", "127.0.0.1"]
}

# HTTP parser
play.http.parser.maxMemoryBuffer = 1M

# Graceful shutdown
play.server.http.idleTimeout = 75s

# Application metadata
app {
  name = "{{.ProjectName}}"
  version = "1.0.0"
  environment = "development"
}

# Thread pool configuration
pekko {
  actor {
    default-dispatcher {
      fork-join-executor {
        parallelism-min = 8
        parallelism-factor = 3.0
        parallelism-max = 64
      }
    }
  }
}
`,
			},

			// conf/logback.xml
			{
				Path: "conf/logback.xml",
				Content: `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <conversionRule conversionWord="coloredLevel"
                    converterClass="play.api.libs.logback.ColoredLevel"/>

    <appender name="STDOUT" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%date [%level] %logger{36} - %message%n%xException</pattern>
        </encoder>
    </appender>

    <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>logs/application.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <fileNamePattern>logs/application-%d{yyyy-MM-dd}.%i.log</fileNamePattern>
            <timeBasedFileNamingAndTriggeringPolicy
                    class="ch.qos.logback.core.rolling.SizeAndTimeBasedFNATP">
                <maxFileSize>100MB</maxFileSize>
            </timeBasedFileNamingAndTriggeringPolicy>
            <maxHistory>30</maxHistory>
        </rollingPolicy>
        <encoder>
            <pattern>%date [%level] %logger{36} - %message%n%xException</pattern>
        </encoder>
    </appender>

    <logger name="play" level="INFO"/>
    <logger name="application" level="DEBUG"/>
    <logger name="controllers" level="DEBUG"/>

    <root level="WARN">
        <appender-ref ref="STDOUT"/>
        <appender-ref ref="FILE"/>
    </root>
</configuration>
`,
			},

			// test/controllers/HealthControllerSpec.scala
			{
				Path: "test/controllers/HealthControllerSpec.scala",
				Content: `package controllers

import org.scalatestplus.play.*
import org.scalatestplus.play.guice.*
import play.api.test.*
import play.api.test.Helpers.*
import play.api.libs.json.*

class HealthControllerSpec extends PlaySpec with GuiceOneAppPerTest with Injecting {

  "HealthController" should {
    "return ok status on health check" in {
      val request = FakeRequest(GET, "/api/v1/health")
      val result  = route(app, request).value

      status(result) mustBe OK
      val json = contentAsJson(result)
      (json \ "status").as[String] mustBe "ok"
      (json \ "uptime").isDefined mustBe true
    }

    "return ready status on readiness check" in {
      val request = FakeRequest(GET, "/api/v1/ready")
      val result  = route(app, request).value

      status(result) mustBe OK
      val json = contentAsJson(result)
      (json \ "status").as[String] mustBe "ready"
    }
  }
}
`,
			},

			// test/controllers/HelloControllerSpec.scala
			{
				Path: "test/controllers/HelloControllerSpec.scala",
				Content: `package controllers

import org.scalatestplus.play.*
import org.scalatestplus.play.guice.*
import play.api.test.*
import play.api.test.Helpers.*
import play.api.libs.json.*

class HelloControllerSpec extends PlaySpec with GuiceOneAppPerTest with Injecting {

  "HelloController GET /api/v1/hello" should {
    "return greeting with default name" in {
      val request = FakeRequest(GET, "/api/v1/hello")
      val result  = route(app, request).value

      status(result) mustBe OK
      val json = contentAsJson(result)
      (json \ "message").as[String] mustBe "Hello, World!"
      (json \ "timestamp").isDefined mustBe true
    }

    "return greeting with custom name" in {
      val request = FakeRequest(GET, "/api/v1/hello?name=Scala")
      val result  = route(app, request).value

      status(result) mustBe OK
      (contentAsJson(result) \ "message").as[String] mustBe "Hello, Scala!"
    }
  }

  "HelloController POST /api/v1/echo" should {
    "echo back the message" in {
      val request = FakeRequest(POST, "/api/v1/echo")
        .withJsonBody(Json.obj("message" -> "hello scala"))

      val result = route(app, request).value

      status(result) mustBe OK
      val json = contentAsJson(result)
      (json \ "echo").as[String] mustBe "hello scala"
      (json \ "timestamp").isDefined mustBe true
    }

    "return 400 for invalid JSON" in {
      val request = FakeRequest(POST, "/api/v1/echo")
        .withJsonBody(Json.obj("wrong_field" -> "value"))

      val result = route(app, request).value
      status(result) mustBe BAD_REQUEST
    }
  }
}
`,
			},

			// Dockerfile
			{
				Path: "Dockerfile",
				Content: `# Build stage
FROM sbtscala/scala-sbt:eclipse-temurin-17.0.10_7_1.9.8_3.3.1 AS builder

WORKDIR /build

# Cache dependencies
COPY build.sbt project/ ./
RUN sbt update

# Build application
COPY app conf public ./
COPY app ./app
COPY conf ./conf
RUN sbt dist

# Unzip distribution
RUN unzip target/universal/*.zip -d /dist && \
    mv /dist/$(ls /dist) /dist/app

# Run stage
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /dist/app /app
RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 9000

HEALTHCHECK --interval=30s --timeout=3s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9000/api/v1/health || exit 1

ENTRYPOINT ["/app/bin/{{.ProjectName}}", "-Dhttp.port=9000"]
`,
			},

			// Makefile
			{
				Path: "Makefile",
				Content: `.PHONY: help compile test run package docker-build docker-run clean

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

compile: ## 编译代码
	sbt compile

test: ## 运行测试
	sbt test

test-coverage: ## 生成测试覆盖率报告
	sbt coverage test coverageReport

run: ## 运行开发服务器 (端口 9000)
	sbt run

package: ## 打包为分发包
	sbt dist

docker-build: ## 构建 Docker 镜像
	docker build -t {{.ProjectName}}:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 9000:9000 \
		-e APPLICATION_SECRET="changeme_in_production_use_32+_chars" \
		--name {{.ProjectName}} {{.ProjectName}}:latest

clean: ## 清理构建产物
	sbt clean
	rm -rf logs/ target/

.DEFAULT_GOAL := help
`,
			},

			// .github/workflows/ci.yml
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

      - name: 设置 JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: 缓存 sbt
        uses: actions/cache@v3
        with:
          path: |
            ~/.sbt
            ~/.ivy2/cache
            ~/.coursier
          key: sbt-${{ hashFiles('**/build.sbt', '**/plugins.sbt') }}
          restore-keys: sbt-

      - name: 运行测试
        run: sbt test

      - name: 生成覆盖率报告
        run: sbt coverage test coverageReport

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3

  build:
    name: 构建
    runs-on: ubuntu-latest
    needs: [test]

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: 缓存 sbt
        uses: actions/cache@v3
        with:
          path: |
            ~/.sbt
            ~/.ivy2/cache
            ~/.coursier
          key: sbt-${{ hashFiles('**/build.sbt', '**/plugins.sbt') }}
          restore-keys: sbt-

      - name: 打包应用
        run: sbt dist

      - name: 构建 Docker 镜像
        run: docker build -t {{.ProjectName}}:latest .
`,
			},

			// .gitignore
			{
				Path: ".gitignore",
				Content: `# sbt
target/
project/target/
project/project/
.bsp/

# IDE
.idea/
*.iml
.vscode/
*.swp
*.swo
*~

# Play
logs/
RUNNING_PID

# OS
.DS_Store
Thumbs.db

# Environment
.env
.env.local
`,
			},

			// README.md
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

Scala Play Framework 服务 - 使用 dev-template 生成

## 特性

- ✅ Scala 3.3 + Play Framework 3.0
- ✅ play-json 类型安全 JSON
- ✅ Prometheus metrics (JVM + 自定义)
- ✅ 结构化 Logback 日志
- ✅ ScalaTest + scalatestplus-play
- ✅ sbt 构建
- ✅ Docker 支持
- ✅ GitHub Actions CI/CD

## 快速开始

### 前置要求

- JDK 17+
- sbt 1.9+

### 编译项目

` + "```bash" + `
sbt compile
` + "```" + `

### 运行服务

` + "```bash" + `
sbt run
` + "```" + `

服务将在 http://localhost:9000 启动

### 运行测试

` + "```bash" + `
make test
` + "```" + `

### 运行测试并生成覆盖率报告

` + "```bash" + `
make test-coverage
` + "```" + `

## API 端点

- ` + "`GET /api/v1/health`" + ` - 健康检查
- ` + "`GET /api/v1/ready`" + ` - 就绪检查
- ` + "`GET /metrics`" + ` - Prometheus metrics
- ` + "`GET /api/v1/hello?name=xxx`" + ` - 示例 API
- ` + "`POST /api/v1/echo`" + ` - Echo API

## 生产配置

修改 ` + "`conf/application.conf`" + ` 中的配置：

` + "```bash" + `
# 设置应用 Secret (必须修改！)
play.http.secret.key = "your-strong-secret-here-minimum-32-chars"

# 限制允许的 Host
play.filters.hosts.allowed = ["yourdomain.com"]
` + "```" + `

## Docker

` + "```bash" + `
make docker-build
make docker-run
` + "```" + `

## 项目结构

` + "```" + `
.
├── app/
│   ├── controllers/
│   │   ├── HealthController.scala    # 健康检查
│   │   ├── HelloController.scala     # 示例控制器
│   │   └── MetricsController.scala  # Prometheus metrics
│   └── models/Dtos.scala            # 数据模型
├── conf/
│   ├── routes                       # 路由配置
│   ├── application.conf             # 应用配置 (HOCON)
│   └── logback.xml                  # 日志配置
├── test/controllers/                # 测试
├── project/
│   ├── plugins.sbt                  # sbt 插件
│   └── build.properties             # sbt 版本
├── build.sbt                        # 构建定义
├── Dockerfile
└── Makefile
` + "```" + `

## License

MIT
`,
			},
		},
	}
}
