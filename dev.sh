#!/bin/bash
# バックエンドとフロントエンドを並列で起動する開発用スクリプト

set -e

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"

cleanup() {
  echo ""
  echo "停止中..."
  kill "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null
  wait "$BACKEND_PID" "$FRONTEND_PID" 2>/dev/null
  echo "停止しました"
}
trap cleanup EXIT INT TERM

echo "バックエンドを起動中 (port 8080)..."
cd "$ROOT_DIR/backend"
go run main.go &
BACKEND_PID=$!

echo "フロントエンドを起動中 (port 3000)..."
cd "$ROOT_DIR/frontend"
npm run dev &
FRONTEND_PID=$!

echo ""
echo "起動完了"
echo "  Backend:  http://localhost:8080"
echo "  Frontend: http://localhost:3000"
echo ""
echo "停止するには Ctrl+C を押してください"

wait "$BACKEND_PID" "$FRONTEND_PID"
