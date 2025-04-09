import { makeStyles, Paper, Typography } from "@material-ui/core";
import { FC, memo } from "react";
import { ResourceState } from "~~/model/application_live_state_pb";
import { HealthStatusIcon } from "./health-status-icon";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "inline-flex",
    flexDirection: "column",
    padding: theme.spacing(2),
    width: 300,
    cursor: "pointer",
  },
  nameLine: {
    display: "flex",
  },
  name: {
    marginLeft: theme.spacing(0.5),
  },
}));

type Props = {
  resource: ResourceState.AsObject;
  onClick: (resource: ResourceState.AsObject) => void;
};

export const ResourceNode: FC<Props> = memo(function ResourceNode({
  resource,
  onClick,
}) {
  const classes = useStyles();
  return (
    <Paper square className={classes.root} onClick={() => onClick(resource)}>
      <Typography variant="caption">
        {findMetadataByKey(resource.resourceMetadataMap, "Kind")}
      </Typography>
      <div className={classes.nameLine}>
        <HealthStatusIcon health={resource.healthStatus} />
        <Typography variant="subtitle2" className={classes.name}>
          {resource.name}
        </Typography>
      </div>
    </Paper>
  );
});
