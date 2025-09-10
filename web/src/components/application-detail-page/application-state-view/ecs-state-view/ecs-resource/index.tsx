import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { ECSResourceState } from "~/types/applications-live-state";
import { ECSResourceHealthStatusIcon } from "./health-status-icon";

export interface ECSResourceProps {
  resource: ECSResourceState.AsObject;
  onClick: (resource: ECSResourceState.AsObject) => void;
}

export const ECSResource: FC<ECSResourceProps> = memo(function ECSResource({
  resource,
  onClick,
}) {
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
        <ECSResourceHealthStatusIcon health={resource.healthStatus} />
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
});
