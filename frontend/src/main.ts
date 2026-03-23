import { Viewport } from "./viewport";
import { TileManager } from "./tileManager";
import { loadFromUrl, saveToUrl, DEFAULT_STATE } from "./params";

const canvas = document.getElementById("canvas") as HTMLCanvasElement;
const ctx = canvas.getContext("2d")!;
const colorSchemeEl = document.getElementById("colorScheme") as HTMLSelectElement;
const maxIterEl = document.getElementById("maxIter") as HTMLInputElement;
const resetBtn = document.getElementById("resetBtn") as HTMLButtonElement;

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

// ビューポートの現在状態を URL に保存する
function syncUrl() {
  saveToUrl({
    cx: viewport.centerX,
    cy: viewport.centerY,
    scale: viewport.scale,
    colorScheme: colorSchemeEl.value,
    maxIter: parseInt(maxIterEl.value, 10),
  });
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
  syncUrl();
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
  syncUrl();
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
  syncUrl();
}, { passive: false });

canvas.addEventListener("touchend", () => {
  lastTouches = null;
});

// ---- コントロール ----

colorSchemeEl.addEventListener("change", () => {
  tileManager.setColorScheme(colorSchemeEl.value);
  render();
  syncUrl();
});

maxIterEl.addEventListener("change", () => {
  const v = parseInt(maxIterEl.value, 10);
  if (!isNaN(v) && v > 0) {
    tileManager.setMaxIter(v);
    render();
    syncUrl();
  }
});

// リセットボタン: 初期状態に戻す
resetBtn.addEventListener("click", () => {
  viewport.centerX = DEFAULT_STATE.cx;
  viewport.centerY = DEFAULT_STATE.cy;
  viewport.scale = DEFAULT_STATE.scale;

  colorSchemeEl.value = DEFAULT_STATE.colorScheme;
  maxIterEl.value = DEFAULT_STATE.maxIter.toString();

  tileManager.setColorScheme(DEFAULT_STATE.colorScheme);
  tileManager.setMaxIter(DEFAULT_STATE.maxIter);

  render();
  syncUrl();
});

// ---- 初期化: URL パラメータから状態を復元 ----

const initialState = loadFromUrl();
viewport.centerX = initialState.cx;
viewport.centerY = initialState.cy;
viewport.scale = initialState.scale;

colorSchemeEl.value = initialState.colorScheme;
maxIterEl.value = initialState.maxIter.toString();

tileManager.setColorScheme(initialState.colorScheme);
tileManager.setMaxIter(initialState.maxIter);

window.addEventListener("resize", resize);
resize();
syncUrl();
