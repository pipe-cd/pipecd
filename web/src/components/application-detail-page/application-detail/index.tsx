import {
  Box,
  Button,
  Link,
  Paper,
  Typography,
  Chip,
  Skeleton,
} from "@mui/material";
import SyncIcon from "@mui/icons-material/Cached";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import { SerializedError } from "@reduxjs/toolkit";
import dayjs from "dayjs";
import { FC, Fragment, memo, useMemo, useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import ReactMarkdown from "react-markdown";
import { AppSyncStatus } from "~/components/app-sync-status";
import { DetailTableRow } from "~/components/detail-table-row";
import { SplitButton } from "~/components/split-button";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { unwrapResult, useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  Application,
  ApplicationDeploymentReference,
  ApplicationSyncStatus,
  fetchApplication,
  selectById as selectApplicationById,
  syncApplication,
} from "~/modules/applications";
import { SyncStrategy } from "~/modules/deployments";
import { selectPipedById } from "~/modules/pipeds";
import { AppLiveState } from "./app-live-state";
import { OutOfSyncReason, InvalidConfigReason } from "./sync-state-reason";
import { ArtifactVersion } from "~~/model/common_pb";
import { CopyIconButton } from "~/components/copy-icon-button";

export interface ApplicationDetailProps {
  applicationId: string;
}

const useIsSyncingApplication = (
  applicationId: string | undefined
): boolean => {
  return useAppSelector<boolean>((state) => {
    if (!applicationId) {
      return false;
    }

    return state.applications.syncing[applicationId];
  });
};

const ERROR_MESSAGE = "It was unable to fetch the application.";

const ArtifactVersions: FC<{
  deployment: ApplicationDeploymentReference.AsObject;
}> = ({ deployment }) => {
  const defaultDisplayLimit = 4;

  const [showMore, setShowMore] = useState(false);

  if (deployment.versionsList.length === 0) {
    return <span>{deployment.version}</span>;
  }

  const buildLinkableArtifactVersion = (
    v: ArtifactVersion.AsObject
  ): React.ReactChild => {
    return v.name === "" ? (
      <>
        <span>{v.version}</span>
        <br />
      </>
    ) : (
      <>
        <Link
          href={v.url.includes("://") ? v.url : `//${v.url}`}
          target="_blank"
          rel="noreferrer"
        >
          {v.name}:{v.version}
          <OpenInNewIcon
            sx={{
              fontSize: 16,
              verticalAlign: "text-bottom",
              marginLeft: 0.5,
            }}
          />
        </Link>
        <br />
      </>
    );
  };

  if (deployment.versionsList.length <= defaultDisplayLimit) {
    return (
      <>
        {deployment.versionsList.map((v) => (
          <Fragment key={`${v.name}:${v.version}`}>
            {buildLinkableArtifactVersion(v)}
          </Fragment>
        ))}
      </>
    );
  }

  return !showMore ? (
    <>
      {deployment.versionsList.map((v, idx) => {
        if (idx >= defaultDisplayLimit) return;
        return (
          <Fragment key={`${v.name}:${v.version}`}>
            {buildLinkableArtifactVersion(v)}
          </Fragment>
        );
      })}
      <Box
        component={"span"}
        sx={{ cursor: "pointer" }}
        onClick={() => setShowMore(!showMore)}
      >
        show more...
      </Box>
    </>
  ) : (
    <>
      {deployment.versionsList.map((v) => (
        <Fragment key={`${v.name}:${v.version}`}>
          {buildLinkableArtifactVersion(v)}
        </Fragment>
      ))}
      <Box
        component={"span"}
        sx={{ cursor: "pointer" }}
        onClick={() => setShowMore(!showMore)}
      >
        show less...
      </Box>
    </>
  );
};

const MostRecentlySuccessfulDeployment: FC<{
  deployment?: ApplicationDeploymentReference.AsObject;
}> = ({ deployment }) => {
  if (!deployment) {
    return <Skeleton height={105} width={500} />;
  }

  const date = dayjs(deployment.startedAt * 1000);

  return (
    <>
      <Box display="flex" alignItems="baseline">
        <Typography variant="subtitle1">Latest Deployment</Typography>
        <Typography variant="body2" ml={1}>
          <Link
            component={RouterLink}
            to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.deploymentId}`}
          >
            more details
          </Link>
        </Typography>
      </Box>
      <Box pl={2}>
        <table>
          <tbody>
            <DetailTableRow
              label="Deployed At"
              value={<span title={date.format()}>{date.fromNow()}</span>}
            />
            <DetailTableRow
              label="Artifact Versions"
              value={<ArtifactVersions deployment={deployment} />}
            />
            <DetailTableRow label="Summary" value={deployment.summary} />
          </tbody>
        </table>
      </Box>
    </>
  );
};

const syncOptions = ["Sync", "Quick Sync", "Pipeline Sync"];
const syncStrategyByIndex: SyncStrategy[] = [
  SyncStrategy.AUTO,
  SyncStrategy.QUICK_SYNC,
  SyncStrategy.PIPELINE,
];

enum PIPED_VERSION {
  V0 = "v0",
  V1 = "v1",
}

export const ApplicationDetail: FC<ApplicationDetailProps> = memo(
  function ApplicationDetail({ applicationId }) {
    const dispatch = useAppDispatch();

    const [app, fetchApplicationError] = useAppSelector<
      [Application.AsObject | undefined, SerializedError | null]
    >((state) => [
      selectApplicationById(state.applications, applicationId),
      state.applications.fetchApplicationError,
    ]);

    const piped = useAppSelector(selectPipedById(app?.pipedId));
    const isSyncing = useIsSyncingApplication(app?.id);
    const description = app?.description.replace(/\\\n/g, "  \n") || "";

    const handleSync = (index: number): void => {
      if (app) {
        dispatch(
          syncApplication({
            applicationId: app.id,
            syncStrategy: syncStrategyByIndex[index],
          })
        )
          .then(unwrapResult)
          .catch(() => undefined);
      }
    };

    const pipedVersion = useMemo(() => {
      if (!app?.platformProvider) return PIPED_VERSION.V1;
      if (app?.deployTargetsByPluginMap?.length) return PIPED_VERSION.V1;

      return PIPED_VERSION.V0;
    }, [app?.deployTargetsByPluginMap.length, app?.platformProvider]);

    if (fetchApplicationError) {
      return (
        <Paper
          square
          elevation={1}
          sx={{
            padding: 2,
            display: "flex",
            zIndex: "appBar",
            position: "relative",
            flexDirection: "column",
          }}
        >
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
        sx={{
          padding: 2,
          display: "flex",
          zIndex: "appBar",
          position: "relative",
          flexDirection: "column",
          opacity: app?.disabled ? 0.6 : 1,
        }}
      >
        <Box flex={1}>
          <Box display="flex" alignItems="baseline">
            <Typography variant="h5">
              {app ? app.name : <Skeleton width={100} />}
            </Typography>
            {app?.labelsMap.map(([key, value], i) => (
              <Chip
                label={key + ": " + value}
                variant="outlined"
                sx={{ ml: 1 }}
                key={i}
              />
            ))}
          </Box>

          {app ? (
            <>
              <Box display="flex" alignItems="center" gap={1}>
                <AppSyncStatus
                  syncState={app.syncState}
                  deploying={app.deploying}
                  size="large"
                />
                <AppLiveState applicationId={applicationId} />
              </Box>

              {app.syncState &&
                app.syncState.status == ApplicationSyncStatus.OUT_OF_SYNC && (
                  <OutOfSyncReason
                    summary={app.syncState.shortReason}
                    detail={app.syncState.reason}
                  />
                )}

              {app.syncState &&
                app.syncState.status ==
                  ApplicationSyncStatus.INVALID_CONFIG && (
                  <InvalidConfigReason
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
          <Box flex={1}>
            {app && piped ? (
              <table>
                <tbody>
                  <DetailTableRow
                    label="Application ID"
                    value={
                      <>
                        {applicationId}
                        <CopyIconButton
                          name="Application ID"
                          value={applicationId}
                          size="small"
                        />
                      </>
                    }
                  />
                  {pipedVersion === PIPED_VERSION.V0 && (
                    <DetailTableRow
                      label="Kind"
                      value={APPLICATION_KIND_TEXT[app.kind]}
                    />
                  )}
                  <DetailTableRow label="Piped" value={piped.name} />
                  {pipedVersion === PIPED_VERSION.V0 && (
                    <DetailTableRow
                      label="Platform Provider"
                      value={app.platformProvider}
                    />
                  )}
                  {pipedVersion === PIPED_VERSION.V1 && (
                    <DetailTableRow
                      label="Deploy Targets"
                      value={app?.deployTargetsByPluginMap
                        ?.map(([pluginName, { deployTargetsList }]) =>
                          deployTargetsList.map(
                            (deployTarget) => `${deployTarget} - ${pluginName}`
                          )
                        )
                        .join(", ")}
                    />
                  )}

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
                          <OpenInNewIcon
                            sx={{
                              fontSize: 16,
                              verticalAlign: "text-bottom",
                              marginLeft: 0.5,
                            }}
                          />
                        </Link>
                      }
                    />
                  )}
                </tbody>
              </table>
            ) : (
              <Skeleton height={105} width={500} />
            )}
          </Box>

          <Box flex={1}>
            <MostRecentlySuccessfulDeployment
              deployment={app?.mostRecentlySuccessfulDeployment}
            />
          </Box>
        </Box>

        {app && (
          <Box
            borderLeft="2px solid"
            borderColor="divider"
            pl={2}
            fontSize={"body2.fontSize"}
          >
            <ReactMarkdown linkTarget="_blank">{description}</ReactMarkdown>
          </Box>
        )}

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
