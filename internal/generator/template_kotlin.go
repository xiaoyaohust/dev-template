package generator

// GetKotlinServiceTemplate 返回 Kotlin Spring Boot 服务模板
func GetKotlinServiceTemplate() *Template {
	return &Template{
		Name:        "kotlin-service",
		Description: "Kotlin Spring Boot 服务模板 - 现代 JVM 微服务",
		Features:    "Spring Boot 3, Kotlin Coroutines, Gradle Kotlin DSL, ktlint, Micrometer Prometheus, Docker, K8s",
		Files: []FileTemplate{
			// build.gradle.kts
			{
				Path: "build.gradle.kts",
				Content: `import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.springframework.boot") version "3.2.1"
    id("io.spring.dependency-management") version "1.1.4"
    kotlin("jvm") version "1.9.21"
    kotlin("plugin.spring") version "1.9.21"
    id("org.jlleitschuh.gradle.ktlint") version "12.1.0"
    jacoco
}

group = "com.example"
version = "1.0.0"
java.sourceCompatibility = JavaVersion.VERSION_17

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("org.springframework.boot:spring-boot-starter-actuator")
    implementation("org.springframework.boot:spring-boot-starter-validation")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
    implementation("io.micrometer:micrometer-registry-prometheus")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-reactor")
    implementation("io.projectreactor.kotlin:reactor-kotlin-extensions")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
    testImplementation("org.jetbrains.kotlinx:kotlinx-coroutines-test")
}

tasks.withType<KotlinCompile> {
    kotlinOptions {
        freeCompilerArgs += "-Xjsr305=strict"
        jvmTarget = "17"
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
}

tasks.test {
    finalizedBy(tasks.jacocoTestReport)
}

tasks.jacocoTestReport {
    reports {
        xml.required = true
        html.required = true
    }
}
`,
			},

			// settings.gradle.kts
			{
				Path:    "settings.gradle.kts",
				Content: `rootProject.name = "{{.ProjectName}}"` + "\n",
			},

			// Gradle wrapper properties
			{
				Path: "gradle/wrapper/gradle-wrapper.properties",
				Content: `distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-8.5-bin.zip
networkTimeout=10000
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
`,
			},

			// Application.kt
			{
				Path: "src/main/kotlin/com/example/Application.kt",
				Content: `package com.example

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class Application

fun main(args: Array<String>) {
    runApplication<Application>(*args)
}
`,
			},

			// AppConfig.kt
			{
				Path: "src/main/kotlin/com/example/config/AppConfig.kt",
				Content: `package com.example.config

import org.springframework.boot.context.properties.ConfigurationProperties
import org.springframework.boot.context.properties.EnableConfigurationProperties
import org.springframework.context.annotation.Configuration

@Configuration
@EnableConfigurationProperties(AppConfig::class)
@ConfigurationProperties(prefix = "app")
data class AppConfig(
    var name: String = "{{.ProjectName}}",
    var version: String = "1.0.0",
    var environment: String = "development",
)
`,
			},

			// HealthController.kt
			{
				Path: "src/main/kotlin/com/example/controller/HealthController.kt",
				Content: `package com.example.controller

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import java.time.Duration

@RestController
@RequestMapping("/api/v1")
class HealthController {
    private val startTime = System.currentTimeMillis()

    @GetMapping("/health")
    fun health(): Map<String, Any> {
        val uptime = Duration.ofMillis(System.currentTimeMillis() - startTime)
        return mapOf(
            "status" to "ok",
            "uptime" to uptime.toString(),
        )
    }

    @GetMapping("/ready")
    fun ready(): Map<String, String> = mapOf("status" to "ready")
}
`,
			},

			// HelloController.kt
			{
				Path: "src/main/kotlin/com/example/controller/HelloController.kt",
				Content: `package com.example.controller

import com.example.dto.EchoRequest
import com.example.dto.EchoResponse
import com.example.dto.HelloResponse
import io.micrometer.core.instrument.Counter
import io.micrometer.core.instrument.MeterRegistry
import io.micrometer.core.instrument.Timer
import jakarta.validation.Valid
import org.slf4j.LoggerFactory
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestBody
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.bind.annotation.RestController
import java.time.Instant

@RestController
@RequestMapping("/api/v1")
class HelloController(meterRegistry: MeterRegistry) {
    private val logger = LoggerFactory.getLogger(javaClass)

    private val helloCounter: Counter = Counter.builder("hello.requests.total")
        .description("Total hello requests")
        .register(meterRegistry)

    private val helloTimer: Timer = Timer.builder("hello.requests.duration")
        .description("Hello request duration in seconds")
        .register(meterRegistry)

    @GetMapping("/hello")
    fun hello(
        @RequestParam(defaultValue = "World") name: String,
    ): ResponseEntity<HelloResponse> =
        helloTimer.recordCallable {
            helloCounter.increment()
            logger.info("Processing hello request: name={}", name)
            ResponseEntity.ok(
                HelloResponse(
                    message = "Hello, $name!",
                    timestamp = Instant.now().toString(),
                ),
            )
        }!!

    @PostMapping("/echo")
    fun echo(
        @Valid @RequestBody request: EchoRequest,
    ): ResponseEntity<EchoResponse> {
        logger.info("Processing echo request: message={}", request.message)
        return ResponseEntity.ok(
            EchoResponse(
                echo = request.message,
                timestamp = Instant.now().toString(),
            ),
        )
    }
}
`,
			},

			// Dtos.kt
			{
				Path: "src/main/kotlin/com/example/dto/Dtos.kt",
				Content: `package com.example.dto

import jakarta.validation.constraints.NotBlank
import jakarta.validation.constraints.Size

data class HelloResponse(
    val message: String,
    val timestamp: String,
)

data class EchoRequest(
    @field:NotBlank(message = "Message cannot be blank")
    @field:Size(min = 1, max = 1000, message = "Message must be between 1 and 1000 characters")
    val message: String,
)

data class EchoResponse(
    val echo: String,
    val timestamp: String,
)
`,
			},

			// application.yml
			{
				Path: "src/main/resources/application.yml",
				Content: `spring:
  application:
    name: {{.ProjectName}}

server:
  port: 8080
  shutdown: graceful

app:
  name: {{.ProjectName}}
  version: 1.0.0
  environment: development

management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus
  endpoint:
    health:
      show-details: always
  metrics:
    export:
      prometheus:
        enabled: true

logging:
  level:
    root: INFO
    com.example: DEBUG
`,
			},

			// logback-spring.xml
			{
				Path: "src/main/resources/logback-spring.xml",
				Content: `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <include resource="org/springframework/boot/logging/logback/defaults.xml"/>

    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
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
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <root level="INFO">
        <appender-ref ref="CONSOLE"/>
        <appender-ref ref="FILE"/>
    </root>
</configuration>
`,
			},

			// ApplicationTests.kt
			{
				Path: "src/test/kotlin/com/example/ApplicationTests.kt",
				Content: `package com.example

import org.junit.jupiter.api.Test
import org.springframework.boot.test.context.SpringBootTest

@SpringBootTest
class ApplicationTests {

    @Test
    fun contextLoads() {
        // Verify the application context loads successfully
    }
}
`,
			},

			// HelloControllerTest.kt
			{
				Path: "src/test/kotlin/com/example/controller/HelloControllerTest.kt",
				Content: `package com.example.controller

import com.example.dto.EchoRequest
import com.fasterxml.jackson.module.kotlin.jacksonObjectMapper
import org.junit.jupiter.api.Test
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc
import org.springframework.boot.test.context.SpringBootTest
import org.springframework.http.MediaType
import org.springframework.test.web.servlet.MockMvc
import org.springframework.test.web.servlet.get
import org.springframework.test.web.servlet.post

@SpringBootTest
@AutoConfigureMockMvc
class HelloControllerTest {

    @Autowired
    private lateinit var mockMvc: MockMvc

    private val objectMapper = jacksonObjectMapper()

    @Test
    fun ` + "`hello returns greeting with default name`" + `() {
        mockMvc.get("/api/v1/hello")
            .andExpect {
                status { isOk() }
                jsonPath("$.message") { value("Hello, World!") }
                jsonPath("$.timestamp") { exists() }
            }
    }

    @Test
    fun ` + "`hello returns greeting with custom name`" + `() {
        mockMvc.get("/api/v1/hello?name=Kotlin")
            .andExpect {
                status { isOk() }
                jsonPath("$.message") { value("Hello, Kotlin!") }
            }
    }

    @Test
    fun ` + "`echo returns message`" + `() {
        val request = EchoRequest(message = "test message")
        mockMvc.post("/api/v1/echo") {
            contentType = MediaType.APPLICATION_JSON
            content = objectMapper.writeValueAsString(request)
        }.andExpect {
            status { isOk() }
            jsonPath("$.echo") { value("test message") }
            jsonPath("$.timestamp") { exists() }
        }
    }

    @Test
    fun ` + "`echo rejects blank message`" + `() {
        val request = mapOf("message" to "")
        mockMvc.post("/api/v1/echo") {
            contentType = MediaType.APPLICATION_JSON
            content = objectMapper.writeValueAsString(request)
        }.andExpect {
            status { isBadRequest() }
        }
    }
}
`,
			},

			// k8s/deployment.yaml
			{
				Path: "k8s/deployment.yaml",
				Content: `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ProjectName}}
  labels:
    app: {{.ProjectName}}
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.ProjectName}}
  template:
    metadata:
      labels:
        app: {{.ProjectName}}
    spec:
      containers:
      - name: {{.ProjectName}}
        image: {{.ProjectName}}:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: SPRING_PROFILES_ACTIVE
          value: "production"
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /actuator/health/liveness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /actuator/health/readiness
            port: 8080
          initialDelaySeconds: 20
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ProjectName}}
  labels:
    app: {{.ProjectName}}
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
    protocol: TCP
    name: http
  selector:
    app: {{.ProjectName}}
`,
			},

			// Dockerfile
			{
				Path: "Dockerfile",
				Content: `# Build stage
FROM gradle:8.5-jdk17 AS builder

WORKDIR /build

COPY build.gradle.kts settings.gradle.kts ./
COPY gradle gradle
RUN gradle dependencies --no-daemon || true

COPY src ./src
RUN gradle bootJar --no-daemon

# Run stage
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /build/build/libs/*.jar app.jar
RUN chown appuser:appgroup app.jar

USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/actuator/health || exit 1

ENTRYPOINT ["java", "-jar", "/app/app.jar"]
`,
			},

			// Makefile
			{
				Path: "Makefile",
				Content: `.PHONY: help build test run lint docker-build docker-run clean

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## 构建应用
	./gradlew build

test: ## 运行测试
	./gradlew test

test-coverage: ## 运行测试并生成覆盖率报告
	./gradlew test jacocoTestReport

run: ## 运行应用
	./gradlew bootRun

lint: ## 运行 ktlint 检查
	./gradlew ktlintCheck

lint-fix: ## 自动修复 ktlint 问题
	./gradlew ktlintFormat

docker-build: ## 构建 Docker 镜像
	docker build -t {{.ProjectName}}:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 --name {{.ProjectName}} {{.ProjectName}}:latest

clean: ## 清理构建产物
	./gradlew clean
	rm -rf logs/

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

      - name: 设置 Gradle
        uses: gradle/gradle-build-action@v2

      - name: 运行测试
        run: ./gradlew test

      - name: 生成覆盖率报告
        run: ./gradlew jacocoTestReport

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        with:
          files: build/reports/jacoco/test/jacocoTestReport.xml

  lint:
    name: 代码检查
    runs-on: ubuntu-latest

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: 设置 Gradle
        uses: gradle/gradle-build-action@v2

      - name: 运行 ktlint
        run: ./gradlew ktlintCheck

  build:
    name: 构建
    runs-on: ubuntu-latest
    needs: [test, lint]

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'

      - name: 设置 Gradle
        uses: gradle/gradle-build-action@v2

      - name: 构建 JAR
        run: ./gradlew bootJar

      - name: 构建 Docker 镜像
        run: docker build -t {{.ProjectName}}:latest .
`,
			},

			// .gitignore
			{
				Path: ".gitignore",
				Content: `# Gradle
.gradle/
build/
!gradle/wrapper/gradle-wrapper.jar

# IDE
.idea/
*.iml
.vscode/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
logs/
*.log

# Environment
.env
.env.local
`,
			},

			// README.md
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

Kotlin Spring Boot 服务 - 使用 dev-template 生成

## 特性

- ✅ Spring Boot 3.2 + Kotlin 1.9
- ✅ Kotlin Coroutines (reactor-kotlin-extensions)
- ✅ Gradle Kotlin DSL
- ✅ ktlint 代码规范
- ✅ Spring Actuator 健康检查
- ✅ Micrometer + Prometheus metrics
- ✅ JUnit 5 + Spring Boot Test (Kotlin DSL)
- ✅ JaCoCo 覆盖率
- ✅ Docker 支持
- ✅ Kubernetes 配置
- ✅ GitHub Actions CI/CD

## 快速开始

### 前置要求

- JDK 17+
- Gradle 8.5+ (或使用项目自带的 gradlew)

### 初始化 Gradle Wrapper (首次使用)

` + "```bash" + `
gradle wrapper
` + "```" + `

### 构建项目

` + "```bash" + `
./gradlew build
` + "```" + `

### 运行服务

` + "```bash" + `
./gradlew bootRun
` + "```" + `

服务将在 http://localhost:8080 启动

### 运行测试

` + "```bash" + `
make test
` + "```" + `

### 代码检查

` + "```bash" + `
make lint      # 检查
make lint-fix  # 自动修复
` + "```" + `

## API 端点

- ` + "`GET /api/v1/health`" + ` - 自定义健康检查
- ` + "`GET /api/v1/ready`" + ` - 就绪检查
- ` + "`GET /actuator/health`" + ` - Spring Actuator 健康检查
- ` + "`GET /actuator/prometheus`" + ` - Prometheus metrics
- ` + "`GET /api/v1/hello?name=xxx`" + ` - 示例 API
- ` + "`POST /api/v1/echo`" + ` - Echo API

## Docker

` + "```bash" + `
make docker-build
make docker-run
` + "```" + `

## Kubernetes

` + "```bash" + `
kubectl apply -f k8s/deployment.yaml
` + "```" + `

## 项目结构

` + "```" + `
.
├── build.gradle.kts                   # Gradle 构建配置
├── settings.gradle.kts                # 项目设置
├── src/
│   ├── main/kotlin/com/example/
│   │   ├── Application.kt             # 入口
│   │   ├── config/AppConfig.kt        # 配置
│   │   ├── controller/               # 控制器
│   │   └── dto/Dtos.kt               # 数据类
│   ├── main/resources/
│   │   ├── application.yml           # 配置文件
│   │   └── logback-spring.xml        # 日志配置
│   └── test/kotlin/                  # 测试
├── k8s/                              # Kubernetes 配置
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
