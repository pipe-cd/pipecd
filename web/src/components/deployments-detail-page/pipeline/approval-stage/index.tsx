import { Box, Paper, Typography } from "@mui/material";
import WaitIcon from "@mui/icons-material/PauseCircleOutline";
import { FC, memo } from "react";

export interface ApprovalStageProps {
  id: string;
  name: string;
  active: boolean;
  onClick: (stageId: string, stageName: string) => void;
}

export const ApprovalStage: FC<ApprovalStageProps> = memo(
  function ApprovalStage({ id, name, onClick, active }) {
    function handleOnClick(): void {
      onClick(id, name);
    }

    return (
      <Paper
        square
        onClick={handleOnClick}
        sx={(theme) => ({
          flex: 1,
          display: "inline-flex",
          cursor: "pointer",
          padding: theme.spacing(2),
          "&:hover": {
            backgroundColor: theme.palette.action.hover,
          },
          backgroundColor: active
            ? // NOTE: 12%
              theme.palette.primary.main + "1e"
            : undefined,
        })}
      >
        <Box
          sx={{
            display: "flex",
            justifyContent: "flex-start",
            alignItems: "center",
          }}
        >
          <WaitIcon sx={{ color: "warning.main" }} />
          <Typography
            variant="subtitle2"
            ml={1}
            sx={{ fontFamily: "fontFamilyMono" }}
          >
            {name}
          </Typography>
        </Box>
      </Paper>
    );
  }
);
