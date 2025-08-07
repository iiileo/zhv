package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	APIURL string `json:"api_url"`
	Model  string `json:"model"`
	APIKey string `json:"api_key"`
}

// LoadConfig 加载配置，优先级：环境变量 > 配置文件 > 默认值
func LoadConfig() (*Config, error) {
	config := &Config{
		APIURL: "https://api.openai.com/v1",
		Model:  "gpt-3.5-turbo",
		APIKey: "",
	}

	// 尝试从配置文件读取
	if err := loadFromFile(config); err != nil {
		// 配置文件不存在或读取失败时不报错，使用默认值
	}

	// 从环境变量覆盖配置
	loadFromEnv(config)

	return config, nil
}

// loadFromFile 从配置文件加载
func loadFromFile(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".zhv", "setting.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, config)
}

// loadFromEnv 从环境变量加载
func loadFromEnv(config *Config) {
	if apiURL := os.Getenv("ZHV_API_URL"); apiURL != "" {
		config.APIURL = apiURL
	}
	if model := os.Getenv("ZHV_MODEL"); model != "" {
		config.Model = model
	}
	if apiKey := os.Getenv("ZHV_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}
}

// SaveConfig 保存配置到文件
func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".zhv")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "setting.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// IsValid 检查配置是否有效
func (c *Config) IsValid() bool {
	return c.APIKey != "" && c.APIURL != "" && c.Model != ""
}
