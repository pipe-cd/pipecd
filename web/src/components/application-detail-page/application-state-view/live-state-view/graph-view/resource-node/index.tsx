import { Box, Paper, Typography } from "@mui/material";
import { FC, memo } from "react";
import { ResourceState } from "~~/model/application_live_state_pb";
import { HealthStatusIcon } from "./health-status-icon";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";

type Props = {
  resource: ResourceState.AsObject;
  onClick: (resource: ResourceState.AsObject) => void;
};

export const ResourceNode: FC<Props> = memo(function ResourceNode({
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
      <Typography variant="caption">
        {findMetadataByKey(resource.resourceMetadataMap, "Kind")}
      </Typography>
      <Box sx={{ display: "flex" }}>
        <HealthStatusIcon health={resource.healthStatus} />
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
