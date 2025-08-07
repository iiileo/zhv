# ZHV - 中文转英文变量名推荐工具

ZHV (中文变量) 是一个帮助开发者将中文词汇转换为符合编程规范的英文变量名的命令行工具。

## 功能特性

- 🎯 **智能转换**: 使用AI模型将中文词汇转换为合适的英文变量名
- 🎨 **多种命名风格**: 支持驼峰、帕斯卡、蛇形、短横线等命名风格
- ⚙️ **灵活配置**: 支持环境变量和配置文件两种配置方式
- 🔌 **OpenAI兼容**: 支持任何OpenAI兼容的API接口
- 📝 **多个推荐**: 为每个输入提供多个变量名选项

## 安装

### 方式 1: 下载预编译二进制文件（推荐）

从 [Releases 页面](https://github.com/iiileo/zhv/releases) 下载适合您系统的二进制文件：

| 操作系统 | 架构 | 下载文件 |
|---------|------|----------|
| Linux | x86_64 | `zhv-linux-amd64.tar.gz` |
| Linux | ARM64 | `zhv-linux-arm64.tar.gz` |
| macOS | x86_64 | `zhv-darwin-amd64.tar.gz` |
| macOS | ARM64 (M1/M2) | `zhv-darwin-arm64.tar.gz` |
| Windows | x86_64 | `zhv-windows-amd64.zip` |
| FreeBSD | x86_64 | `zhv-freebsd-amd64.tar.gz` |

**Linux/macOS 安装步骤：**
```bash
# 下载并解压（以 Linux x86_64 为例）
wget https://github.com/iiileo/zhv/releases/latest/download/zhv-linux-amd64.tar.gz
tar -xzf zhv-linux-amd64.tar.gz

# 移动到 PATH 目录
sudo mv zhv-linux-amd64 /usr/local/bin/zhv

# 验证安装
zhv --help
```

**Windows 安装步骤：**
1. 下载 `zhv-windows-amd64.zip`
2. 解压到任意目录
3. 将解压目录添加到系统 PATH 环境变量
4. 在命令提示符中运行 `zhv --help` 验证安装

### 方式 2: Go 安装

```bash
go install github.com/iiileo/zhv@latest
```

### 方式 3: 从源码编译

```bash
git clone https://github.com/iiileo/zhv.git
cd zhv
make build
# 或者
go build
```

## 配置

### 方式1: 环境变量

```bash
export ZHV_API_URL="your-api-url"       # API地址
export ZHV_MODEL="your-model"           # 模型名称
export ZHV_KEY="your-api-key"           # API密钥
```

### 方式2: 配置文件

使用命令行设置：

```bash
zhv config set api_url "your-api-url"
zhv config set model "your-model"
zhv config set api_key "your-api-key"
```

配置文件位置：`~/.zhv/setting.json`

### 方式3: 查看当前配置

```bash
zhv config show
```

## 使用方法

### 基本用法

```bash
# 转换单个词汇
zhv 用户名称

# 转换短语
zhv "数据库连接池"

# 指定命名风格
zhv -s snake 文件上传状态
zhv -s pascal 用户管理系统
zhv -s kebab 前端组件名称

# 详细输出
zhv -v 购物车商品
```

### 命名风格

| 风格 | 参数 | 示例 | 说明 |
|------|------|------|------|
| 驼峰命名法 | `camel` | `userName` | 默认风格，首字母小写 |
| 帕斯卡命名法 | `pascal` | `UserName` | 首字母大写 |
| 蛇形命名法 | `snake` | `user_name` | 下划线分隔 |
| 短横线命名法 | `kebab` | `user-name` | 短横线分隔 |

### 示例输出

```bash
$ zhv -s camel 用户个人信息

中文: 用户个人信息
风格: 驼峰命名法 (camelCase)
推荐的变量名:
  1. userPersonalInfo
  2. userProfile
  3. personalInfo
  4. userDetails
  5. profileInfo
```

## 支持的AI服务

本工具支持任何兼容OpenAI API格式的服务，包括但不限于：

- OpenAI GPT-3.5/GPT-4
- Azure OpenAI Service
- 国内AI服务（如智谱AI、百度文心等）
- 自部署的本地模型服务

## 命令参考

### 主命令

```bash
zhv [选项] <中文文本>
```

**选项：**
- `-s, --style <风格>`: 指定命名风格 (camel|pascal|snake|kebab)
- `-v, --verbose`: 显示详细信息
- `-h, --help`: 显示帮助信息

### 配置命令

```bash
# 设置配置
zhv config set <key> <value>

# 显示配置
zhv config show

# 配置项说明
# api_url: API服务地址
# model: 使用的模型名称
# api_key: API密钥
```

## 项目结构

```
zhv/
├── main.go              # 主程序入口和CLI界面
├── config/              # 配置管理
│   └── config.go        # 配置文件和环境变量处理
├── client/              # API客户端
│   └── openai.go        # OpenAI兼容接口实现
├── converter/           # 转换器
│   └── converter.go     # 中文转英文变量名核心逻辑
├── go.mod               # Go模块定义
└── README.md            # 项目说明文档
```

## 开发指南

### 环境要求

- Go 1.24.2 或更高版本

### 本地开发

```bash
# 克隆项目
git clone https://github.com/iiileo/zhv.git
cd zhv

# 安装依赖
make deps
# 或者
go mod tidy

# 构建项目
make build
# 或者
go build

# 运行测试
make test
# 或者
go test ./...

# 构建所有平台（测试跨平台编译）
make build-all

# 查看版本信息
make version
./zhv version

# 清理构建文件
make clean
```

### 可用的 Make 命令

| 命令 | 说明 |
|------|------|
| `make build` | 构建当前平台的二进制文件 |
| `make build-all` | 构建所有支持平台的二进制文件 |
| `make test` | 运行测试 |
| `make deps` | 下载和整理依赖 |
| `make clean` | 清理构建文件 |
| `make install` | 安装到本地 (/usr/local/bin) |
| `make uninstall` | 从本地卸载 |
| `make version` | 显示版本信息 |
| `make fmt` | 格式化代码 |
| `make lint` | 运行代码检查（需要 golangci-lint） |
| `make help` | 显示帮助信息 |

### 贡献代码

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

### 发布流程

项目使用 GitHub Actions 进行自动化构建和发布：

- **自动构建**：推送标签时自动构建所有平台的二进制文件
- **自动发布**：创建 GitHub Release 并上传构建产物
- **版本管理**：支持语义化版本和预发布版本

详细的发布指南请参考 [RELEASE.md](RELEASE.md)。

#### 快速发布

```bash
# 创建并推送标签
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions 会自动构建和发布
```

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 常见问题

### Q: 如何配置使用不同的AI服务？

A: 只需要设置对应服务的API地址、模型名称和API密钥即可。例如使用Azure OpenAI：

```bash
zhv config set api_url "https://your-resource.openai.azure.com/openai/deployments/your-deployment"
zhv config set model "gpt-35-turbo"
zhv config set api_key "your-azure-api-key"
```

### Q: 支持哪些命名风格？

A: 目前支持四种主流的编程命名风格：
- camel: 驼峰命名法 (默认)
- pascal: 帕斯卡命名法
- snake: 蛇形命名法  
- kebab: 短横线命名法

### Q: 如何获得更准确的变量名推荐？

A: 建议：
1. 提供准确且具体的中文描述
2. 使用专业术语和标准表达
3. 避免过于复杂的长句
4. 根据编程语言选择合适的命名风格

### Q: 配置文件存储在哪里？

A: 配置文件存储在用户主目录下的 `.zhv/setting.json` 文件中。

## 更新日志

### v1.0.0
- 首次发布
- 支持中文转英文变量名
- 支持多种命名风格
- 支持环境变量和配置文件
- 兼容OpenAI API格式
