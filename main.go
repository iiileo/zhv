package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/iiileo/zhv/config"
	"github.com/iiileo/zhv/converter"
	"github.com/spf13/cobra"
)

var (
	style   string
	verbose bool

	// 版本信息，在构建时通过 ldflags 注入
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

// rootCmd 根命令
var rootCmd = &cobra.Command{
	Use:   "zhv [中文文本]",
	Short: "中文转英文变量名推荐工具",
	Long: `ZHV (中文变量) 是一个帮助开发者将中文词汇转换为符合编程规范的英文变量名的工具。

支持多种命名风格：
  - camel: 驼峰命名法 (userName)
  - pascal: 帕斯卡命名法 (UserName)  
  - snake: 蛇形命名法 (user_name)
  - kebab: 短横线命名法 (user-name)

配置方式：
  1. 环境变量：ZHV_API_URL, ZHV_MODEL, ZHV_KEY
  2. 配置文件：~/.zhv/setting.json

示例：
  zhv 用户名称
  zhv -s snake 数据库连接
  zhv -s pascal "文件上传状态"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 合并所有参数为一个中文文本
		chineseText := strings.Join(args, " ")

		if err := convertAndDisplay(chineseText); err != nil {
			fmt.Fprintf(os.Stderr, "错误: %v\n", err)
			os.Exit(1)
		}
	},
}

// configCmd 配置命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "配置管理",
	Long:  "管理ZHV的配置信息",
}

// setConfigCmd 设置配置
var setConfigCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "设置配置项",
	Long: `设置配置项，支持的配置项：
  - api_url: API地址
  - model: 模型名称
  - api_key: API密钥`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		if err := setConfig(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "设置配置失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("配置 %s 已设置\n", key)
	},
}

// showConfigCmd 显示配置
var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Run: func(cmd *cobra.Command, args []string) {
		if err := showConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "显示配置失败: %v\n", err)
			os.Exit(1)
		}
	},
}

// versionCmd 版本命令
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示ZHV的版本信息，包括版本号、构建时间和Git提交哈希",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ZHV (中文变量名推荐工具)\n")
		fmt.Printf("版本: %s\n", Version)
		fmt.Printf("构建时间: %s\n", BuildTime)
		fmt.Printf("Git提交: %s\n", GitCommit)
		fmt.Printf("Go版本: %s\n", runtime.Version())
		fmt.Printf("平台: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	// 根命令的标志
	rootCmd.PersistentFlags().StringVarP(&style, "style", "s", "camel", "命名风格 (camel|pascal|snake|kebab)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出")

	// 配置命令
	configCmd.AddCommand(setConfigCmd)
	configCmd.AddCommand(showConfigCmd)
	rootCmd.AddCommand(configCmd)

	// 版本命令
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行命令失败: %v\n", err)
		os.Exit(1)
	}
}

// convertAndDisplay 转换并显示结果
func convertAndDisplay(chineseText string) error {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	// 检查配置是否有效
	if !cfg.IsValid() {
		return fmt.Errorf(`配置不完整，请设置以下配置：

方式1 - 使用环境变量：
  export ZHV_API_URL="your-api-url"
  export ZHV_MODEL="your-model"  
  export ZHV_KEY="your-api-key"

方式2 - 使用配置文件：
  zhv config set api_url "your-api-url"
  zhv config set model "your-model"
  zhv config set api_key "your-api-key"

方式3 - 查看当前配置：
  zhv config show`)
	}

	if verbose {
		fmt.Printf("使用配置: API=%s, Model=%s\n", cfg.APIURL, cfg.Model)
		fmt.Printf("转换文本: %s\n", chineseText)
		fmt.Printf("命名风格: %s\n\n", style)
	}

	// 创建转换器
	conv := converter.NewConverter(cfg)

	// 显示基本信息
	fmt.Printf("中文: %s\n", chineseText)
	fmt.Printf("风格: %s\n", getStyleDescription(style))
	fmt.Println("正在生成变量名推荐...")
	fmt.Println()

	// 使用流式输出
	var mu sync.Mutex
	var streamContent strings.Builder
	var hasContent bool

	err = conv.ConvertToVariableNameStream(chineseText, style,
		func(content string) {
			// 实时显示内容
			mu.Lock()
			defer mu.Unlock()

			if !hasContent {
				fmt.Print("AI回复: \n")
				hasContent = true
			}

			fmt.Print(content)
			streamContent.WriteString(content)
		},
		func(results []string) {
			// 完成后显示整理后的结果
			mu.Lock()
			defer mu.Unlock()

			fmt.Println()

			if len(results) == 0 {
				fmt.Println("未找到合适的变量名推荐")
				return
			}

			fmt.Println("推荐的变量名:")
			for i, result := range results {
				fmt.Printf("  %d. %s\n", i+1, result)
			}
		},
	)

	if err != nil {
		return fmt.Errorf("转换失败: %w", err)
	}

	return nil
}

// setConfig 设置配置
func setConfig(key, value string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	switch key {
	case "api_url":
		cfg.APIURL = value
	case "model":
		cfg.Model = value
	case "api_key":
		cfg.APIKey = value
	default:
		return fmt.Errorf("未知的配置项: %s", key)
	}

	return config.SaveConfig(cfg)
}

// showConfig 显示配置
func showConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	fmt.Println("当前配置:")
	fmt.Printf("  API地址: %s\n", cfg.APIURL)
	fmt.Printf("  模型: %s\n", cfg.Model)

	if cfg.APIKey != "" {
		// 隐藏API密钥的大部分内容
		maskedKey := maskAPIKey(cfg.APIKey)
		fmt.Printf("  API密钥: %s\n", maskedKey)
	} else {
		fmt.Printf("  API密钥: (未设置)\n")
	}

	fmt.Printf("\n配置状态: ")
	if cfg.IsValid() {
		fmt.Println("✓ 配置完整")
	} else {
		fmt.Println("✗ 配置不完整")
	}

	return nil
}

// maskAPIKey 隐藏API密钥
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

// getStyleDescription 获取风格描述
func getStyleDescription(style string) string {
	switch style {
	case "camel":
		return "驼峰命名法 (camelCase)"
	case "pascal":
		return "帕斯卡命名法 (PascalCase)"
	case "snake":
		return "蛇形命名法 (snake_case)"
	case "kebab":
		return "短横线命名法 (kebab-case)"
	default:
		return "驼峰命名法 (camelCase)"
	}
}
