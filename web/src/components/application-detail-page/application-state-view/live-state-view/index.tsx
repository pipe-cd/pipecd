import { Box, makeStyles } from "@material-ui/core";
import clsx from "clsx";
import { FC, useEffect, useMemo, useState } from "react";

import { ResourceState } from "~~/model/application_live_state_pb";
import DeploymentTabBar from "./deployment-tab-bar";
import GraphView from "./graph-view";

const useStyles = makeStyles(() => ({
  root: {
    display: "flex",
    flex: 1,
    justifyContent: "center",
    overflow: "hidden",
    position: "relative",
  },
  floatLeft: {
    position: "absolute",
    left: 0,
    top: 0,
    zIndex: 100,
  },
}));

type Props = {
  resources: ResourceState.AsObject[];
};

export const LiveStateView: FC<Props> = ({ resources: allResources }) => {
  const classes = useStyles();
  const [deploymentTabSelected, setDeploymentTabSelected] = useState("");

  const resourcesByDeployment = useMemo(() => {
    const resourceMapByDeployment: Record<
      string,
      ResourceState.AsObject[]
    > = {};
    allResources.forEach((r) => {
      const deployTarget = r.deployTarget;
      if (!deployTarget) return;

      if (!resourceMapByDeployment[deployTarget]) {
        resourceMapByDeployment[deployTarget] = [];
      }
      resourceMapByDeployment[deployTarget].push(r);
    });
    return resourceMapByDeployment;
  }, [allResources]);

  const deploymentList = useMemo(() => Object.keys(resourcesByDeployment), [
    resourcesByDeployment,
  ]);

  useEffect(() => {
    if (deploymentTabSelected == "" && deploymentList.length > 0) {
      setDeploymentTabSelected(deploymentList[0]);
    }
  }, [deploymentList, deploymentTabSelected]);

  const resources = useMemo(() => {
    if (!deploymentTabSelected) return [];
    return resourcesByDeployment[deploymentTabSelected];
  }, [deploymentTabSelected, resourcesByDeployment]);

  return (
    <div className={clsx(classes.root)}>
      <Box className={classes.floatLeft}>
        <DeploymentTabBar
          tabs={deploymentList}
          selectedTab={deploymentTabSelected}
          onSelectTab={setDeploymentTabSelected}
        />
      </Box>

      {deploymentTabSelected && (
        <GraphView key={deploymentTabSelected} resources={resources} />
      )}
    </div>
  );
};
