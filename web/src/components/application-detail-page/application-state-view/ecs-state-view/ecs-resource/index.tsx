import { makeStyles, Paper, Typography } from "@material-ui/core";
import { FC, memo } from "react";
import { ECSResourceState } from "~/modules/applications-live-state";
import { ECSResourceHealthStatusIcon } from "./health-status-icon";

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

export interface ECSResourceProps {
  resource: ECSResourceState.AsObject;
  onClick: (resource: ECSResourceState.AsObject) => void;
}

export const ECSResource: FC<ECSResourceProps> = memo(
  function ECSResource({ resource, onClick }) {
    const classes = useStyles();
    return (
      <Paper square className={classes.root} onClick={() => onClick(resource)}>
        <Typography variant="caption">{resource.kind}</Typography>
        <div className={classes.nameLine}>
          <ECSResourceHealthStatusIcon health={resource.healthStatus} />
          <Typography variant="subtitle2" className={classes.name}>
            {resource.name}
          </Typography>
        </div>
      </Paper>
    );
  }
);
