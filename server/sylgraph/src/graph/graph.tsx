import { useEffect, useRef, useState } from "react";
import * as d3 from "d3";
import type { IGraphData, ILink, INode } from "./data";
import { getSafeScale, getTextStyle, getTextXY } from "./utils";
import { getColors, getNodeColor, setHexOpacity } from "./colors";

interface OwnProps {
  data: IGraphData | undefined;
  openDoc: (id: number) => void;
}

const fadeFrom = 2;
const fadeTill = 5;
const magnateScaler = 2;
const defaultFontSize = 18;
const unfocusedNodeAlpha = 0.2;
const focusedTransitionDuration = 700;
const defaultFontFx = (size: number = defaultFontSize) => `${size}px Sans`;
const defaultFont = defaultFontFx();

export const Graph = ({ data, openDoc }: OwnProps) => {
  const dref = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [nodes, setNodes] = useState<INode[]>([]);
  const hoveredNode = useRef<INode | undefined>(undefined);
  const [links, setLinks] = useState<ILink[]>([]);
  const transformRef = useRef(d3.zoomIdentity);

  const [dim, setDim] = useState<{ width: number; height: number }>({
    width: 500,
    height: 500,
  });

  const handleResize = () => {
    const r = dref.current?.getClientRects();
    if (r?.length) {
      const dr = r[0];
      setDim({ height: dr.height, width: dr.width });
    }
  };

  useEffect(() => {
    window.addEventListener("resize", handleResize);
    setTimeout(() => {
      handleResize();
    });
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    if (data) {
      setNodes(data.nodes);
      setLinks(data.links);
    }
  }, [data]);

  useEffect(() => {
    // Setup d3 force simulation
    // d
    const centerX = dim.width / 2,
      centerY = dim.height / 2;
    const simulation = d3
      .forceSimulation(nodes)
      .alpha(1)
      .alphaDecay(0.005)
      .force(
        "link",
        d3.forceLink<INode, ILink>(links).id((d) => d.id),
      )
      .force("charge", d3.forceManyBody().strength(-10))
      .force("collide", d3.forceCollide().radius(10))
      // .force("radial", d3.forceRadial(centerX, centerY).strength(0.01))
      .force("x", d3.forceX(0).strength(0.01))
      .force("y", d3.forceY(0).strength(0.01))
      .force("center", d3.forceCenter(centerX, centerY).strength(0.9)); // center of the canvas

    // hover active states &&  transition time track
    let prevRenderTime: Date = new Date();
    let delta: number = 0;
    let hoveredRenderElapsed: number = 0;
    let elaspedPercentage: number = 0;
    let hoveredNodesMap: Map<number, boolean> = new Map();
    const updateHoveredMap = (node: INode) => {
      const id = node.id;

      const m = new Map<number, boolean>();
      links.forEach((link) => {
        const sId = (link.source as INode).id;
        const tId = (link.target as INode).id;
        if (sId === id) {
          m.set(tId, true);
        } else if (tId === id) {
          m.set(sId, true);
        }
      });
      hoveredNodesMap = m;
    };

    let dragOffsetX = 0;
    let dragOffsetY = 0;

    const onClick = (event: MouseEvent) => {
      const rect = canvasRef.current?.getBoundingClientRect();

      const canvasX = event.x - (rect?.left || 0);
      const canvasY = event.y - (rect?.top || 0);

      const graphX = transformRef.current.invertX(canvasX);
      const graphY = transformRef.current.invertY(canvasY);

      const clickedNode = simulation.find(
        graphX,
        graphY,
        getSafeScale(transformRef.current.k) * magnateScaler,
      );
      if (clickedNode) {
        openDoc(clickedNode.id);
      }
    };
    const mouseHover = (event: MouseEvent) => {
      const rect = canvasRef.current?.getBoundingClientRect();

      const canvasX = event.x - (rect?.left || 0);
      const canvasY = event.y - (rect?.top || 0);

      const graphX = transformRef.current.invertX(canvasX);
      const graphY = transformRef.current.invertY(canvasY);

      const newHoveredNode = simulation.find(
        graphX,
        graphY,
        getSafeScale(transformRef.current.k) * magnateScaler,
      );
      if (newHoveredNode) {
        if (hoveredNode.current) {
          if (hoveredNode.current.id !== newHoveredNode.id) {
            updateHoveredMap(newHoveredNode);
          }
        } else {
          updateHoveredMap(newHoveredNode);
          hoveredNode.current = newHoveredNode;
        }
      } else {
        hoveredNode.current = newHoveredNode;
      }
    };

    // Reheat the simulation when drag starts, and fix the subject position.
    const dragstarted = (event: {
      active: any;
      x: number;
      y: number;
      subject: { x: any; y: any; fx: any; fy: any };
    }) => {
      if (!event.active) simulation.alphaTarget(0.3).restart();

      const mouseX = transformRef.current.invertX(event.x);
      const mouseY = transformRef.current.invertY(event.y);

      // Node position in simulation space
      const nodeX = event.subject.x;
      const nodeY = event.subject.y;

      // Store the offset between mouse and node
      dragOffsetX = nodeX - mouseX;
      dragOffsetY = nodeY - mouseY;

      event.subject.fx = nodeX;
      event.subject.fy = nodeY;
    };

    // Update the subject (dragged node) position during drag.
    const dragged = (event: {
      x: number;
      y: number;
      subject: { fx: number; fy: number };
    }) => {
      // Adjust the drag position for zoom scale
      const adjustedX = transformRef.current.invertX(event.x);
      const adjustedY = transformRef.current.invertY(event.y);

      // Update the node's fixed position (fx, fy)
      event.subject.fx = adjustedX + dragOffsetX;
      event.subject.fy = adjustedY + dragOffsetY;
    };

    // Restore the target alpha so the simulation cools after dragging ends.
    // Unfix the subject position now that it’s no longer being dragged.
    const dragended = (event: {
      active: any;
      subject: { fx: null; fy: null };
    }) => {
      if (!event.active) simulation.alphaTarget(0);
      event.subject.fx = null;
      event.subject.fy = null;
    };

    const dragSubject = (event: { x: number; y: number }) => {
      // Invert the mouse coordinates back to the untransformed space
      const x = transformRef.current.invertX(event.x);
      const y = transformRef.current.invertY(event.y);

      // Find the node based on the adjusted (zoomed) coordinates
      return simulation.find(
        x,
        y,
        getSafeScale(transformRef.current.k) * magnateScaler,
      );
    };

    //Setup Zoom Behavior
    const zoom = d3
      .zoom()
      .scaleExtent([0.1, 12]) // Limit zooming scale
      .on("zoom", (event) => {
        transformRef.current = event.transform;
      });

    // Apply zoom behavior to canvas
    if (canvasRef?.current) {
      const selection = d3.select(canvasRef.current as Element);

      // drag
      selection
        .call(
          d3
            .drag()
            .container(canvasRef?.current)
            .subject(dragSubject)
            .on("start", dragstarted)
            .on("drag", dragged)
            .on("end", dragended),
        )
        .call(zoom)
        .call(zoom.transform, d3.zoomIdentity);
    }

    const ctx = canvasRef?.current?.getContext("2d");

    if (ctx) {
      const colors = getColors();
      canvasRef.current!.style.background = colors.background;
      ctx.font = defaultFont;

      let alpha = 1;
      const render = () => {
        ctx.textAlign = "center";
        ctx.textBaseline = "middle";
        const isNodeHoverActive = !!hoveredNode?.current;
        // time track
        const now = new Date();
        if (
          isNodeHoverActive
            ? hoveredRenderElapsed < focusedTransitionDuration
            : hoveredRenderElapsed > 0
        ) {
          // elapsedTimePercentageCalcFocused
          elaspedPercentage =
            (hoveredRenderElapsed / focusedTransitionDuration) * 1;

          alpha = 1 - elaspedPercentage;
          if (alpha < unfocusedNodeAlpha) {
            alpha = unfocusedNodeAlpha;
          }

          delta = now.getMilliseconds() - prevRenderTime.getMilliseconds();
          if (delta > 0) {
            if (isNodeHoverActive) {
              hoveredRenderElapsed += delta;
            } else {
              hoveredRenderElapsed -= delta;
            }
          }
        }
        prevRenderTime = now;

        ctx.clearRect(
          0,
          0,
          canvasRef.current!.width,
          canvasRef.current!.height,
        ); // clear canvas

        // Draw links
        ctx.beginPath();

        const { x: tx, y: ty, k: scale } = transformRef.current;
        const safeScale = getSafeScale(scale);

        links.forEach((link) => {
          if ((link.source as INode)?.id && (link.target as INode)?.id) {
            const sourceNode = link.source as INode,
              targetNode = link.target as INode;
            // Apply zoom/pan transform to node positions
            const sx = sourceNode.x! * scale + tx;
            const sy = sourceNode.y! * scale + ty;
            const tx_ = targetNode.x! * scale + tx;
            const ty_ = targetNode.y! * scale + ty;

            ctx.stroke();
            ctx.beginPath();
            ctx.moveTo(sx, sy);
            ctx.lineTo(tx_, ty_);
            if (
              hoveredNode.current &&
              (hoveredNode.current.id === targetNode.id ||
                hoveredNode.current.id === sourceNode.id)
            ) {
              ctx.strokeStyle = colors.primary;
            } else {
              ctx.lineWidth = 0.4;
              ctx.strokeStyle = setHexOpacity(colors.tertiary, alpha);
            }
          }
        });

        // Draw nodes
        let labelDrawParams: [string, number, number] | undefined;
        nodes.forEach((node) => {
          const nx = node.x! * scale + tx;
          const ny = node.y! * scale + ty;
          ctx.beginPath();
          const radius = safeScale * node.val;
          ctx.arc(nx, ny, radius, 0, 2 * Math.PI);
          const isCurrentNodeHovered = node.id === hoveredNode.current?.id;
          if (isCurrentNodeHovered) {
            labelDrawParams = getTextXY(
              node.name,
              nx,
              ny,
              radius,
              defaultFontSize,
            );
            ctx.fillStyle = colors.accent;
            ctx.arc(nx, ny, radius + 2, 0, 2 * Math.PI);
            ctx.fill();
            ctx.beginPath();
            ctx.fillStyle = colors.primary;
            ctx.arc(nx, ny, radius, 0, 2 * Math.PI);
            ctx.fill();
          } else if (hoveredNodesMap.has(node.id)) {
            ctx.fillStyle = colors.secondary;
            ctx.fill();
          } else {
            ctx.fillStyle = getNodeColor(colors.secondary, alpha, node);
            // ctx.fillStyle = "pink";
            ctx.fill();
          }
          // label and ignore if hoveredNode
          if (scale > fadeFrom && !isCurrentNodeHovered) {
            ctx.beginPath();
            let fontScale;
            [ctx.fillStyle, fontScale] = getTextStyle(
              colors.tertiary,
              scale,
              fadeFrom,
              fadeTill,
            );
            ctx.font = fontScale
              ? defaultFontFx(Math.ceil(defaultFontSize * fontScale))
              : defaultFont;
            ctx.fillText(
              ...getTextXY(node.name, nx, ny, radius, defaultFontSize),
            );
            ctx.fill();
          }
        });
        // draw hoveredNode label
        if (labelDrawParams) {
          ctx.font = defaultFont;
          ctx.beginPath();
          ctx.fillStyle = colors.accent;
          ctx.fillText(...labelDrawParams);
          ctx.fill();
        }
      };

      let requestAnimationId: number;
      const tick = () => {
        render();
        requestAnimationId = requestAnimationFrame(tick); // Continue the simulation
      };

      tick(); // Start the simulation loop

      addEventListener("mousemove", mouseHover);
      addEventListener("click", onClick);

      return () => {
        removeEventListener("mousemove", mouseHover);
        removeEventListener("click", onClick);
        // Clean up the simulation on component unmount
        simulation.stop();
        cancelAnimationFrame(requestAnimationId);
      };
    }
  }, [nodes, links]);

  return (
    <div ref={dref} className="h-full">
      <canvas width={dim.width} height={dim.height} ref={canvasRef}></canvas>
    </div>
  );
};
