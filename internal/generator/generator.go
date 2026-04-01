package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Generator 项目生成器
type Generator struct {
	templates map[string]*Template
}

// Template 模板定义
type Template struct {
	Name        string
	Description string
	Features    string
	Files       []FileTemplate
}

// FileTemplate 文件模板
type FileTemplate struct {
	Path    string // 相对路径
	Content string // 文件内容
	IsDir   bool   // 是否是目录
}

// TemplateData 模板数据
type TemplateData struct {
	ProjectName string
	ModulePath  string
}

// New 创建新的生成器
func New() *Generator {
	return &Generator{
		templates: initTemplates(),
	}
}

// IsValidTemplate 检查模板类型是否有效
func IsValidTemplate(templateType string) bool {
	gen := New()
	_, exists := gen.templates[templateType]
	return exists
}

// ListTemplates 列出所有模板
func ListTemplates() []*Template {
	gen := New()
	var templates []*Template
	for _, tmpl := range gen.templates {
		templates = append(templates, tmpl)
	}
	return templates
}

// Generate 生成项目
func (g *Generator) Generate(templateType, projectName, outputPath string) error {
	tmpl, exists := g.templates[templateType]
	if !exists {
		return fmt.Errorf("模板不存在: %s", templateType)
	}

	// 准备模板数据
	data := TemplateData{
		ProjectName: projectName,
		ModulePath:  fmt.Sprintf("github.com/yourname/%s", projectName),
	}

	// 创建根目录
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("创建项目目录失败: %w", err)
	}

	// 生成所有文件
	for _, fileTemplate := range tmpl.Files {
		filePath := filepath.Join(outputPath, fileTemplate.Path)

		// 如果是目录，创建目录
		if fileTemplate.IsDir {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return fmt.Errorf("创建目录失败 %s: %w", filePath, err)
			}
			continue
		}

		// 确保父目录存在
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %w", dir, err)
		}

		// 渲染模板内容
		content, err := renderTemplate(fileTemplate.Content, data)
		if err != nil {
			return fmt.Errorf("渲染模板失败 %s: %w", filePath, err)
		}

		// 写入文件
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件失败 %s: %w", filePath, err)
		}

		fmt.Printf("  ✓ %s\n", fileTemplate.Path)
	}

	return nil
}

// renderTemplate 渲染模板内容
func renderTemplate(content string, data TemplateData) (string, error) {
	tmpl, err := template.New("file").Parse(content)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// initTemplates 初始化所有模板
func initTemplates() map[string]*Template {
	return map[string]*Template{
		"go-service":     GetGoServiceTemplate(),
		"python-api":     GetPythonAPITemplate(),
		"java-service":   GetJavaServiceTemplate(),
		"kotlin-service": GetKotlinServiceTemplate(),
		"scala-service":  GetScalaServiceTemplate(),
	}
}
