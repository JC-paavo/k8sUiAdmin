#!/bin/bash
set -e

echo "============================================"
echo "  K8s UI Admin - 构建脚本"
echo "============================================"

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
echo "项目目录: $PROJECT_ROOT"

echo ""
echo "========== 1. 构建前端 =========="
cd "$PROJECT_ROOT/frontend"
npm install
npm run build
echo "前端构建完成: $PROJECT_ROOT/frontend/dist"

echo ""
echo "========== 2. 复制前端到后端 =========="
rm -rf "$PROJECT_ROOT/backend/cmd/dist"
cp -r "$PROJECT_ROOT/frontend/dist" "$PROJECT_ROOT/backend/cmd/dist"
echo "已复制到: $PROJECT_ROOT/backend/cmd/dist"

echo ""
echo "========== 3. 构建后端 =========="
cd "$PROJECT_ROOT/backend"
go mod download
CGO_ENABLED=1 go build -o "$PROJECT_ROOT/k8s-ui-admin" ./cmd/...
echo "后端构建完成: $PROJECT_ROOT/k8s-ui-admin"

echo ""
echo "========== 4. 清理临时文件 =========="
rm -rf "$PROJECT_ROOT/backend/cmd/dist"

echo ""
echo "============================================"
echo "  构建完成！"
echo "  可执行文件: $PROJECT_ROOT/k8s-ui-admin"
echo "  运行: ./k8s-ui-admin"
echo "============================================"
