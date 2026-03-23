import { Viewport } from "./viewport";
import { TileManager } from "./tileManager";

const canvas = document.getElementById("canvas") as HTMLCanvasElement;
const ctx = canvas.getContext("2d")!;
const colorSchemeEl = document.getElementById("colorScheme") as HTMLSelectElement;
const maxIterEl = document.getElementById("maxIter") as HTMLInputElement;

const viewport = new Viewport();
const tileManager = new TileManager();

// Canvas をウィンドウサイズに合わせる
function resize() {
  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;
  render();
}

// タイルを描画する
function render() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  tileManager.render(ctx, viewport, canvas.width, canvas.height, render);
}

// ---- マウス操作 ----

let dragging = false;
let lastX = 0;
let lastY = 0;

canvas.addEventListener("mousedown", (e) => {
  dragging = true;
  lastX = e.clientX;
  lastY = e.clientY;
  canvas.classList.add("grabbing");
});

window.addEventListener("mousemove", (e) => {
  if (!dragging) return;
  viewport.pan(e.clientX - lastX, e.clientY - lastY);
  lastX = e.clientX;
  lastY = e.clientY;
  render();
});

window.addEventListener("mouseup", () => {
  dragging = false;
  canvas.classList.remove("grabbing");
});

canvas.addEventListener("wheel", (e) => {
  e.preventDefault();
  const factor = e.deltaY < 0 ? 0.8 : 1.25;
  viewport.zoom(factor, e.clientX, e.clientY, canvas.width, canvas.height);
  render();
}, { passive: false });

// ---- タッチ操作 ----

let lastTouches: TouchList | null = null;

canvas.addEventListener("touchstart", (e) => {
  e.preventDefault();
  lastTouches = e.touches;
}, { passive: false });

canvas.addEventListener("touchmove", (e) => {
  e.preventDefault();
  if (!lastTouches) return;

  if (e.touches.length === 1 && lastTouches.length === 1) {
    // 1本指: パン
    const dx = e.touches[0].clientX - lastTouches[0].clientX;
    const dy = e.touches[0].clientY - lastTouches[0].clientY;
    viewport.pan(dx, dy);
  } else if (e.touches.length === 2 && lastTouches.length === 2) {
    // 2本指: ピンチズーム
    const prevDist = Math.hypot(
      lastTouches[0].clientX - lastTouches[1].clientX,
      lastTouches[0].clientY - lastTouches[1].clientY
    );
    const currDist = Math.hypot(
      e.touches[0].clientX - e.touches[1].clientX,
      e.touches[0].clientY - e.touches[1].clientY
    );
    const factor = prevDist / currDist;
    const cx = (e.touches[0].clientX + e.touches[1].clientX) / 2;
    const cy = (e.touches[0].clientY + e.touches[1].clientY) / 2;
    viewport.zoom(factor, cx, cy, canvas.width, canvas.height);
  }

  lastTouches = e.touches;
  render();
}, { passive: false });

canvas.addEventListener("touchend", () => {
  lastTouches = null;
});

// ---- コントロール ----

colorSchemeEl.addEventListener("change", () => {
  tileManager.setColorScheme(colorSchemeEl.value);
  render();
});

maxIterEl.addEventListener("change", () => {
  const v = parseInt(maxIterEl.value, 10);
  if (!isNaN(v) && v > 0) {
    tileManager.setMaxIter(v);
    render();
  }
});

// ---- 初期化 ----

window.addEventListener("resize", resize);
resize();
