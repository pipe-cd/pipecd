import { makeStyles } from "@material-ui/core";
import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router";
import { DeploymentDetail } from "../../components/deployment-detail";
import { LogViewer } from "../../components/log-viewer";
import { Pipeline } from "../../components/pipeline";
import { AppState } from "../../modules";
import {
  Deployment,
  fetchDeploymentById,
  isDeploymentRunning,
  selectById as selectDeploymentById,
} from "../../modules/deployments";
import { useInterval } from "../../utils/use-interval";
import { clearActiveStage } from "../../modules/active-stage";

const FETCH_INTERVAL = 4000;

const useStyles = makeStyles({
  root: {
    display: "flex",
    flexDirection: "column",
    alignItems: "stretch",
    flex: 1,
  },
  main: {
    flex: 1,
  },
  loading: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
});

export const DeploymentDetailPage: FC = memo(function DeploymentDetailPage() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const { deploymentId } = useParams<{ deploymentId: string }>();
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectDeploymentById(state.deployments, deploymentId)
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
      <div className={classes.main}>
        <DeploymentDetail deploymentId={deploymentId} />
        <Pipeline deploymentId={deploymentId} />
      </div>
      <LogViewer />
    </div>
  );
});
