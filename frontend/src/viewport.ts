/** Canvas 上の複素数平面ビューポートを管理する */
export class Viewport {
  /** 複素数平面の中心座標 */
  centerX: number;
  centerY: number;
  /** 1ピクセルあたりの複素数平面の幅 */
  scale: number;

  constructor(centerX = -0.5, centerY = 0, scale = 0.004) {
    this.centerX = centerX;
    this.centerY = centerY;
    this.scale = scale;
  }

  /** Canvas 座標を複素数平面座標に変換する */
  toComplex(px: number, py: number, width: number, height: number) {
    return {
      x: this.centerX + (px - width / 2) * this.scale,
      y: this.centerY + (py - height / 2) * this.scale,
    };
  }

  /** 現在のビューポートの複素数範囲を返す */
  bounds(width: number, height: number) {
    return {
      xMin: this.centerX - (width / 2) * this.scale,
      xMax: this.centerX + (width / 2) * this.scale,
      yMin: this.centerY - (height / 2) * this.scale,
      yMax: this.centerY + (height / 2) * this.scale,
    };
  }

  /** ピクセル単位でパンする */
  pan(dx: number, dy: number) {
    this.centerX -= dx * this.scale;
    this.centerY -= dy * this.scale;
  }

  /** 指定した Canvas 座標を中心にズームする */
  zoom(factor: number, px: number, py: number, width: number, height: number) {
    const { x, y } = this.toComplex(px, py, width, height);
    this.scale *= factor;
    this.centerX = x + (this.centerX - x) * factor;
    this.centerY = y + (this.centerY - y) * factor;
  }
}
