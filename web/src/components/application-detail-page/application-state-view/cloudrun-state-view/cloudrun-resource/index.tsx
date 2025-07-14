import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { CloudRunResourceState } from "~/types/applications-live-state";
import { CloudRunResourceHealthStatusIcon } from "./health-status-icon";

export interface CloudRunResourceProps {
  resource: CloudRunResourceState.AsObject;
  onClick: (resource: CloudRunResourceState.AsObject) => void;
}

export const CloudRunResource: FC<CloudRunResourceProps> = memo(
  function CloudRunResource({ resource, onClick }) {
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
          <CloudRunResourceHealthStatusIcon health={resource.healthStatus} />
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
