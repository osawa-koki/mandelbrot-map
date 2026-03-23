# ---- フロントエンドビルド ----
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- バックエンドビルド ----
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# フロントエンドのビルド成果物を組み込む
COPY --from=frontend-builder /app/frontend/dist ./dist
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# ---- 本番イメージ ----
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/backend/server .
COPY --from=backend-builder /app/backend/dist ./dist
EXPOSE 8080
CMD ["./server"]
