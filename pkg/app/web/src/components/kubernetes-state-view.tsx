import { makeStyles, Box } from "@material-ui/core";
import React, { FC } from "react";
import { KubernetesResourceState } from "../modules/applications-live-state";
import { KubernetesResource } from "./kubernetes-resource";
import dagre from "dagre";
import { theme } from "../theme";

interface Line {
  baseX: number;
  baseY: number;
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  width: number;
  height: number;
}

const useStyles = makeStyles(() => ({
  container: {
    position: "relative",
  },
}));

interface Props {
  resources: KubernetesResourceState[];
}

const NODE_HEIGHT = 72;
const NODE_WIDTH = 300;
const STROKE_WIDTH = 2;
const SVG_RENDER_PADDING = STROKE_WIDTH * 2;

export const KubernetesStateView: FC<Props> = ({ resources }) => {
  const classes = useStyles();

  const graph = new dagre.graphlib.Graph<{ name: string; kind: string }>();
  graph.setGraph({ rankdir: "LR", align: "UL" });
  graph.setDefaultEdgeLabel(() => ({}));
  resources.forEach((resource) => {
    graph.setNode(resource.id, {
      name: resource.name,
      kind: resource.kind,
      height: NODE_HEIGHT,
      width: NODE_WIDTH,
    });
    if (resource.parentIdsList.length > 0) {
      resource.parentIdsList.forEach((parentId) => {
        graph.setEdge(parentId, resource.id);
      });
    }
  });

  // Update after change graph
  dagre.layout(graph);

  const nodes = graph
    .nodes()
    .map((v) => graph.node(v))
    .filter(Boolean);

  const lines: Line[] = [];
  graph.edges().forEach((v) => {
    const edge = graph.edge(v);
    if (edge.points.length > 1) {
      for (let i = 1; i < edge.points.length; i++) {
        const line: any = {
          baseX: Math.round(Math.min(edge.points[i - 1].x, edge.points[i].x)),
          baseY: Math.round(Math.min(edge.points[i - 1].y, edge.points[i].y)),
        };
        line.x1 = Math.round(edge.points[i - 1].x - line.baseX);
        line.y1 = Math.round(edge.points[i - 1].y - line.baseY);
        line.x2 = Math.round(edge.points[i].x - line.baseX);
        line.y2 = Math.round(edge.points[i].y - line.baseY);
        line.width = Math.max(line.x1, line.x2);
        line.height = Math.max(line.y1, line.y2);
        lines.push(line as Line);
      }
    }
  });

  return (
    <div className={classes.container}>
      {nodes.map((node) => (
        <Box key={node.name} position="absolute" top={node.y} left={node.x}>
          <KubernetesResource name={node.name} kind={node.kind} />
        </Box>
      ))}
      {graph.edges().map((v, i) => {
        const edge = graph.edge(v);
        let baseX = 10000;
        let baseY = 10000;
        let svgWidth = 0;
        let svgHeight = 0;
        edge.points.forEach((p) => {
          baseX = Math.min(baseX, p.x);
          baseY = Math.min(baseY, p.y);
          svgWidth = Math.max(svgWidth, p.x);
          svgHeight = Math.max(svgHeight, p.y);
        });
        baseX = Math.round(baseX);
        baseY = Math.round(baseY);
        // NOTE: Add padding to SVG sizes for showing edges completely.
        // If you use the same size as the polyline points, it may hide the some strokes.
        svgWidth = Math.ceil(svgWidth - baseX) + SVG_RENDER_PADDING;
        svgHeight = Math.ceil(svgHeight - baseY) + SVG_RENDER_PADDING;
        return (
          <svg
            key={`edge-${i}`}
            style={{
              position: "absolute",
              top: baseY + NODE_HEIGHT / 2,
              left: baseX + NODE_WIDTH / 2,
              zIndex: -1,
            }}
            width={svgWidth}
            height={svgHeight}
          >
            <polyline
              points={edge.points.reduce((prev, current) => {
                return (
                  prev +
                  `${Math.round(current.x - baseX) + STROKE_WIDTH},${
                    Math.round(current.y - baseY) + STROKE_WIDTH
                  } `
                );
              }, "")}
              strokeWidth={STROKE_WIDTH}
              stroke={theme.palette.divider}
              fill="transparent"
            />
          </svg>
        );
      })}
    </div>
  );
};
