#!/bin/bash

# 测试构建脚本
# 用于验证GitHub Actions工作流的本地测试

set -e

echo "🔧 开始测试构建流程..."

# 检查Go环境
echo "📋 检查环境..."
echo "Go版本: $(go version)"
echo "当前目录: $(pwd)"

# 清理之前的构建
echo "🧹 清理之前的构建文件..."
make clean

# 测试基本构建
echo "🔨 测试基本构建..."
make build

# 测试版本命令
echo "📝 测试版本命令..."
./zhv version
echo ""

# 测试帮助命令
echo "❓ 测试帮助命令..."
./zhv --help
echo ""

# 测试配置命令
echo "⚙️ 测试配置命令..."
./zhv config show || echo "配置显示正常（可能未配置）"
echo ""

# 测试跨平台构建（如果make build-all可用）
echo "🌍 测试跨平台构建..."
if make build-all; then
    echo "✅ 跨平台构建成功"
    echo "📦 生成的文件："
    ls -la zhv-* 2>/dev/null || echo "没有找到跨平台构建文件"
else
    echo "⚠️ 跨平台构建失败或不可用"
fi

echo ""
echo "🎉 构建测试完成！"
echo ""
echo "📁 当前目录文件："
ls -la zhv* 2>/dev/null || echo "没有找到构建文件"
