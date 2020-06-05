import React, { FC } from "react";
import { makeStyles, Paper, Typography, Box, Link } from "@material-ui/core";
import { DeploymentStatus, Commit } from "pipe/pkg/app/web/model/deployment_pb";
import { StatusIcon } from "./deployment-status-icon";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { LabeledText } from "./labeled-text";

const useStyles = makeStyles(theme => ({
  container: {
    padding: theme.spacing(2)
  },
  textMargin: {
    marginLeft: theme.spacing(1)
  },
  env: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1)
  }
}));

interface Props {
  name: string;
  env: string;
  status: DeploymentStatus;
  pipedId: string;
  description: string;
  commit: Commit.AsObject;
}

export const DeploymentDetail: FC<Props> = ({
  name,
  pipedId,
  env,
  status,
  description,
  commit
}) => {
  const classes = useStyles();
  return (
    <Paper square elevation={1} className={classes.container}>
      <Box display="flex">
        <Box flex={1}>
          <Box alignItems="center" display="flex">
            <StatusIcon status={status} />
            <Typography className={classes.textMargin} variant="h6">
              {DEPLOYMENT_STATE_TEXT[status]}
            </Typography>
            <Typography className={classes.textMargin} variant="h6">
              {name}
            </Typography>
            <Typography variant="subtitle2" className={classes.env}>
              {env}
            </Typography>
          </Box>
          <LabeledText label="piped" text={pipedId} />
          <LabeledText label="Description" text={description} />
        </Box>
        <Box flex={2}>
          <Box display="flex">
            <Typography variant="subtitle2" color="textSecondary">
              COMMIT
            </Typography>
            <Box display="flex">
              <Typography variant="body2" className={classes.textMargin}>
                {commit.message}
              </Typography>
              <span className={classes.textMargin}>
                (<Link variant="body2">{`${commit.hash}`}</Link>)
              </span>
            </Box>
          </Box>
          {/* TODO: Display createAt time as text */}
          <LabeledText label="CREATED AT" text={`${commit.createdAt}`} />
          <LabeledText label="TRIGGERED BY" text={commit.author} />
        </Box>
      </Box>
    </Paper>
  );
};
