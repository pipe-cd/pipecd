import { Cached, CheckCircle, Error, Info, Warning } from "@mui/icons-material";
import { FC } from "react";
import { ApplicationSyncStatus } from "~/modules/applications";

export interface SyncStatusIconProps {
  status: ApplicationSyncStatus;
}

export const SyncStatusIcon: FC<SyncStatusIconProps> = ({ status }) => {
  switch (status) {
    case ApplicationSyncStatus.UNKNOWN:
      return <Info sx={{ color: "grey.500" }} />;
    case ApplicationSyncStatus.SYNCED:
      return <CheckCircle sx={{ color: "success.main" }} />;
    case ApplicationSyncStatus.DEPLOYING:
      return (
        <Cached
          sx={{
            color: "info.main",
            animation: "spin 3s linear infinite",
            "@keyframes spin": {
              "0%": { transform: "rotate(360deg)" },
              "100%": { transform: "rotate(0deg)" },
            },
          }}
        />
      );
    case ApplicationSyncStatus.OUT_OF_SYNC:
      return <Error sx={{ color: "error.main" }} />;
    case ApplicationSyncStatus.INVALID_CONFIG:
      return <Warning sx={{ color: "warning.light" }} />;
  }
};
