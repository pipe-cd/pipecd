import {
  Box,
  Button,
  Link,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import SyncIcon from "@material-ui/icons/Cached";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import Skeleton from "@material-ui/lab/Skeleton/Skeleton";
import { SerializedError } from "@reduxjs/toolkit";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { APPLICATION_HEALTH_STATUS_TEXT } from "../constants/health-status-text";
import { PAGE_PATH_DEPLOYMENTS } from "../constants/path";
import { UI_TEXT_REFRESH } from "../constants/ui-text";
import { AppState } from "../modules";
import {
  Application,
  ApplicationDeploymentReference,
  fetchApplication,
  selectById as selectApplicationById,
  syncApplication,
} from "../modules/applications";
import {
  ApplicationLiveState,
  selectById as selectLiveStateById,
} from "../modules/applications-live-state";
import { SyncStrategy } from "../modules/deployments";
import {
  Environment,
  selectById as selectEnvById,
} from "../modules/environments";
import { Piped, selectById as selectPipeById } from "../modules/pipeds";
import { DetailTableRow } from "./detail-table-row";
import { ApplicationHealthStatusIcon } from "./health-status-icon";
import { SplitButton } from "./split-button";
import { SyncStateReason } from "./sync-state-reason";
import { SyncStatusIcon } from "./sync-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    display: "flex",
    zIndex: theme.zIndex.appBar,
    position: "relative",
    flexDirection: "column",
  },
  content: {
    flex: 1,
  },
  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  syncStatusText: {
    marginLeft: theme.spacing(0.5),
    marginRight: theme.spacing(1),
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
  latestDeploymentTable: {
    paddingLeft: theme.spacing(2),
  },
  latestDeploymentLink: {
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  applicationId: string;
}

const useIsSyncingApplication = (
  applicationId: string | undefined
): boolean => {
  return useSelector<AppState, boolean>((state) => {
    if (!applicationId) {
      return false;
    }

    return state.applications.syncing[applicationId];
  });
};

const ERROR_MESSAGE = "It was unable to fetch the application.";

const MostRecentlySuccessfulDeployment: FC<{
  deployment?: ApplicationDeploymentReference.AsObject;
}> = ({ deployment }) => {
  const classes = useStyles();

  if (!deployment) {
    return <Skeleton height={105} width={500} />;
  }

  const date = dayjs(deployment.startedAt * 1000);

  return (
    <>
      <Box display="flex" alignItems="baseline">
        <Typography variant="subtitle1">Latest Deployment</Typography>
        <Typography variant="body2" className={classes.latestDeploymentLink}>
          <Link
            component={RouterLink}
            to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.deploymentId}`}
          >
            more details
          </Link>
        </Typography>
      </Box>
      <table className={classes.latestDeploymentTable}>
        <tbody>
          <DetailTableRow
            label="Deployed At"
            value={<span title={date.format()}>{date.fromNow()}</span>}
          />
          <DetailTableRow label="Version" value={deployment.version} />
          <DetailTableRow label="Summary" value={deployment.summary} />
        </tbody>
      </table>
    </>
  );
};

const syncOptions = ["Sync", "Quick Sync", "Pipeline Sync"];
const syncStrategyByIndex: SyncStrategy[] = [
  SyncStrategy.AUTO,
  SyncStrategy.QUICK_SYNC,
  SyncStrategy.PIPELINE,
];

export const ApplicationDetail: FC<Props> = memo(function ApplicationDetail({
  applicationId,
}) {
  const classes = useStyles();
  const dispatch = useDispatch();

  const [app, liveState, fetchApplicationError] = useSelector<
    AppState,
    [
      Application | undefined,
      ApplicationLiveState | undefined,
      SerializedError | null
    ]
  >((state) => [
    selectApplicationById(state.applications, applicationId),
    selectLiveStateById(state.applicationLiveState, applicationId),
    state.applications.fetchApplicationError,
  ]);

  const [pipe, env] = useSelector<
    AppState,
    [Environment | undefined, Piped | undefined]
  >((state) => [
    app ? selectEnvById(state.environments, app.envId) : undefined,
    app ? selectPipeById(state.pipeds, app.pipedId) : undefined,
  ]);

  const isSyncing = useIsSyncingApplication(app?.id);

  const handleSync = (index: number): void => {
    if (app) {
      dispatch(
        syncApplication({
          applicationId: app.id,
          syncStrategy: syncStrategyByIndex[index],
        })
      );
    }
  };

  if (fetchApplicationError) {
    return (
      <Paper square elevation={1} className={classes.root}>
        <Box
          height={200}
          display="flex"
          flexDirection="column"
          justifyContent="center"
          alignItems="center"
        >
          <Typography variant="body1">{ERROR_MESSAGE}</Typography>
          <Button
            color="primary"
            onClick={() => {
              dispatch(fetchApplication(applicationId));
            }}
          >
            {UI_TEXT_REFRESH}
          </Button>
        </Box>
      </Paper>
    );
  }

  return (
    <Paper square elevation={1} className={classes.root}>
      <Box flex={1}>
        <Box display="flex" alignItems="baseline">
          <Typography variant="h5">
            {app ? app.name : <Skeleton width={100} />}
          </Typography>
          <Typography variant="subtitle2" className={classes.env}>
            {env ? env.name : <Skeleton width={100} />}
          </Typography>
        </Box>

        {app?.syncState ? (
          <>
            <Box display="flex" alignItems="center">
              <SyncStatusIcon status={app.syncState.status} />
              <Box display="flex" alignItems="baseline">
                <Typography variant="h6" className={classes.syncStatusText}>
                  {APPLICATION_SYNC_STATUS_TEXT[app.syncState.status]}
                </Typography>
              </Box>

              {liveState ? (
                <>
                  <ApplicationHealthStatusIcon
                    health={liveState.healthStatus}
                  />
                  <Typography variant="h6" className={classes.syncStatusText}>
                    {APPLICATION_HEALTH_STATUS_TEXT[liveState.healthStatus]}
                  </Typography>
                </>
              ) : (
                <Skeleton height={32} width={100} />
              )}
            </Box>

            <SyncStateReason
              summary={app.syncState.shortReason}
              detail={app.syncState.reason}
            />
          </>
        ) : (
          <Skeleton height={32} width={200} />
        )}
      </Box>

      <Box mt={1} display="flex">
        <div className={classes.content}>
          {app && pipe ? (
            <table>
              <tbody>
                <DetailTableRow
                  label="Kind"
                  value={APPLICATION_KIND_TEXT[app.kind]}
                />
                <DetailTableRow label="Piped" value={pipe.name} />
                <DetailTableRow
                  label="Cloud Provider"
                  value={app.cloudProvider}
                />

                {app.gitPath && (
                  <DetailTableRow
                    label="Configuration Directory"
                    value={
                      <Link
                        href={app.gitPath.url}
                        target="_blank"
                        rel="noreferrer"
                      >
                        {app.gitPath.path}
                        <OpenInNewIcon className={classes.linkIcon} />
                      </Link>
                    }
                  />
                )}
              </tbody>
            </table>
          ) : (
            <Skeleton height={105} width={500} />
          )}
        </div>

        <div className={classes.content}>
          <MostRecentlySuccessfulDeployment
            deployment={app?.mostRecentlySuccessfulDeployment}
          />
        </div>
      </Box>

      <Box top={0} right={0} pr={2} pt={2} position="absolute">
        <SplitButton
          label="select sync strategy"
          color="primary"
          loading={isSyncing}
          onClick={handleSync}
          options={syncOptions}
          startIcon={<SyncIcon />}
        />
      </Box>
    </Paper>
  );
});
