import { Box } from "@mui/material";
import { FC, useEffect, useMemo, useState } from "react";

import { ResourceState } from "~~/model/application_live_state_pb";
import DeployTargetTabBar from "./deploy-target-tab-bar";
import GraphView from "./graph-view";

type Props = {
  resources: ResourceState.AsObject[];
};

export const LiveStateView: FC<Props> = ({ resources: allResources }) => {
  const [deployTargetTabSelected, setDeployTargetTabSelected] = useState("");

  const resourcesByDeployTarget = useMemo(() => {
    return allResources.reduce((all, resource) => {
      const deployTarget = resource.deployTarget;
      if (!deployTarget) return all;

      if (!all[deployTarget]) {
        all[deployTarget] = [];
      }
      all[deployTarget].push(resource);
      return all;
    }, {} as Record<string, ResourceState.AsObject[]>);
  }, [allResources]);

  const deployTargetList = useMemo(
    () => Object.keys(resourcesByDeployTarget).sort(),
    [resourcesByDeployTarget]
  );

  // Fall back to the first deploy target when nothing is selected yet, or when
  // the selected one is no longer in the live state (e.g. after a refresh).
  const selectedDeployTarget = deployTargetList.includes(
    deployTargetTabSelected
  )
    ? deployTargetTabSelected
    : deployTargetList[0] ?? "";

  // Remember the fallback, so a refresh that adds an earlier-sorted target
  // doesn't switch away from the one being shown.
  useEffect(() => {
    if (selectedDeployTarget !== deployTargetTabSelected) {
      setDeployTargetTabSelected(selectedDeployTarget);
    }
  }, [selectedDeployTarget, deployTargetTabSelected]);

  const resources = useMemo(() => {
    if (!selectedDeployTarget) return [];
    return resourcesByDeployTarget[selectedDeployTarget];
  }, [selectedDeployTarget, resourcesByDeployTarget]);

  return (
    <Box
      sx={{
        display: "flex",
        flex: 1,
        justifyContent: "center",
        overflow: "hidden",
        position: "relative",
      }}
    >
      <Box sx={{ position: "absolute", left: 0, top: 0, zIndex: 100 }}>
        <DeployTargetTabBar
          tabs={deployTargetList}
          selectedTab={selectedDeployTarget}
          onSelectTab={setDeployTargetTabSelected}
        />
      </Box>

      {selectedDeployTarget && (
        <GraphView key={selectedDeployTarget} resources={resources} />
      )}
    </Box>
  );
};
