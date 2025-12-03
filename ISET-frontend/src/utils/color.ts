export const hexToHSL = (hex: string, lightnessAdjustment: number) => {
  // 移除十六进制颜色值的前缀（如果有的话）
  let cleanHex = hex.startsWith("#") ? hex.slice(1) : hex;

  // 将三位的十六进制颜色扩展为六位
  if (cleanHex.length === 3) {
    cleanHex = cleanHex
      .split("")
      .map(char => char + char)
      .join("");
  }

  // 解析十六进制颜色为RGB
  const r = parseInt(cleanHex.substr(0, 2), 16) / 255;
  const g = parseInt(cleanHex.substr(2, 2), 16) / 255;
  const b = parseInt(cleanHex.substr(4, 2), 16) / 255;

  // 将RGB转换为HSL
  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  const l = (max + min) / 2;
  const d = max - min;
  const h =
    d === 0
      ? 0
      : max === r
      ? 60 * (((g - b) / d) % 6)
      : max === g
      ? 60 * ((b - r) / d + 2)
      : 60 * ((r - g) / d + 4);
  const s = l > 0 && l < 1 ? d / (1 - Math.abs(2 * l - 1)) : 0;

  // 应用深浅调整
  const adjustedLightness = Math.max(0, Math.min(1, l + lightnessAdjustment));

  // 确保亮度值在0到1之间，然后转换为百分比
  const hsl = `hsl(${Math.round(h)},${Math.round(s * 100)}%, ${Math.round(
    adjustedLightness * 100
  )}%)`;

  return hsl;
};
