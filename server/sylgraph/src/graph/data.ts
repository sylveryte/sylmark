import type { SimulationLinkDatum, SimulationNodeDatum } from "d3";

export interface IGraphData {
  nodes: INode[];
  links: ILink[];
}

export const NodeKind = {
  File: 1,
  Tag: 2,
  UnresolvedFile: 3,
};

export interface INode extends SimulationNodeDatum {
  id: number;
  name: string;
  kind: number;
  val: number;
}

export interface ILink extends SimulationLinkDatum<INode> {}

export function genRandomTree(N = 300, reverse = false) {
  return {
    nodes: [...Array(N).keys()].map((i) => ({ id: i })),
    links: [...Array(N).keys()]
      .filter((id) => id)
      .map((id) => ({
        [reverse ? "target" : "source"]: id,
        [reverse ? "source" : "target"]: Math.round(Math.random() * (id - 1)),
      })),
  };
}
