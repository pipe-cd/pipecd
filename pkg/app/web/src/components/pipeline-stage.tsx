import React, { FC } from "react";
import { makeStyles, Paper, Typography, Box } from "@material-ui/core";
import { StageStatus } from "pipe/pkg/app/web/model/deployment_pb";
import { StageStatusIcon } from "./stage-status-icon";

const useStyles = makeStyles(theme => ({
  container: {
    display: "inline-flex"
  },
  name: {
    marginLeft: theme.spacing(1)
  }
}));

interface Props {
  name: string;
  status: StageStatus;
}

export const PipelineStage: FC<Props> = ({ name, status }) => {
  const classes = useStyles();
  return (
    <Paper square className={classes.container}>
      <Box alignItems="center" display="flex" justifyContent="center" p={2}>
        <StageStatusIcon status={status} />
        <Typography variant="subtitle2" className={classes.name}>
          <Box fontFamily="Roboto Mono">{name}</Box>
        </Typography>
      </Box>
    </Paper>
  );
};
