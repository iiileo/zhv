package converter

import (
	"fmt"
	"strings"

	"github.com/iiileo/zhv/client"
	"github.com/iiileo/zhv/config"
)

// Converter 中文转英文变量名转换器
type Converter struct {
	client *client.OpenAIClient
}

// NewConverter 创建新的转换器
func NewConverter(cfg *config.Config) *Converter {
	return &Converter{
		client: client.NewOpenAIClient(cfg),
	}
}

// ConvertToVariableName 将中文转换为英文变量名
func (c *Converter) ConvertToVariableName(chineseText string, style string) ([]string, error) {
	prompt := c.buildPrompt(chineseText, style)

	messages := []client.Message{
		{
			Role:    "system",
			Content: "你是一位资深的软件工程师和编程规范专家，精通多种编程语言的命名约定。你的任务是将中文概念准确转换为地道的英文变量名，确保：1) 语义准确表达原始概念；2) 遵循目标命名风格；3) 符合国际编程最佳实践；4) 使用简洁明了的英语词汇；5）请不要回复无关的信息，仅回复变量名，不要回复任何其他信息。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	response, err := c.client.Chat(messages)
	if err != nil {
		return nil, fmt.Errorf("AI请求失败: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("AI未返回有效响应")
	}

	content := response.Choices[0].Message.Content
	return c.parseResponse(content), nil
}

// buildPrompt 构建提示词
func (c *Converter) buildPrompt(chineseText, style string) string {
	var styleDesc, examples string
	switch style {
	case "camel":
		styleDesc = "驼峰命名法 (camelCase)"
		examples = "userName, userProfile, dataCount, isActive"
	case "pascal":
		styleDesc = "帕斯卡命名法 (PascalCase)"
		examples = "UserName, UserProfile, DataCount, IsActive"
	case "snake":
		styleDesc = "蛇形命名法 (snake_case)"
		examples = "user_name, user_profile, data_count, is_active"
	case "kebab":
		styleDesc = "短横线命名法 (kebab-case)"
		examples = "user-name, user-profile, data-count, is-active"
	default:
		styleDesc = "驼峰命名法 (camelCase)"
		examples = "userName, userProfile, dataCount, isActive"
	}

	return fmt.Sprintf(`作为专业的变量命名助手，为中文词汇"%s"生成高质量的英文变量名。

## 要求
- 命名风格: %s
- 参考示例: %s
- 生成3-5个选项
- 使用地道英语，避免中式英语
- 符合编程最佳实践

## 输出格式
每行一个变量名，格式: 变量名 - 说明
示例:
userName - 用户名称
accountName - 账户名称
userAccount - 用户账户

## 命名原则
1. 语义准确: 准确表达概念含义
2. 简洁明了: 避免冗长或复杂的词汇
3. 约定俗成: 使用业界通用术语
4. 上下文适配: 考虑在代码中的使用场景

现在开始生成变量名:`, chineseText, styleDesc, examples)
}

// parseResponse 解析AI响应
func (c *Converter) parseResponse(content string) []string {
	lines := strings.Split(content, "\n")
	var results []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 移除序号、项目符号等前缀
		line = strings.TrimPrefix(line, "- ")
		line = strings.TrimPrefix(line, "• ")

		// 查找数字前缀并移除
		for i := 1; i <= 10; i++ {
			prefix := fmt.Sprintf("%d. ", i)
			if strings.HasPrefix(line, prefix) {
				line = strings.TrimPrefix(line, prefix)
				break
			}
		}

		// 提取变量名（在 - 之前的部分）
		if parts := strings.Split(line, " - "); len(parts) >= 2 {
			varName := strings.TrimSpace(parts[0])
			if varName != "" && c.isValidVariableName(varName) {
				results = append(results, varName)
			}
		} else if c.isValidVariableName(line) {
			// 如果整行就是一个变量名
			results = append(results, line)
		}
	}

	// 如果解析失败，返回原始内容的清理版本
	if len(results) == 0 {
		cleanContent := strings.ReplaceAll(content, "\n", " ")
		cleanContent = strings.TrimSpace(cleanContent)
		if cleanContent != "" {
			results = append(results, cleanContent)
		}
	}

	return results
}

// isValidVariableName 检查是否为有效的变量名
func (c *Converter) isValidVariableName(name string) bool {
	if name == "" {
		return false
	}

	// 基本检查：只包含字母、数字、下划线、短横线
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}

	// 不能以数字开头
	if name[0] >= '0' && name[0] <= '9' {
		return false
	}

	// 长度检查
	if len(name) > 50 {
		return false
	}

	return true
}

// ConvertToVariableNameStream 流式将中文转换为英文变量名
func (c *Converter) ConvertToVariableNameStream(chineseText string, style string, onContent func(string), onComplete func([]string)) error {
	prompt := c.buildPrompt(chineseText, style)

	messages := []client.Message{
		{
			Role:    "system",
			Content: "你是一位资深的软件工程师和编程规范专家，精通多种编程语言的命名约定。你的任务是将中文概念准确转换为地道的英文变量名，确保：1) 语义准确表达原始概念；2) 遵循目标命名风格；3) 符合国际编程最佳实践；4) 使用简洁明了的英语词汇；5）请不要回复无关的信息，仅回复变量名，不要回复任何其他信息。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	responseChan, errorChan := c.client.ChatStream(messages)

	var fullContent strings.Builder

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				// 通道关闭，处理完整内容
				results := c.parseResponse(fullContent.String())
				onComplete(results)
				return nil
			}

			if len(response.Choices) > 0 {
				content := response.Choices[0].Delta.Content
				if content != "" {
					fullContent.WriteString(content)
					onContent(content)
				}
			}

		case err := <-errorChan:
			if err != nil {
				return fmt.Errorf("AI流式请求失败: %w", err)
			}
		}
	}
}
