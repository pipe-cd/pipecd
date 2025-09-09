import {
  Cached,
  CheckCircle,
  Error,
  IndeterminateCheckBox,
  Stop,
  Block,
} from "@mui/icons-material";
import { FC } from "react";
import { StageStatus } from "~~/model/deployment_pb";

export interface StageStatusIconProps {
  status: StageStatus;
}

export const StageStatusIcon: FC<StageStatusIconProps> = ({ status }) => {
  switch (status) {
    case StageStatus.STAGE_SUCCESS:
      return <CheckCircle sx={{ color: "success.main" }} />;
    case StageStatus.STAGE_FAILURE:
      return <Error sx={{ color: "error.main" }} />;
    case StageStatus.STAGE_CANCELLED:
      return <Stop sx={{ color: "error.main" }} />;
    case StageStatus.STAGE_NOT_STARTED_YET:
      return <IndeterminateCheckBox sx={{ color: "grey.500" }} />;
    case StageStatus.STAGE_RUNNING:
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
    case StageStatus.STAGE_SKIPPED:
      return <Block sx={{ color: "grey.500" }} />;
    case StageStatus.STAGE_EXITED:
      return <CheckCircle sx={{ color: "success.main" }} />;
  }
};
