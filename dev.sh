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

# 指定ポートで動作中のプロセスを kill する
kill_port() {
  local port=$1
  local pid
  pid=$(lsof -ti :"$port" 2>/dev/null || true)
  if [ -n "$pid" ]; then
    echo "port $port を使用中のプロセス (PID: $pid) を停止します..."
    kill "$pid" 2>/dev/null || true
    sleep 1
  fi
}

kill_port 8080
kill_port 3000

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
