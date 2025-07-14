import { Box, Typography } from "@mui/material";
import Skeleton from "@mui/material/Skeleton";
import { FC, memo } from "react";
import { APPLICATION_HEALTH_STATUS_TEXT } from "~/constants/health-status-text";
import { UI_TEXT_NOT_AVAILABLE_TEXT } from "~/constants/ui-text";
import { ApplicationLiveState } from "~/queries/application-live-state/use-get-application-state-by-id";
import { ApplicationHealthStatusIcon } from "../health-status-icon";

export interface AppLiveStateProps {
  liveState?: ApplicationLiveState;
  liveStateLoading: boolean;
}

export const AppLiveState: FC<AppLiveStateProps> = memo(function AppLiveState({
  liveState,
  liveStateLoading = false,
}) {
  if (liveStateLoading && liveState === undefined) {
    return <Skeleton height={32} width={100} />;
  }

  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
      }}
    >
      {liveState ? (
        <ApplicationHealthStatusIcon health={liveState.healthStatus} />
      ) : null}
      <Typography
        variant="h6"
        sx={{
          ml: 0.5,
        }}
      >
        {liveState
          ? APPLICATION_HEALTH_STATUS_TEXT[liveState.healthStatus]
          : UI_TEXT_NOT_AVAILABLE_TEXT}
      </Typography>
    </Box>
  );
});
