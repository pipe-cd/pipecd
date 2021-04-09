import {
  Box,
  Button,
  CircularProgress,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { UI_TEXT_REFRESH } from "../../constants/ui-text";
import { AppState } from "../../modules";
import {
  Application,
  selectById as selectAppById,
} from "../../modules/applications";
import {
  ApplicationLiveState,
  fetchApplicationStateById,
  selectById as selectLiveStateById,
  selectHasError,
} from "../../modules/applications-live-state";
import { KubernetesStateView } from "../kubernetes-state-view";

export interface ApplicationStateViewProps {
  applicationId: string;
}

const ERROR_MESSAGE = "It was unable to fetch the latest state of application.";
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
    const dispatch = useDispatch();
    const [hasError, liveState, app] = useSelector<
      AppState,
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
        <div className={classes.container}>
          <CircularProgress />
        </div>
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
