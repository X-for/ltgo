#!/bin/bash

# 项目名称
APP_NAME="ltgo"
# 源码入口
MAIN_PATH="./cmd/ltgo"
# 输出目录
DIST_DIR="dist"

# 清理旧的构建
rm -rf $DIST_DIR
mkdir -p $DIST_DIR

echo "🚀 Starting build process for $APP_NAME..."

# 1. Windows (amd64)
echo "📦 Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o "$DIST_DIR/${APP_NAME}-windows-amd64.exe" $MAIN_PATH

# 2. Linux (amd64)
echo "🐧 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$DIST_DIR/${APP_NAME}-linux-amd64" $MAIN_PATH

# 3. macOS (Intel)
echo "🍎 Building for macOS (Intel)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o "$DIST_DIR/${APP_NAME}-darwin-amd64" $MAIN_PATH

# 4. macOS (Apple Silicon M1/M2)
echo "🍎 Building for macOS (Apple Silicon)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o "$DIST_DIR/${APP_NAME}-darwin-arm64" $MAIN_PATH

echo "✅ Build complete! All binaries are in '$DIST_DIR/' directory."
ls -lh $DIST_DIR
