import bz from "bezier-easing";
import type { INode } from "./data";
import { setHexOpacity } from "./colors";

const easeFunction = bz(0.03, 0.65, 1, 0.42);
const SAFEMINSCALE = 0.01;
const SAFEMAXSCALE = 2.5;

export function getSafeScale(scale: number): number {
  let smoothScale = easeFunction(scale) * 2;
  if (smoothScale < SAFEMINSCALE) smoothScale = SAFEMINSCALE;
  if (smoothScale > SAFEMAXSCALE) smoothScale = SAFEMAXSCALE;
  return smoothScale;
}

// no need simulation.findNode works
export function findNode(nodes: INode[], x: number, y: number): INode | null {
  let i;
  for (i = nodes.length - 1; i >= 0; --i) {
    const node = nodes[i],
      dx = x - (node?.x || 0),
      dy = y - (node?.y || 0),
      distSq = dx * dx + dy * dy;
    const rSq = node.val * node.val * Math.PI;
    if (distSq < rSq) {
      return node;
    }
  }

  return null;
}

export function getTextStyle(
  color: string,
  scale: number,
  fadeFrom: number,
  fadeTill: number,
): [string, number] {
  if (scale > fadeTill) {
    return [color, 0];
  }
  const perc = (scale - fadeFrom) / (fadeTill - fadeFrom);
  return [setHexOpacity(color, perc), perc * 0.5 + 0.5];
}

export function getTextXY(
  text: string,
  x: number,
  y: number,
  radius: number,
  defaultFontSize: number,
): [string, number, number] {
  y += radius + defaultFontSize; //16 is font size

  return [text, x, y];
}

export function paintText(
  ctx: CanvasRenderingContext2D,
  x: number,
  y: number,
  width: number,
  _lineHeight: number,
  text: string,
  textSize: number,
) {
  if (text.length > widthToCharLength(width, textSize)) {
    ctx.fillText(text, x, y);
    return;
  }

  const splitted = text.split(" ");
  let tempText = "";
  for (let i = 0; i < splitted.length; i++) {
    if ((i = 0)) {
      tempText = splitted[0];
    } else {
      // todo

      ctx.fillText(text, x, y);
    }

    // determine if it's last and we should commit
    if (i == splitted.length - 1) {
      if (tempText.length + splitted[i]) {
      }
    }
  }
}

function widthToCharLength(width: number, textSize: number) {
  return Math.floor(width / textSize);
}
