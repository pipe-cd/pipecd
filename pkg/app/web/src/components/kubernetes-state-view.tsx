import { makeStyles, Box } from "@material-ui/core";
import React, { FC } from "react";
import { KubernetesResourceState } from "../modules/applications-live-state";
import { KubernetesResource } from "./kubernetes-resource";
import dagre from "dagre";

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

export const KubernetesStateView: FC<Props> = ({ resources }) => {
  const classes = useStyles();

  const graph = new dagre.graphlib.Graph<{ name: string; kind: string }>();
  graph.setGraph({ rankdir: "LR" });
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

  const nodes = graph.nodes().map((v) => graph.node(v));

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
        <Box position="absolute" top={node.y} left={node.x}>
          <KubernetesResource name={node.name} kind={node.kind} />
        </Box>
      ))}
      {graph.edges().map((v) => {
        const edge = graph.edge(v);
        let baseX = 10000;
        let baseY = 10000;
        edge.points.forEach((p) => {
          baseX = Math.round(Math.min(baseX, p.x));
          baseY = Math.round(Math.min(baseY, p.y));
        });
        return (
          <svg
            style={{
              position: "absolute",
              top: baseY + NODE_HEIGHT / 2,
              left: baseX + NODE_WIDTH / 2,
              zIndex: -1,
            }}
          >
            <polyline
              points={edge.points.reduce((prev, current) => {
                return (
                  prev +
                  `${Math.round(current.x - baseX + 1)},${Math.round(
                    current.y - baseY + 1
                  )} `
                );
              }, "")}
              strokeWidth={2}
              stroke="black"
              fill="white"
            />
          </svg>
        );
      })}
    </div>
  );
};
