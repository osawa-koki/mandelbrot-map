# Mandelbrot Map

🗺️🗺️🗺️ マンデルブロ集合の地図アプリ！  

![成果物](./fruit.gif)  
（無料プランなので動作は遅い...）  

通常の地図アプリと同様の操作感でマンデルブロ集合を探索できる。  

## 機能

- マウスドラッグ・タッチドラッグによるパン移動
- マウスホイール・ピンチによるズーム
- カラースキーム切替（Classic / Fire / Ocean）
- 最大反復回数の調整
- ビューポート状態の URL パラメータ保持（ブックマーク・共有対応）
- 初期値リセットボタン

## アーキテクチャ

| ディレクトリ | 役割 | 技術スタック |
|-------------|------|------------|
| `backend/` | タイル画像生成 API | Go / Gin |
| `frontend/` | Canvas タイルマップ | TypeScript / Vite |

### データフロー

フロントエンドがビューポートをタイルに分割し、各タイルを `GET /api/tile` で並列リクエストする。バックエンドが goroutine で並列計算して PNG 画像を返す。

## 開発環境

**前提条件**: Go 1.25+、Node.js 22+

```bash
./dev.sh
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000

`dev.sh` は起動時に依存パッケージの自動インストールと使用中ポートの解放も行う。

## API リファレンス

### `GET /api/tile`

マンデルブロ集合タイル画像（PNG）を返す。すべてのパラメータが必須。

| パラメータ | 型 | 説明 |
|-----------|-----|------|
| `x_min` | float64 | 複素数平面 X 最小値 |
| `x_max` | float64 | 複素数平面 X 最大値 |
| `y_min` | float64 | 複素数平面 Y 最小値 |
| `y_max` | float64 | 複素数平面 Y 最大値 |
| `width` | int | 画像幅（px） |
| `height` | int | 画像高さ（px） |
| `max_iter` | int | 最大反復回数 |
| `color_scheme` | string | `classic` / `fire` / `ocean` |

**例**

```bash
curl "http://localhost:8080/api/tile?x_min=-2.5&x_max=1.0&y_min=-1.25&y_max=1.25&width=256&height=256&max_iter=256&color_scheme=classic" -o tile.png
```

## デプロイ（Render）

`render.yaml` による設定済み。Render ダッシュボードでリポジトリを接続するだけでデプロイできる。

### 環境変数

| 変数 | 説明 |
|------|------|
| `PORT` | サーバーポート（デフォルト: `8080`） |
| `GIN_MODE` | Gin 動作モード（本番: `release`） |

### ローカルで Docker ビルドを確認する

```bash
docker build -t mandelbrot-map .
docker run -p 8080:8080 mandelbrot-map
```

## 開発フロー

`CLAUDE.md` に記載されている開発ワークフローに従う。

1. Issue 作成（`/create-issue`）
2. ブランチ作成（`issues/番号_説明`）
3. 実装
4. PR 作成（`/create-pr`）
5. PR マージ（`/merge-pr`）
