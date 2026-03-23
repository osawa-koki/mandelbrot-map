import { Viewport } from "./viewport";

const TILE_SIZE = 256;
const API_ENDPOINT = "/api/tile";

interface PendingTile {
  screenX: number;
  screenY: number;
  controller: AbortController;
}

/** タイルの取得・キャッシュ・描画を管理する */
export class TileManager {
  private cache = new Map<string, ImageBitmap>();
  private pending = new Map<string, PendingTile>();
  private colorScheme = "classic";
  private maxIter = 256;

  setColorScheme(scheme: string) {
    this.colorScheme = scheme;
    this.clearCache();
  }

  setMaxIter(iter: number) {
    this.maxIter = iter;
    this.clearCache();
  }

  clearCache() {
    for (const bitmap of this.cache.values()) bitmap.close();
    this.cache.clear();
    this.abortAll();
  }

  /** 現在のビューポートに必要なタイルを Canvas に描画する */
  render(
    ctx: CanvasRenderingContext2D,
    viewport: Viewport,
    width: number,
    height: number,
    onTileLoaded: () => void
  ) {
    // 不要なリクエストをキャンセル
    const needed = new Set<string>();
    const tiles = this.visibleTiles(viewport, width, height);

    for (const { key } of tiles) needed.add(key);

    for (const [key, pending] of this.pending) {
      if (!needed.has(key)) {
        pending.controller.abort();
        this.pending.delete(key);
      }
    }

    // キャッシュ済みタイルを描画し、未取得タイルをリクエスト
    for (const { key, screenX, screenY, params } of tiles) {
      const cached = this.cache.get(key);
      if (cached) {
        ctx.drawImage(cached, screenX, screenY, TILE_SIZE, TILE_SIZE);
      } else if (!this.pending.has(key)) {
        this.fetchTile(key, screenX, screenY, params, onTileLoaded);
      }
    }
  }

  /** 表示領域を覆うタイル一覧を計算する */
  private visibleTiles(viewport: Viewport, width: number, height: number) {
    const results: {
      key: string;
      screenX: number;
      screenY: number;
      params: string;
    }[] = [];

    // 左上タイルの起点を求める
    const tileScale = viewport.scale * TILE_SIZE;
    const originX =
      Math.floor(
        (viewport.centerX - (width / 2) * viewport.scale) / tileScale
      ) * tileScale;
    const originY =
      Math.floor(
        (viewport.centerY - (height / 2) * viewport.scale) / tileScale
      ) * tileScale;

    for (
      let cx = originX;
      cx < viewport.centerX + (width / 2) * viewport.scale;
      cx += tileScale
    ) {
      for (
        let cy = originY;
        cy < viewport.centerY + (height / 2) * viewport.scale;
        cy += tileScale
      ) {
        const xMin = cx;
        const xMax = cx + tileScale;
        const yMin = cy;
        const yMax = cy + tileScale;

        const key = `${xMin.toFixed(10)},${xMax.toFixed(10)},${yMin.toFixed(10)},${yMax.toFixed(10)},${this.maxIter},${this.colorScheme}`;

        const screenX = Math.round(
          (cx - viewport.centerX) / viewport.scale + width / 2
        );
        const screenY = Math.round(
          (cy - viewport.centerY) / viewport.scale + height / 2
        );

        const params = new URLSearchParams({
          x_min: xMin.toString(),
          x_max: xMax.toString(),
          y_min: yMin.toString(),
          y_max: yMax.toString(),
          width: TILE_SIZE.toString(),
          height: TILE_SIZE.toString(),
          max_iter: this.maxIter.toString(),
          color_scheme: this.colorScheme,
        }).toString();

        results.push({ key, screenX, screenY, params });
      }
    }

    return results;
  }

  private async fetchTile(
    key: string,
    screenX: number,
    screenY: number,
    params: string,
    onLoaded: () => void
  ) {
    const controller = new AbortController();
    this.pending.set(key, { screenX, screenY, controller });

    try {
      const res = await fetch(`${API_ENDPOINT}?${params}`, {
        signal: controller.signal,
      });
      if (!res.ok) return;

      const blob = await res.blob();
      const bitmap = await createImageBitmap(blob);

      this.cache.set(key, bitmap);
      this.pending.delete(key);
      onLoaded();
    } catch (e) {
      if ((e as Error).name !== "AbortError") {
        console.error("タイル取得失敗:", e);
      }
      this.pending.delete(key);
    }
  }

  private abortAll() {
    for (const { controller } of this.pending.values()) controller.abort();
    this.pending.clear();
  }
}
