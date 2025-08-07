# 发布指南

本文档描述如何使用 GitHub Actions 自动构建和发布 ZHV 项目。

## 🚀 自动发布流程

### 概述

项目配置了完整的 CI/CD 流程，当推送标签时会自动：

1. **跨平台构建** - 构建多个操作系统和架构的二进制文件
2. **创建 Release** - 自动创建 GitHub Release
3. **上传文件** - 将构建的二进制文件上传到 Release

### 支持的平台

| 操作系统 | 架构 | 文件格式 |
|---------|------|----------|
| Linux | x86_64 | `.tar.gz` |
| Linux | ARM64 | `.tar.gz` |
| macOS | x86_64 | `.tar.gz` |
| macOS | ARM64 (M1/M2) | `.tar.gz` |
| Windows | x86_64 | `.zip` |
| FreeBSD | x86_64 | `.tar.gz` |

## 📦 如何发布新版本

### 1. 准备发布

确保代码已准备好发布：

```bash
# 更新依赖
go mod tidy

# 运行测试
make test

# 本地构建测试
make build-all

# 代码格式化
make fmt

# 代码检查（如果安装了 golangci-lint）
make lint
```

### 2. 创建和推送标签

使用语义化版本标签（如 v1.0.0）：

```bash
# 创建带注释的标签
git tag -a v1.0.0 -m "Release v1.0.0

- 添加功能 A
- 修复 Bug B
- 改进性能"

# 推送标签到远程仓库
git push origin v1.0.0
```

### 3. 自动构建和发布

推送标签后，GitHub Actions 会自动：

1. 触发构建工作流
2. 为所有支持的平台构建二进制文件
3. 创建 GitHub Release
4. 上传构建产物
5. 生成变更日志

### 4. 验证发布

发布完成后，检查：

- [ ] GitHub Release 页面显示新版本
- [ ] 所有平台的二进制文件都已上传
- [ ] 变更日志生成正确
- [ ] 下载链接正常工作

## 🛠 本地开发和测试

### 环境设置

```bash
# 克隆项目
git clone https://github.com/iiileo/zhv.git
cd zhv

# 安装依赖
make deps
```

### 本地构建

```bash
# 构建当前平台
make build

# 构建所有平台（测试跨平台构建）
make build-all

# 清理构建文件
make clean
```

### 测试版本信息

```bash
# 构建并测试版本命令
make build
./zhv version
```

## 🔧 工作流配置

### 工作流文件

主要配置文件：`.github/workflows/release.yml`

### 触发条件

- **标签推送**：推送以 `v` 开头的标签（如 `v1.0.0`）
- **手动触发**：在 GitHub Actions 页面手动运行

### 构建特性

- **跨平台编译**：支持 Linux、macOS、Windows、FreeBSD
- **版本嵌入**：构建时注入版本、构建时间、Git 提交信息
- **自动压缩**：Linux/macOS 使用 tar.gz，Windows 使用 zip
- **缓存优化**：Go 模块缓存加速构建

### 发布特性

- **自动变更日志**：基于 Git 提交生成
- **预发布检测**：包含 rc/beta/alpha 的标签标记为预发布
- **下载说明**：自动生成平台特定的下载和安装说明

## 🏷 版本管理

### 版本号规范

使用 [语义化版本](https://semver.org/lang/zh-CN/)：

- **主版本号**：不兼容的 API 修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 标签格式

```bash
# 正式版本
v1.0.0
v1.2.3

# 预发布版本
v1.0.0-rc.1    # 发布候选
v1.0.0-beta.1  # 测试版
v1.0.0-alpha.1 # 内测版
```

### 分支策略

- **main**：主分支，包含稳定代码
- **develop**：开发分支，用于集成新功能
- **feature/**：功能分支
- **hotfix/**：热修复分支

## 🔍 故障排除

### 常见问题

1. **构建失败**
   ```bash
   # 检查 Go 版本
   go version
   
   # 检查依赖
   go mod verify
   ```

2. **标签推送失败**
   ```bash
   # 检查远程仓库
   git remote -v
   
   # 重新推送标签
   git push origin v1.0.0 --force
   ```

3. **工作流权限问题**
   - 确保仓库设置中启用了 Actions
   - 检查 `GITHUB_TOKEN` 权限

### 调试工作流

1. 查看 GitHub Actions 日志
2. 检查工作流文件语法
3. 验证环境变量和密钥

## 📚 相关资源

- [GitHub Actions 文档](https://docs.github.com/cn/actions)
- [Go 交叉编译文档](https://golang.org/doc/install/source#environment)
- [语义化版本规范](https://semver.org/lang/zh-CN/)
- [Git 标签文档](https://git-scm.com/book/zh/v2/Git-%E5%9F%BA%E7%A1%80-%E6%89%93%E6%A0%87%E7%AD%BE)

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request
5. 等待代码审查

发布新版本前，确保：
- [ ] 所有测试通过
- [ ] 代码审查完成
- [ ] 更新文档
- [ ] 准备发布说明
