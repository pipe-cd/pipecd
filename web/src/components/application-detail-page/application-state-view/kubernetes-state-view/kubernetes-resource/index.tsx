import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { KubernetesResourceState } from "~/modules/applications-live-state";
import { KubernetesResourceHealthStatusIcon } from "./health-status-icon";

export interface KubernetesResourceProps {
  resource: KubernetesResourceState.AsObject;
  onClick: (resource: KubernetesResourceState.AsObject) => void;
}

export const KubernetesResource: FC<KubernetesResourceProps> = memo(
  function KubernetesResource({ resource, onClick }) {
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
          <KubernetesResourceHealthStatusIcon health={resource.healthStatus} />
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
