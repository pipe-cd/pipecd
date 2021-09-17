import {
  Box,
  Link,
  Button,
  CircularProgress,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { FC, memo, useEffect } from "react";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import {
  Application,
  ApplicationKind,
  selectById as selectAppById,
} from "~/modules/applications";
import {
  ApplicationLiveState,
  fetchApplicationStateById,
  selectById as selectLiveStateById,
  selectHasError,
} from "~/modules/applications-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";

const FETCH_INTERVAL = 4000;

export interface ApplicationStateViewProps {
  applicationId: string;
}

const ERROR_MESSAGE = "It was unable to fetch the latest state of application.";
const COMING_SOON_MESSAGE = "Live state for this kind of application is not yet available.";
const FEATURE_STATUS_INTRO = "PipeCD feature status"
const DISABLED_APPLICATION_MESSAGE =
  "This application is currently disabled. You can enable it from the application list page.";

const useStyles = makeStyles(() => ({
  container: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
}));

export const ApplicationStateView: FC<ApplicationStateViewProps> = memo(
  function ApplicationStateView({ applicationId }) {
    const classes = useStyles();
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

    useEffect(() => {
      if (app?.kind === ApplicationKind.KUBERNETES) {
        dispatch(fetchApplicationStateById(app.id));
      }
    }, [app, dispatch]);

    useInterval(
      () => {
        // Only fetch kubernetes application.
        if (app?.kind === ApplicationKind.KUBERNETES) {
          dispatch(fetchApplicationStateById(app.id));
        }
      },
      // Only fetch kubernetes application.
      app?.kind === ApplicationKind.KUBERNETES && hasError === false ? FETCH_INTERVAL : null
    );

    if (app?.disabled) {
      return (
        <Box
          display="flex"
          justifyContent="center"
          alignItems="center"
          flex={1}
        >
          <Typography variant="h6" component="span">
            {DISABLED_APPLICATION_MESSAGE}
          </Typography>
        </Box>
      );
    }

    if (hasError) {
      return (
        <Box className={classes.container} flexDirection="column">
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
          {app?.kind === ApplicationKind.KUBERNETES ? (
            <div className={classes.container}>
              <CircularProgress />
            </div>
          ) : (
            <Box className={classes.container} flexDirection="column">
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

    switch (liveState.kind) {
      case ApplicationKind.KUBERNETES: {
        const resources = liveState.kubernetes?.resourcesList || [];
        return <KubernetesStateView resources={resources} />;
      }
      default:
    }

    // NOTE: other resource types are not implemented.
    return null;
  }
);
