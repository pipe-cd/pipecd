import { Box, Link, Button, CircularProgress, Typography } from "@mui/material";
import { FC, memo, useEffect, useMemo } from "react";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import {
  Application,
  ApplicationKind,
  PIPED_VERSION,
  selectById as selectAppById,
} from "~/modules/applications";
import {
  ApplicationLiveState,
  fetchApplicationStateById,
  selectById as selectLiveStateById,
  selectHasError,
} from "~/modules/applications-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";
import { CloudRunStateView } from "./cloudrun-state-view";
import { ECSStateView } from "./ecs-state-view";
import { LambdaStateView } from "./lambda-state-view";
import { checkPipedAppVersion } from "~/utils/common";
import { LiveStateView } from "./live-state-view";

const isDisplayLiveState = (app: Application.AsObject | undefined): boolean => {
  const result = checkPipedAppVersion(app);
  if (result[PIPED_VERSION.V1]) return true;

  return (
    app?.kind === ApplicationKind.KUBERNETES ||
    app?.kind === ApplicationKind.CLOUDRUN ||
    app?.kind === ApplicationKind.ECS ||
    app?.kind === ApplicationKind.LAMBDA
  );
};

const FETCH_INTERVAL = 4000;

export interface ApplicationStateViewProps {
  applicationId: string;
}

const ERROR_MESSAGE = "It was unable to fetch the latest state of application.";
const COMING_SOON_MESSAGE =
  "Live state for this kind of application is not yet available.";
const FEATURE_STATUS_INTRO = "PipeCD feature status";
const DISABLED_APPLICATION_MESSAGE =
  "This application is currently disabled. You can enable it from the application list page.";

export const ApplicationStateView: FC<ApplicationStateViewProps> = memo(
  function ApplicationStateView({ applicationId }) {
    const dispatch = useAppDispatch();
    const [hasError, liveState, app] = useAppSelector<
      [
        boolean,
        ApplicationLiveState | undefined,
        Application.AsObject | undefined
      ]
    >((state) => [
      selectHasError(state.applicationLiveState, applicationId),
      selectLiveStateById(state.applicationLiveState, applicationId),
      selectAppById(state.applications, applicationId),
    ]);

    const pipedVersion = useMemo(() => checkPipedAppVersion(app), [app]);

    useEffect(() => {
      if (app && isDisplayLiveState(app)) {
        dispatch(fetchApplicationStateById(app.id));
      }
    }, [app, dispatch]);

    useInterval(
      () => {
        // Fetch only supported kind applications.
        if (app && isDisplayLiveState(app)) {
          dispatch(fetchApplicationStateById(app.id));
        }
      },
      // Fetch only supported kind applications.
      isDisplayLiveState(app) && hasError === false ? FETCH_INTERVAL : null
    );

    if (app?.disabled) {
      return (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flex: 1,
          }}
        >
          <Typography variant="h6" component="span">
            {DISABLED_APPLICATION_MESSAGE}
          </Typography>
        </Box>
      );
    }

    if (hasError) {
      return (
        <Box
          sx={{
            flexDirection: "column",
            flex: 1,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <Typography variant="body1">{ERROR_MESSAGE}</Typography>
          <Button
            color="primary"
            onClick={() => {
              dispatch(fetchApplicationStateById(applicationId));
            }}
          >
            {UI_TEXT_REFRESH}
          </Button>
        </Box>
      );
    }

    if (!liveState) {
      return (
        <>
          {isDisplayLiveState(app) ? (
            <Box
              sx={{
                flex: 1,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <CircularProgress />
            </Box>
          ) : (
            <Box
              sx={{
                flexDirection: "column",
                flex: 1,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
              }}
            >
              <Typography variant="body1">{COMING_SOON_MESSAGE}</Typography>
              <Link
                href="https://pipecd.dev/docs/feature-status/"
                target="_blank"
                rel="noreferrer"
              >
                {FEATURE_STATUS_INTRO}
              </Link>
            </Box>
          )}
        </>
      );
    }

    if (pipedVersion[PIPED_VERSION.V1]) {
      const resources = liveState.applicationLiveState?.resourcesList || [];
      return <LiveStateView resources={resources} />;
    }

    switch (liveState.kind) {
      case ApplicationKind.KUBERNETES: {
        const resources = liveState.kubernetes?.resourcesList || [];
        return <KubernetesStateView resources={resources} />;
      }
      case ApplicationKind.CLOUDRUN: {
        const resources = liveState.cloudrun?.resourcesList || [];
        return <CloudRunStateView resources={resources} />;
      }
      case ApplicationKind.ECS: {
        const resources = liveState.ecs?.resourcesList || [];
        return <ECSStateView resources={resources} />;
      }
      case ApplicationKind.LAMBDA: {
        const resources = liveState.lambda?.resourcesList || [];
        return <LambdaStateView resources={resources} />;
      }
      default:
    }

    // NOTE: other resource types are not implemented.
    return null;
  }
);
