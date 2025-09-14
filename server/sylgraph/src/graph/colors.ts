import { NodeKind, type INode } from "./data";

export type TColorFx = (alpha?: number) => string;
export interface IColors {
  background: string; // background
  primary: string; // node highlight
  secondary: string; // node
  tertiary: string; // links
  accent: string; // node highlight label
}
export function isPreferredThemeDark() {
  return window.matchMedia("(prefers-color-scheme: dark)").matches;
}

export function getColors(): IColors {
  return isPreferredThemeDark() ? colors.dark[0] : colors.light[0];
}

export function getNodeColor(
  color: string,
  alpha: number,
  node: INode,
): string {
  if (node.kind === NodeKind.UnresolvedFile && alpha > 0.2) {
    return setHexOpacity(color, 0.2);
  } else if (node.kind === NodeKind.Tag) {
    return setHexOpacity("#67C090", alpha);
  }
  return setHexOpacity(color, alpha);
}

export const colors = {
  dark: [
    {
      background: "#1c202a",
      primary: "#777777",
      secondary: "#0b8494",
      tertiary: "#888888",
      accent: "#F05A7e",
    },
  ],
  light: [
    {
      background: "#f1e2d9",
      primary: "#444444",
      secondary: "#20464d",
      tertiary: "#888888",
      accent: "#F05A7e",
    },
  ],
};

export const setHexOpacity = (hex: string, alpha: number) =>
  `${hex}${Math.floor(alpha * 255)
    .toString(16)
    .padStart(2, "0")}`;

console.log("cool", setHexOpacity("#777777", 50));
