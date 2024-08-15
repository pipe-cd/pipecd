import { makeStyles, Paper, Typography } from "@material-ui/core";
import { FC, memo } from "react";
import { LambdaResourceState } from "~/modules/applications-live-state";
import { LambdaResourceHealthStatusIcon } from "./health-status-icon";

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

export interface LambdaResourceProps {
  resource: LambdaResourceState.AsObject;
  onClick: (resource: LambdaResourceState.AsObject) => void;
}

export const LambdaResource: FC<LambdaResourceProps> = memo(function LambdaResource({
  resource,
  onClick,
}) {
  const classes = useStyles();
  return (
    <Paper square className={classes.root} onClick={() => onClick(resource)}>
      <Typography variant="caption">{resource.kind}</Typography>
      <div className={classes.nameLine}>
        <LambdaResourceHealthStatusIcon health={resource.healthStatus} />
        <Typography variant="subtitle2" className={classes.name}>
          {resource.name}
        </Typography>
      </div>
    </Paper>
  );
});
