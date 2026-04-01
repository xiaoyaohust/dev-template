package generator

// GetJavaServiceTemplate 返回 Java Spring Boot 服务模板
func GetJavaServiceTemplate() *Template {
	return &Template{
		Name:        "java-service",
		Description: "Java Spring Boot 服务模板 - 企业级微服务",
		Features:    "Spring Boot 3, JUnit 5, Checkstyle, Logback, Prometheus, Docker, K8s",
		Files: []FileTemplate{
			// pom.xml
			{
				Path: "pom.xml",
				Content: `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0
         https://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
        <version>3.2.1</version>
        <relativePath/>
    </parent>

    <groupId>com.example</groupId>
    <artifactId>{{.ProjectName}}</artifactId>
    <version>1.0.0</version>
    <name>{{.ProjectName}}</name>
    <description>Spring Boot 服务</description>

    <properties>
        <java.version>17</java.version>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    </properties>

    <dependencies>
        <!-- Spring Boot Starter Web -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>

        <!-- Spring Boot Starter Actuator -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-actuator</artifactId>
        </dependency>

        <!-- Spring Boot Starter Validation -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-validation</artifactId>
        </dependency>

        <!-- Micrometer Prometheus -->
        <dependency>
            <groupId>io.micrometer</groupId>
            <artifactId>micrometer-registry-prometheus</artifactId>
        </dependency>

        <!-- Lombok -->
        <dependency>
            <groupId>org.projectlombok</groupId>
            <artifactId>lombok</artifactId>
            <optional>true</optional>
        </dependency>

        <!-- Spring Boot Starter Test -->
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-test</artifactId>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <!-- Spring Boot Maven Plugin -->
            <plugin>
                <groupId>org.springframework.boot</groupId>
                <artifactId>spring-boot-maven-plugin</artifactId>
                <configuration>
                    <excludes>
                        <exclude>
                            <groupId>org.projectlombok</groupId>
                            <artifactId>lombok</artifactId>
                        </exclude>
                    </excludes>
                </configuration>
            </plugin>

            <!-- Maven Checkstyle Plugin -->
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-checkstyle-plugin</artifactId>
                <version>3.3.1</version>
                <configuration>
                    <configLocation>checkstyle.xml</configLocation>
                    <consoleOutput>true</consoleOutput>
                    <failsOnError>true</failsOnError>
                </configuration>
                <executions>
                    <execution>
                        <phase>validate</phase>
                        <goals>
                            <goal>check</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>

            <!-- Maven Surefire Plugin (Tests) -->
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>3.2.3</version>
            </plugin>

            <!-- JaCoCo (Coverage) -->
            <plugin>
                <groupId>org.jacoco</groupId>
                <artifactId>jacoco-maven-plugin</artifactId>
                <version>0.8.11</version>
                <executions>
                    <execution>
                        <goals>
                            <goal>prepare-agent</goal>
                        </goals>
                    </execution>
                    <execution>
                        <id>report</id>
                        <phase>test</phase>
                        <goals>
                            <goal>report</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
        </plugins>
    </build>
</project>
`,
			},

			// 主应用类
			{
				Path: "src/main/java/com/example/Application.java",
				Content: `package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * 应用程序主类
 */
@SpringBootApplication
public class Application {

    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
`,
			},

			// 配置类
			{
				Path: "src/main/java/com/example/config/AppConfig.java",
				Content: `package com.example.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;
import lombok.Data;

/**
 * 应用配置
 */
@Data
@Configuration
@ConfigurationProperties(prefix = "app")
public class AppConfig {

    private String name = "{{.ProjectName}}";
    private String version = "1.0.0";
    private String environment = "development";
}
`,
			},

			// Hello Controller
			{
				Path: "src/main/java/com/example/controller/HelloController.java",
				Content: `package com.example.controller;

import com.example.dto.HelloResponse;
import com.example.dto.EchoRequest;
import com.example.dto.EchoResponse;
import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import jakarta.validation.Valid;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.Instant;

/**
 * Hello 控制器
 */
@Slf4j
@RestController
@RequestMapping("/api/v1")
public class HelloController {

    private final Counter helloCounter;
    private final Timer helloTimer;

    public HelloController(MeterRegistry meterRegistry) {
        this.helloCounter = Counter.builder("hello.requests.total")
                .description("Hello 请求总数")
                .register(meterRegistry);

        this.helloTimer = Timer.builder("hello.requests.duration")
                .description("Hello 请求延迟")
                .register(meterRegistry);
    }

    /**
     * Hello 端点
     */
    @GetMapping("/hello")
    public ResponseEntity<HelloResponse> hello(
            @RequestParam(defaultValue = "世界") String name) {

        return helloTimer.recordCallable(() -> {
            helloCounter.increment();
            log.info("处理 Hello 请求: name={}", name);

            HelloResponse response = new HelloResponse();
            response.setMessage("你好, " + name + "!");
            response.setTimestamp(Instant.now().toString());

            return ResponseEntity.ok(response);
        });
    }

    /**
     * Echo 端点
     */
    @PostMapping("/echo")
    public ResponseEntity<EchoResponse> echo(@Valid @RequestBody EchoRequest request) {
        log.info("处理 Echo 请求: message={}", request.getMessage());

        EchoResponse response = new EchoResponse();
        response.setEcho(request.getMessage());
        response.setTimestamp(Instant.now().toString());

        return ResponseEntity.ok(response);
    }
}
`,
			},

			// DTO 类
			{
				Path: "src/main/java/com/example/dto/HelloResponse.java",
				Content: `package com.example.dto;

import lombok.Data;

/**
 * Hello 响应
 */
@Data
public class HelloResponse {
    private String message;
    private String timestamp;
}
`,
			},

			{
				Path: "src/main/java/com/example/dto/EchoRequest.java",
				Content: `package com.example.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;
import lombok.Data;

/**
 * Echo 请求
 */
@Data
public class EchoRequest {

    @NotBlank(message = "消息不能为空")
    @Size(min = 1, max = 1000, message = "消息长度必须在 1 到 1000 之间")
    private String message;
}
`,
			},

			{
				Path: "src/main/java/com/example/dto/EchoResponse.java",
				Content: `package com.example.dto;

import lombok.Data;

/**
 * Echo 响应
 */
@Data
public class EchoResponse {
    private String echo;
    private String timestamp;
}
`,
			},

			// 测试类
			{
				Path: "src/test/java/com/example/ApplicationTests.java",
				Content: `package com.example;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

/**
 * 应用启动测试
 */
@SpringBootTest
class ApplicationTests {

    @Test
    void contextLoads() {
        // 测试应用上下文加载
    }
}
`,
			},

			{
				Path: "src/test/java/com/example/controller/HelloControllerTest.java",
				Content: `package com.example.controller;

import com.example.dto.EchoRequest;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.*;

/**
 * Hello 控制器测试
 */
@SpringBootTest
@AutoConfigureMockMvc
class HelloControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @Test
    void testHelloDefault() throws Exception {
        mockMvc.perform(get("/api/v1/hello"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.message").value("你好, 世界!"))
                .andExpect(jsonPath("$.timestamp").exists());
    }

    @Test
    void testHelloWithName() throws Exception {
        mockMvc.perform(get("/api/v1/hello")
                        .param("name", "Java"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.message").value("你好, Java!"));
    }

    @Test
    void testEcho() throws Exception {
        EchoRequest request = new EchoRequest();
        request.setMessage("测试消息");

        mockMvc.perform(post("/api/v1/echo")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.echo").value("测试消息"))
                .andExpect(jsonPath("$.timestamp").exists());
    }

    @Test
    void testEchoValidation() throws Exception {
        EchoRequest request = new EchoRequest();
        request.setMessage(""); // 空消息

        mockMvc.perform(post("/api/v1/echo")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }
}
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

# Actuator 配置
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

# 日志配置
logging:
  level:
    root: INFO
    com.example: DEBUG
  pattern:
    console: "%d{yyyy-MM-dd HH:mm:ss} - %msg%n"
    file: "%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
`,
			},

			// logback-spring.xml
			{
				Path: "src/main/resources/logback-spring.xml",
				Content: `<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <include resource="org/springframework/boot/logging/logback/defaults.xml"/>

    <!-- 控制台输出 -->
    <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
        <encoder>
            <pattern>%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
        </encoder>
    </appender>

    <!-- 文件输出 -->
    <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
        <file>logs/application.log</file>
        <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
            <fileNamePattern>logs/application-%d{yyyy-MM-dd}.%i.log</fileNamePattern>
            <timeBasedFileNamingAndTriggeringPolicy class="ch.qos.logback.core.rolling.SizeAndTimeBasedFNATP">
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

			// Dockerfile
			{
				Path: "Dockerfile",
				Content: `# 构建阶段
FROM maven:3.9-eclipse-temurin-17 AS builder

WORKDIR /build

# 复制 pom.xml 并下载依赖（利用 Docker 缓存）
COPY pom.xml .
RUN mvn dependency:go-offline

# 复制源代码并构建
COPY src ./src
RUN mvn clean package -DskipTests

# 运行阶段
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# 创建非 root 用户
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# 从构建阶段复制 jar 文件
COPY --from=builder /build/target/*.jar app.jar

# 设置文件权限
RUN chown appuser:appgroup app.jar

# 切换到非 root 用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/actuator/health || exit 1

# 运行应用
ENTRYPOINT ["java", "-jar", "/app/app.jar"]
`,
			},

			// checkstyle.xml
			{
				Path: "checkstyle.xml",
				Content: `<?xml version="1.0"?>
<!DOCTYPE module PUBLIC
    "-//Checkstyle//DTD Checkstyle Configuration 1.3//EN"
    "https://checkstyle.org/dtds/configuration_1_3.dtd">

<module name="Checker">
    <property name="charset" value="UTF-8"/>
    <property name="severity" value="warning"/>
    <property name="fileExtensions" value="java, properties, xml"/>

    <!-- 文件长度检查 -->
    <module name="FileLength">
        <property name="max" value="500"/>
    </module>

    <!-- TreeWalker 检查 -->
    <module name="TreeWalker">
        <!-- 命名规范 -->
        <module name="ConstantName"/>
        <module name="LocalFinalVariableName"/>
        <module name="LocalVariableName"/>
        <module name="MemberName"/>
        <module name="MethodName"/>
        <module name="PackageName"/>
        <module name="ParameterName"/>
        <module name="StaticVariableName"/>
        <module name="TypeName"/>

        <!-- 导入检查 -->
        <module name="AvoidStarImport"/>
        <module name="IllegalImport"/>
        <module name="RedundantImport"/>
        <module name="UnusedImports"/>

        <!-- 空格检查 -->
        <module name="WhitespaceAfter"/>
        <module name="WhitespaceAround"/>

        <!-- 块检查 -->
        <module name="EmptyBlock"/>
        <module name="LeftCurly"/>
        <module name="NeedBraces"/>
        <module name="RightCurly"/>

        <!-- 编码检查 -->
        <module name="EmptyStatement"/>
        <module name="EqualsHashCode"/>
        <module name="SimplifyBooleanExpression"/>
        <module name="SimplifyBooleanReturn"/>

        <!-- 类设计检查 -->
        <module name="FinalClass"/>
        <module name="InterfaceIsType"/>
        <module name="VisibilityModifier"/>

        <!-- 其他检查 -->
        <module name="ArrayTypeStyle"/>
        <module name="TodoComment"/>
        <module name="UpperEll"/>
    </module>
</module>
`,
			},

			// Makefile
			{
				Path: "Makefile",
				Content: `.PHONY: help clean compile test package run docker lint

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## 清理构建文件
	./mvnw clean

compile: ## 编译代码
	./mvnw compile

test: ## 运行测试
	./mvnw test

test-coverage: ## 生成测试覆盖率报告
	./mvnw test jacoco:report

package: ## 打包应用
	./mvnw clean package

run: ## 运行应用
	./mvnw spring-boot:run

lint: ## 代码检查
	./mvnw checkstyle:check

docker-build: ## 构建 Docker 镜像
	docker build -t {{.ProjectName}}:latest .

docker-run: ## 运行 Docker 容器
	docker run -p 8080:8080 --name {{.ProjectName}} {{.ProjectName}}:latest

.DEFAULT_GOAL := help
`,
			},

			// Kubernetes 部署配置
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

    steps:
      - name: Checkout 代码
        uses: actions/checkout@v4

      - name: 设置 JDK 17
        uses: actions/setup-java@v4
        with:
          java-version: '17'
          distribution: 'temurin'
          cache: maven

      - name: 运行测试
        run: ./mvnw clean test

      - name: 生成覆盖率报告
        run: ./mvnw jacoco:report

      - name: 上传覆盖率
        uses: codecov/codecov-action@v3
        with:
          files: ./target/site/jacoco/jacoco.xml

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
          cache: maven

      - name: Checkstyle 检查
        run: ./mvnw checkstyle:check

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
          cache: maven

      - name: 构建应用
        run: ./mvnw clean package -DskipTests

      - name: 构建 Docker 镜像
        run: docker build -t {{.ProjectName}}:latest .
`,
			},

			// .gitignore
			{
				Path: ".gitignore",
				Content: `# Maven
target/
pom.xml.tag
pom.xml.releaseBackup
pom.xml.versionsBackup
pom.xml.next
release.properties
dependency-reduced-pom.xml
buildNumber.properties
.mvn/timing.properties
.mvn/wrapper/maven-wrapper.jar

# Java
*.class
*.jar
*.war
*.ear
*.log
hs_err_pid*

# IDE
.idea/
*.iml
.vscode/
.classpath
.project
.settings/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# 日志
logs/
*.log
`,
			},

			// mvnw (Maven Wrapper)
			{
				Path: ".mvn/wrapper/maven-wrapper.properties",
				Content: `distributionUrl=https://repo.maven.apache.org/maven2/org/apache/maven/apache-maven/3.9.6/apache-maven-3.9.6-bin.zip
wrapperUrl=https://repo.maven.apache.org/maven2/org/apache/maven/wrapper/maven-wrapper/3.2.0/maven-wrapper-3.2.0.jar
`,
			},

			// README
			{
				Path: "README.md",
				Content: `# {{.ProjectName}}

Java Spring Boot 服务 - 使用 dev-template 生成

## 特性

- ✅ Spring Boot 3.2
- ✅ Spring Web (REST API)
- ✅ Spring Actuator (健康检查、metrics)
- ✅ Prometheus metrics
- ✅ Logback 日志
- ✅ JUnit 5 测试
- ✅ Checkstyle 代码检查
- ✅ JaCoCo 覆盖率
- ✅ Lombok
- ✅ Docker 支持
- ✅ Kubernetes 配置
- ✅ GitHub Actions CI/CD

## 快速开始

### 前置要求

- JDK 17+
- Maven 3.8+ (或使用项目自带的 mvnw)

### 构建项目

` + "```bash" + `
./mvnw clean install
` + "```" + `

### 运行服务

` + "```bash" + `
./mvnw spring-boot:run
` + "```" + `

服务将在 http://localhost:8080 启动

### 运行测试

` + "```bash" + `
make test
` + "```" + `

### 代码检查

` + "```bash" + `
make lint
` + "```" + `

## API 端点

- ` + "`GET /actuator/health`" + ` - 健康检查
- ` + "`GET /actuator/metrics`" + ` - Metrics
- ` + "`GET /actuator/prometheus`" + ` - Prometheus metrics
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

## Kubernetes

部署到 Kubernetes:

` + "```bash" + `
kubectl apply -f k8s/deployment.yaml
` + "```" + `

## 配置

配置文件位于 ` + "`src/main/resources/application.yml`" + `

可以通过环境变量或 Spring profiles 覆盖配置。

## 项目结构

` + "```" + `
.
├── src/
│   ├── main/
│   │   ├── java/com/example/
│   │   │   ├── Application.java       # 主类
│   │   │   ├── config/               # 配置
│   │   │   ├── controller/           # 控制器
│   │   │   └── dto/                  # DTO
│   │   └── resources/
│   │       ├── application.yml       # 配置文件
│   │       └── logback-spring.xml    # 日志配置
│   └── test/                         # 测试
├── k8s/                              # Kubernetes 配置
├── Dockerfile                        # Docker 配置
├── pom.xml                           # Maven 配置
├── checkstyle.xml                    # Checkstyle 配置
└── Makefile                          # 构建脚本
` + "```" + `

## 开发工具

### Maven Wrapper

项目包含 Maven Wrapper,无需预装 Maven:

` + "```bash" + `
./mvnw clean install    # Linux/Mac
mvnw.cmd clean install  # Windows
` + "```" + `

### 热重载

开发模式下支持热重载:

` + "```bash" + `
./mvnw spring-boot:run
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
