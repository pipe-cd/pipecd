import { Box } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC, useMemo, useState } from "react";
import { sortedSet } from "~/utils/sorted-set";
import { ResourceState } from "~~/model/application_live_state_pb";
import { ResourceFilterPopover } from "../resource-filter-popover";
import { ResourceNode } from "./resource-node";
import { ResourceDetail } from "../resource-detail";
import dagre from "dagre";
import ResourceEdge from "./resource-edge";

const useStyles = makeStyles((theme) => ({
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
  detailPanel: {
    position: "absolute",
    right: theme.spacing(0),
    top: theme.spacing(0),
    bottom: theme.spacing(0),
    color: theme.palette.grey[500],
    zIndex: 100,
  },
  floatRight: {
    position: "absolute",
    right: 0,
    top: 0,
    zIndex: 10,
  },
}));

type Props = {
  resources: ResourceState.AsObject[];
};

const initFilterState = (
  resources: ResourceState.AsObject[]
): Record<string, boolean> => {
  const types: string[] = sortedSet(resources.map((r) => r.resourceType ?? ""));
  return types.reduce<Record<string, boolean>>((prev, current) => {
    prev[current] = true;
    return prev;
  }, {});
};

const NODE_HEIGHT = 72;
const NODE_WIDTH = 300;
const STROKE_WIDTH = 2;
const SVG_RENDER_PADDING = STROKE_WIDTH * 2;

function useGraph(
  resources: ResourceState.AsObject[],
  filterState: Record<string, boolean>
): dagre.graphlib.Graph<{
  resource: ResourceState.AsObject;
}> {
  const final = useMemo(() => {
    const graph = new dagre.graphlib.Graph<{
      resource: ResourceState.AsObject;
    }>();
    graph.setGraph({ rankdir: "LR", align: "UL" });
    graph.setDefaultEdgeLabel(() => ({}));

    const ignoreMap = resources.reduce<Record<string, boolean>>((prev, r) => {
      const resourceType = r.resourceType;

      prev[r.id] = !filterState[resourceType];
      return prev;
    }, {});

    resources.forEach((resource) => {
      if (ignoreMap[resource.id]) {
        return;
      }

      graph.setNode(resource.id, {
        resource,
        height: NODE_HEIGHT,
        width: NODE_WIDTH,
      });
      if (resource.parentIdsList.length > 0) {
        resource.parentIdsList.forEach((parentId) => {
          if (ignoreMap[parentId]) {
            return;
          }
          graph.setEdge(parentId, resource.id);
        });
      }
    });

    // Update after change graph
    dagre.layout(graph);

    return graph;
  }, [filterState, resources]);

  return final;
}

const GraphView: FC<Props> = ({ resources }) => {
  const classes = useStyles();
  const [filterState, setFilterState] = useState<Record<string, boolean>>(
    initFilterState(resources)
  );

  const [
    selectedResource,
    setSelectedResource,
  ] = useState<ResourceState.AsObject | null>(null);

  const graph = useGraph(resources, filterState);

  const nodes = graph
    .nodes()
    .map((v) => graph.node(v))
    .filter(Boolean);

  const edges = useMemo(() => {
    return graph.edges().map((v) => {
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

      svgWidth = Math.ceil(svgWidth - baseX) + SVG_RENDER_PADDING;
      svgHeight = Math.ceil(svgHeight - baseY) + SVG_RENDER_PADDING;

      const points = edge.points.reduce((prev, current) => {
        const x = Math.round(current.x - baseX) + STROKE_WIDTH;
        const y = Math.round(current.y - baseY) + STROKE_WIDTH;
        return prev + `${x},${y} `;
      }, "");

      return {
        points,
        top: baseY + NODE_HEIGHT / 2,
        left: baseX + NODE_WIDTH / 2,
        width: svgWidth,
        height: svgHeight,
      };
    });
  }, [graph]);

  const graphInstance = graph.graph();

  return (
    <div className={classes.stateViewWrapper}>
      <div className={classes.stateView}>
        {nodes.map((node) => (
          <Box
            key={`${node.resource.id}-${node.resource.name}`}
            position="absolute"
            top={node.y}
            left={node.x}
            zIndex={1}
            data-testid="application-resource"
          >
            <ResourceNode
              resource={node.resource}
              onClick={setSelectedResource}
            />
          </Box>
        ))}
        {edges?.map(({ top, left, width, height, points }, i) => (
          <ResourceEdge
            key={points + i}
            top={top}
            left={left}
            width={width}
            height={height}
            points={points}
          />
        ))}
        {graphInstance && (
          <div
            style={{
              width: (graphInstance.width ?? 0) + NODE_WIDTH,
              height: (graphInstance.height ?? 0) + NODE_HEIGHT,
            }}
          />
        )}
      </div>
      <Box className={classes.floatRight}>
        <ResourceFilterPopover
          filterState={filterState}
          onChange={(state) => setFilterState(state)}
        />
      </Box>
      {selectedResource && (
        <Box className={classes.detailPanel}>
          <ResourceDetail
            resource={selectedResource}
            onClose={() => setSelectedResource(null)}
          />
        </Box>
      )}
    </div>
  );
};

export default GraphView;
