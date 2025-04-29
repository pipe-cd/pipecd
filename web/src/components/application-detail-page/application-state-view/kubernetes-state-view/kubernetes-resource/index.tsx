import { Paper, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC, memo } from "react";
import { KubernetesResourceState } from "~/modules/applications-live-state";
import { KubernetesResourceHealthStatusIcon } from "./health-status-icon";

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

export interface KubernetesResourceProps {
  resource: KubernetesResourceState.AsObject;
  onClick: (resource: KubernetesResourceState.AsObject) => void;
}

export const KubernetesResource: FC<KubernetesResourceProps> = memo(
  function KubernetesResource({ resource, onClick }) {
    const classes = useStyles();
    return (
      <Paper square className={classes.root} onClick={() => onClick(resource)}>
        <Typography variant="caption">{resource.kind}</Typography>
        <div className={classes.nameLine}>
          <KubernetesResourceHealthStatusIcon health={resource.healthStatus} />
          <Typography variant="subtitle2" className={classes.name}>
            {resource.name}
          </Typography>
        </div>
      </Paper>
    );
  }
);
