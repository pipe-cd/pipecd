import React, { FC, memo } from "react";
import { Box, makeStyles, Typography } from "@material-ui/core";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  ApplicationLiveState,
  selectById,
  selectLoadingById,
} from "../modules/applications-live-state";
import Skeleton from "@material-ui/lab/Skeleton";
import { ApplicationHealthStatusIcon } from "./health-status-icon";
import { APPLICATION_HEALTH_STATUS_TEXT } from "../constants/health-status-text";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "../constants/ui-text";

const useStyles = makeStyles((theme) => ({
  liveStateText: {
    marginLeft: theme.spacing(0.5),
  },
}));

interface Props {
  applicationId: string;
}

export const AppLiveState: FC<Props> = memo(function AppLiveState({
  applicationId,
}) {
  const classes = useStyles();
  const [liveState, liveStateLoading] = useSelector<
    AppState,
    [ApplicationLiveState | undefined, boolean]
  >((state) => [
    selectById(state.applicationLiveState, applicationId),
    selectLoadingById(state.applicationLiveState, applicationId),
  ]);

  if (liveStateLoading) {
    return <Skeleton height={32} width={100} />;
  }

  return (
    <Box display="flex" alignItems="center">
      {liveState ? (
        <ApplicationHealthStatusIcon health={liveState.healthStatus} />
      ) : null}
      <Typography variant="h6" className={classes.liveStateText}>
        {liveState
          ? APPLICATION_HEALTH_STATUS_TEXT[liveState.healthStatus]
          : UI_TEXT_NOT_AVAILABLE_TEXT}
      </Typography>
    </Box>
  );
});
