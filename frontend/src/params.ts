/** URL GET パラメータとして保持するアプリ状態 */
export interface AppState {
  cx: number;
  cy: number;
  scale: number;
  colorScheme: string;
  maxIter: number;
}

/** デフォルトの初期状態 */
export const DEFAULT_STATE: AppState = {
  cx: -0.5,
  cy: 0,
  scale: 0.004,
  colorScheme: "classic",
  maxIter: 256,
};

/** 現在の URL パラメータからアプリ状態を復元する。不正値はデフォルトにフォールバックする */
export function loadFromUrl(): AppState {
  const p = new URLSearchParams(window.location.search);

  const cx = parseFloat(p.get("cx") ?? "");
  const cy = parseFloat(p.get("cy") ?? "");
  const scale = parseFloat(p.get("scale") ?? "");
  const maxIter = parseInt(p.get("max_iter") ?? "", 10);
  const colorScheme = p.get("color_scheme") ?? "";

  return {
    cx: isFinite(cx) ? cx : DEFAULT_STATE.cx,
    cy: isFinite(cy) ? cy : DEFAULT_STATE.cy,
    scale: isFinite(scale) && scale > 0 ? scale : DEFAULT_STATE.scale,
    maxIter: isFinite(maxIter) && maxIter > 0 ? maxIter : DEFAULT_STATE.maxIter,
    colorScheme: colorScheme !== "" ? colorScheme : DEFAULT_STATE.colorScheme,
  };
}

/** アプリ状態を URL パラメータに書き込む（履歴は増やさない） */
export function saveToUrl(state: AppState) {
  const p = new URLSearchParams({
    cx: state.cx.toString(),
    cy: state.cy.toString(),
    scale: state.scale.toString(),
    color_scheme: state.colorScheme,
    max_iter: state.maxIter.toString(),
  });
  history.replaceState(null, "", `?${p.toString()}`);
}
