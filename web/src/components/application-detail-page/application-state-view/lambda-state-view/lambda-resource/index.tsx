import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { LambdaResourceState } from "~/modules/applications-live-state";
import { LambdaResourceHealthStatusIcon } from "./health-status-icon";

export interface LambdaResourceProps {
  resource: LambdaResourceState.AsObject;
  onClick: (resource: LambdaResourceState.AsObject) => void;
}

export const LambdaResource: FC<LambdaResourceProps> = memo(
  function LambdaResource({ resource, onClick }) {
    return (
      <Paper
        square
        sx={{
          display: "inline-flex",
          flexDirection: "column",
          padding: 2,
          width: 300,
          cursor: "pointer",
        }}
        onClick={() => onClick(resource)}
      >
        <Typography variant="caption">{resource.kind}</Typography>
        <Box sx={{ display: "flex" }}>
          <LambdaResourceHealthStatusIcon health={resource.healthStatus} />
          <Typography
            variant="subtitle2"
            sx={{
              ml: 0.5,
            }}
          >
            {resource.name}
          </Typography>
        </Box>
      </Paper>
    );
  }
);
