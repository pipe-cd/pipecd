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
import clsx from "clsx";
import dayjs from "dayjs";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link as RouterLink } from "react-router-dom";
import { APPLICATION_KIND_TEXT } from "../../constants/application-kind";
import { PAGE_PATH_DEPLOYMENTS } from "../../constants/path";
import { UI_TEXT_REFRESH } from "../../constants/ui-text";
import { AppState } from "../../modules";
import {
  Application,
  ApplicationDeploymentReference,
  fetchApplication,
  selectById as selectApplicationById,
  syncApplication,
} from "../../modules/applications";
import { SyncStrategy } from "../../modules/deployments";
import { selectEnvById } from "../../modules/environments";
import { selectPipedById } from "../../modules/pipeds";
import { AppLiveState } from "../app-live-state";
import { AppSyncStatus } from "../app-sync-status";
import { DetailTableRow } from "../detail-table-row";
import { SplitButton } from "../split-button";
import { SyncStateReason } from "../sync-state-reason";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    display: "flex",
    zIndex: theme.zIndex.appBar,
    position: "relative",
    flexDirection: "column",
  },
  disabled: {
    opacity: 0.6,
  },
  content: {
    flex: 1,
  },
  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  appSyncState: {
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

export interface ApplicationDetailProps {
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

export const ApplicationDetail: FC<ApplicationDetailProps> = memo(
  function ApplicationDetail({ applicationId }) {
    const classes = useStyles();
    const dispatch = useDispatch();

    const [app, fetchApplicationError] = useSelector<
      AppState,
      [Application.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectApplicationById(state.applications, applicationId),
      state.applications.fetchApplicationError,
    ]);

    const env = useSelector(selectEnvById(app?.envId));
    const piped = useSelector(selectPipedById(app?.pipedId));
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
      <Paper
        square
        elevation={1}
        className={clsx(classes.root, {
          [classes.disabled]: app?.disabled,
        })}
      >
        <Box flex={1}>
          <Box display="flex" alignItems="baseline">
            <Typography variant="h5">
              {app ? app.name : <Skeleton width={100} />}
            </Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {env ? env.name : <Skeleton width={100} />}
            </Typography>
          </Box>

          {app ? (
            <>
              <Box display="flex" alignItems="center">
                <AppSyncStatus
                  syncState={app.syncState}
                  deploying={app.deploying}
                  size="large"
                  className={classes.appSyncState}
                />
                <AppLiveState applicationId={applicationId} />
              </Box>

              {app.syncState && (
                <SyncStateReason
                  summary={app.syncState.shortReason}
                  detail={app.syncState.reason}
                />
              )}
            </>
          ) : (
            <Skeleton height={32} width={200} />
          )}
        </Box>

        <Box mt={1} display="flex">
          <div className={classes.content}>
            {app && piped ? (
              <table>
                <tbody>
                  <DetailTableRow
                    label="Kind"
                    value={APPLICATION_KIND_TEXT[app.kind]}
                  />
                  <DetailTableRow label="Piped" value={piped.name} />
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
            disabled={isSyncing || Boolean(app?.disabled)}
            onClick={handleSync}
            options={syncOptions}
            startIcon={<SyncIcon />}
          />
        </Box>
      </Paper>
    );
  }
);
