import {
  Box,
  IconButton,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";
import clsx from "clsx";
import dagre from "dagre";
import React, { FC, useState } from "react";
import { KubernetesResourceState } from "../modules/applications-live-state";
import { theme } from "../theme";
import { KubernetesResource } from "./kubernetes-resource";

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

const DETAIL_WIDTH = 400;

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    flex: 1,
    justifyContent: "center",
  },
  stateView: {
    position: "relative",
    overflow: "scroll",
    height: "100%",
  },
  stateViewShift: {
    paddingRight: DETAIL_WIDTH,
  },
  detail: {
    width: DETAIL_WIDTH,
    padding: "16px 24px",
    height: "100%",
    position: "absolute",
    right: 0,
    zIndex: 2,
  },
  closeDetailButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
  detailName: {
    paddingRight: theme.spacing(4),
    wordBreak: "break-all",
    paddingBottom: theme.spacing(2),
  },
  detailSectionTitle: {
    color: theme.palette.text.secondary,
    minWidth: 120,
  },
  detailSection: {
    paddingTop: theme.spacing(1),
    display: "flex",
    alignItems: "center",
  },
  detailSectionBody: {
    flex: 1,
    wordBreak: "break-all",
  },
  multilineSection: {
    paddingTop: theme.spacing(1),
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
  const [
    selectedResource,
    setSelectedResource,
  ] = useState<KubernetesResourceState | null>(null);

  const graph = new dagre.graphlib.Graph<{
    resource: KubernetesResourceState;
  }>();
  graph.setGraph({ rankdir: "LR", align: "UL" });
  graph.setDefaultEdgeLabel(() => ({}));
  resources.forEach((resource) => {
    graph.setNode(resource.id, {
      resource,
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

  const graphInstance = graph.graph();

  return (
    <div
      className={clsx(classes.root, {
        [classes.stateViewShift]: selectedResource,
      })}
    >
      <div className={classes.stateView}>
        {nodes.map((node) => (
          <Box
            key={`${node.resource.kind}-${node.resource.name}`}
            position="absolute"
            top={node.y}
            left={node.x}
            zIndex={1}
          >
            <KubernetesResource
              resource={node.resource}
              onClick={setSelectedResource}
            />
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
        {graphInstance && (
          <div
            style={{
              width: (graphInstance.width ?? 0) + NODE_WIDTH,
              height: (graphInstance.height ?? 0) + NODE_HEIGHT,
            }}
          />
        )}
      </div>
      {selectedResource && (
        <Paper className={classes.detail} square>
          <IconButton
            className={classes.closeDetailButton}
            onClick={() => setSelectedResource(null)}
          >
            <CloseIcon />
          </IconButton>
          <Typography variant="h6" className={classes.detailName}>
            {selectedResource.name}
          </Typography>

          <div className={classes.detailSection}>
            <Typography
              variant="subtitle1"
              className={classes.detailSectionTitle}
            >
              Kind
            </Typography>
            <Typography variant="body1" className={classes.detailSectionBody}>
              {selectedResource.kind}
            </Typography>
          </div>

          <div className={classes.detailSection}>
            <Typography
              variant="subtitle1"
              className={classes.detailSectionTitle}
            >
              Namespace
            </Typography>
            <Typography variant="body1" className={classes.detailSectionBody}>
              {selectedResource.namespace}
            </Typography>
          </div>

          <div className={classes.detailSection}>
            <Typography
              variant="subtitle1"
              className={classes.detailSectionTitle}
            >
              Api Version
            </Typography>
            <Typography variant="body1" className={classes.detailSectionBody}>
              {selectedResource.apiVersion}
            </Typography>
          </div>

          <div className={classes.multilineSection}>
            <Typography
              variant="subtitle1"
              className={classes.detailSectionTitle}
            >
              Health Description
            </Typography>
            <Typography variant="body1">
              {selectedResource.healthDescription || "Empty"}
            </Typography>
          </div>
        </Paper>
      )}
    </div>
  );
};
