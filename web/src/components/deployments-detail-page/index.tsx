import { Box, makeStyles } from "@material-ui/core";
import { FC, memo, useEffect } from "react";
import { useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import { clearActiveStage } from "~/modules/active-stage";
import {
  Deployment,
  fetchDeploymentById,
  isDeploymentRunning,
  selectById as selectDeploymentById,
} from "~/modules/deployments";
import { DeploymentDetail } from "./deployment-detail";
import { LogViewer } from "./log-viewer";
import { Pipeline } from "./pipeline";

const FETCH_INTERVAL = 4000;

const useStyles = makeStyles({
  root: {
    display: "flex",
    flexDirection: "column",
    alignItems: "stretch",
    flex: 1,
    overflow: "auto",
  },
});

export const DeploymentDetailPage: FC = memo(function DeploymentDetailPage() {
  const classes = useStyles();
  const dispatch = useAppDispatch();
  const { deploymentId } = useParams<{ deploymentId: string }>();
  const deployment = useAppSelector<Deployment.AsObject | undefined>((state) =>
    selectDeploymentById(state.deployments, deploymentId ?? "")
  );

  const fetchData = (): void => {
    if (deploymentId) {
      dispatch(fetchDeploymentById(deploymentId));
    }
  };

  useEffect(fetchData, [dispatch, deploymentId]);
  useInterval(
    fetchData,
    deploymentId && isDeploymentRunning(deployment?.status)
      ? FETCH_INTERVAL
      : null
  );

  // NOTE: Clear active stage when leave detail page
  useEffect(
    () => () => {
      dispatch(clearActiveStage());
    },
    [dispatch]
  );

  return (
    <div className={classes.root}>
      <Box flex={1}>
        <DeploymentDetail deploymentId={deploymentId ?? ""} />
        <Pipeline deploymentId={deploymentId ?? ""} />
      </Box>
      <LogViewer />
    </div>
  );
});
