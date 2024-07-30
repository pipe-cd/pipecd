import { Box, makeStyles } from "@material-ui/core";
import clsx from "clsx";
import dagre from "dagre";
import { FC, useState } from "react";
import { ECSResourceState } from "~/modules/applications-live-state";
import { theme } from "~/theme";
import { ECSResource } from "./ecs-resource";
import { ECSResourceDetail } from "./ecs-resource-detail";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    flex: 1,
    justifyContent: "center",
    overflow: "hidden",
  },
  stateViewWrapper: {
    flex: 1,
    display: "flex",
    justifyContent: "center",
    overflow: "hidden",
  },
  stateView: {
    position: "relative",
    overflow: "auto",
  },
  closeDetailButton: {
    position: "absolute",
    right: theme.spacing(1),
    top: theme.spacing(1),
    color: theme.palette.grey[500],
  },
}));

export interface ECSStateViewProps {
  resources: ECSResourceState.AsObject[];
}

const NODE_HEIGHT = 72;
const NODE_WIDTH = 300;
const STROKE_WIDTH = 2;
const SVG_RENDER_PADDING = STROKE_WIDTH * 2;

function useGraph(
  resources: ECSResourceState.AsObject[]
): dagre.graphlib.Graph<{
  resource: ECSResourceState.AsObject;
}> {
  const graph = new dagre.graphlib.Graph<{
    resource: ECSResourceState.AsObject;
  }>();
  graph.setGraph({ rankdir: "LR", align: "UL" });
  graph.setDefaultEdgeLabel(() => ({}));

  resources.forEach((resource) => {
    graph.setNode(resource.id, {
      resource,
      height: NODE_HEIGHT,
      width: NODE_WIDTH,
    });
    // 'Service' does not need parent nodes.
    if (resource.kind != "Service" && resource.parentIdsList.length > 0) {
      resource.parentIdsList.forEach((parentId) => {
        graph.setEdge(parentId, resource.id);
      });
    }
  });

  // Update after change graph
  dagre.layout(graph);

  return graph;
}

export const ECSStateView: FC<ECSStateViewProps> = ({ resources }) => {
  const classes = useStyles();
  const [
    selectedResource,
    setSelectedResource,
  ] = useState<ECSResourceState.AsObject | null>(null);

  const graph = useGraph(resources);
  const nodes = graph
    .nodes()
    .map((v) => graph.node(v))
    .filter(Boolean);

  const graphInstance = graph.graph();

  return (
    <div className={clsx(classes.root)}>
      <div className={classes.stateViewWrapper}>
        <div className={classes.stateView}>
          {nodes.map((node) => (
            <Box
              key={`${node.resource.kind}-${node.resource.name}`}
              position="absolute"
              top={node.y}
              left={node.x}
              zIndex={1}
              data-testid="ecs-resource"
            >
              <ECSResource
                resource={node.resource}
                onClick={setSelectedResource}
              />
            </Box>
          ))}
          {
            // render edges
            graph.edges().map((v, i) => {
              const edge = graph.edge(v);
              let baseX = Infinity;
              let baseY = Infinity;
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
            })
          }
          {graphInstance && (
            <div
              style={{
                width: (graphInstance.width ?? 0) + NODE_WIDTH,
                height: (graphInstance.height ?? 0) + NODE_HEIGHT,
              }}
            />
          )}
        </div>
      </div>

      {selectedResource && (
        <ECSResourceDetail
          resource={selectedResource}
          onClose={() => setSelectedResource(null)}
        />
      )}
    </div>
  );
};
