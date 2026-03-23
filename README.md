# Mandelbrot Map

🗺️🗺️🗺️ マンデルブロ集合の地図アプリ！  
通常の地図アプリと同様の操作感でマンデルブロ集合を探索できる。  

## 機能

- スクロール移動（ドラッグ操作）
- ピンチイン・ピンチアウトによるズーム
- マウスホイールによるズーム

## アーキテクチャ

| ディレクトリ | 役割 | 技術スタック |
|-------------|------|------------|
| `backend/` | APIサーバー（マンデルブロ集合の計算） | Go |
| `frontend/` | フロントエンド（描画・操作UI） | TypeScript / React |

## 開発環境

### 必要な環境

- Go 1.22+
- Node.js 20+

### セットアップ

```bash
# バックエンド
cd backend
go mod tidy
go run .

# フロントエンド
cd frontend
npm install
npm run dev
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000

## 開発フロー

`CLAUDE.md` に記載されている開発ワークフローに従うこと。

1. Issue作成 (`/create-issue`)
2. ブランチ作成 (`issues/番号_説明`)
3. 実装
4. PR作成 (`/create-pr`)
5. PRマージ (`/merge-pr`)