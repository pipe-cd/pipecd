import {
  Box,
  Button,
  CircularProgress,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { UI_TEXT_REFRESH } from "../constants/ui-text";
import { AppState } from "../modules";
import {
  ApplicationLiveState,
  clearError,
  selectById as selectLiveStateById,
} from "../modules/applications-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";

interface Props {
  applicationId: string;
}

const ERROR_MESSAGE = "It was unable to fetch the latest state of application.";

const useStyles = makeStyles(() => ({
  container: {
    flex: 1,
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
}));

export const ApplicationStateView: FC<Props> = memo(
  function ApplicationStateView({ applicationId }) {
    const classes = useStyles();
    const dispatch = useDispatch();
    const [hasError, liveState] = useSelector<
      AppState,
      [boolean, ApplicationLiveState | undefined]
    >((state) => [
      state.applicationLiveState.hasError,
      selectLiveStateById(state.applicationLiveState, applicationId),
    ]);

    if (hasError) {
      return (
        <Box className={classes.container} flexDirection="column">
          <Typography variant="body1">{ERROR_MESSAGE}</Typography>
          <Button
            color="primary"
            onClick={() => {
              dispatch(clearError());
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
