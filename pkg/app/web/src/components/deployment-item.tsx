import React, { FC } from "react";
import { makeStyles, Typography } from "@material-ui/core";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  Deployment,
  selectById as selectDeploymentById,
} from "../modules/deployments";
import { StatusIcon } from "./deployment-status-icon";
import {
  Application,
  selectById as selectApplicationById,
} from "../modules/applications";
import dayjs from "dayjs";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
  },
  appName: {
    marginRight: theme.spacing(2),
  },
  head: {
    display: "flex",
    alignItems: "center",
  },
  description: {
    color: theme.palette.text.hint,
  },
  main: {
    flex: 1,
  },
}));

interface Props {
  id: string;
}

export const DeploymentItem: FC<Props> = ({ id }) => {
  const classes = useStyles();
  const deployment = useSelector<AppState, Deployment | undefined>((state) =>
    selectDeploymentById(state.deployments, id)
  );
  const application = useSelector<AppState, Application | undefined>(
    (state) => {
      if (!deployment) {
        return undefined;
      }
      return selectApplicationById(
        state.applications,
        deployment.applicationId
      );
    }
  );

  if (!deployment || !application) {
    return null;
  }

  return (
    <div className={classes.root}>
      <div className={classes.main}>
        <div className={classes.head}>
          <Typography variant="h6" className={classes.appName}>
            {application.name}
          </Typography>
          <StatusIcon status={deployment.status} />
        </div>
        <Typography variant="body1" className={classes.description}>
          {deployment.description}
        </Typography>
      </div>
      <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
    </div>
  );
};
