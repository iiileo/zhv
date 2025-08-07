package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/iiileo/zhv/config"
)

// OpenAIClient OpenAI兼容客户端
type OpenAIClient struct {
	config     *config.Config
	httpClient *http.Client
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest 聊天请求结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Choice 选择结构
type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

// ChatResponse 聊天响应结构
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// StreamDelta 流式增量内容
type StreamDelta struct {
	Content string `json:"content"`
}

// StreamChoice 流式选择结构
type StreamChoice struct {
	Index        int         `json:"index"`
	Delta        StreamDelta `json:"delta"`
	FinishReason *string     `json:"finish_reason"`
}

// StreamResponse 流式响应结构
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// NewOpenAIClient 创建新的OpenAI客户端
func NewOpenAIClient(cfg *config.Config) *OpenAIClient {
	return &OpenAIClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Chat 发送聊天请求
func (c *OpenAIClient) Chat(messages []Message) (*ChatResponse, error) {
	request := ChatRequest{
		Model:       c.config.Model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("编码请求失败: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.APIURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &chatResponse, nil
}

// ChatStream 发送流式聊天请求
func (c *OpenAIClient) ChatStream(messages []Message) (<-chan StreamResponse, <-chan error) {
	responseChan := make(chan StreamResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		request := ChatRequest{
			Model:       c.config.Model,
			Messages:    messages,
			Temperature: 0.7,
			MaxTokens:   1000,
			Stream:      true,
		}

		jsonData, err := json.Marshal(request)
		if err != nil {
			errorChan <- fmt.Errorf("编码请求失败: %w", err)
			return
		}

		req, err := http.NewRequest("POST", c.config.APIURL+"/chat/completions", bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("创建请求失败: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Cache-Control", "no-cache")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("发送请求失败: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// 跳过空行和非数据行
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			// 移除 "data: " 前缀
			data := strings.TrimPrefix(line, "data: ")

			// 检查是否为结束标志
			if data == "[DONE]" {
				break
			}

			// 解析JSON响应
			var streamResp StreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				// 忽略解析错误，继续处理下一行
				continue
			}

			// 发送响应到通道
			select {
			case responseChan <- streamResp:
			case <-time.After(5 * time.Second):
				errorChan <- fmt.Errorf("发送响应超时")
				return
			}
		}

		if err := scanner.Err(); err != nil {
			errorChan <- fmt.Errorf("读取响应流失败: %w", err)
		}
	}()

	return responseChan, errorChan
}
