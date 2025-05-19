import { Box, Chip, ListItem, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";

import dayjs from "dayjs";
import { FC, memo } from "react";
import { Link as RouterLink } from "react-router-dom";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";
import { useAppSelector } from "~/hooks/redux";
import {
  Deployment,
  selectById as selectDeploymentById,
} from "~/modules/deployments";
import { ellipsis } from "~/styles/text";
import { DeploymentStatusIcon } from "~/components/deployment-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    backgroundColor: theme.palette.background.paper,
  },
  info: {
    marginLeft: theme.spacing(1),
  },
  statusText: {
    marginLeft: theme.spacing(1),
    lineHeight: "1.5rem",
    // Fix width to prevent misalignment of application name.
    width: "100px",
  },
  description: {
    ...ellipsis,
    color: theme.palette.text.secondary, // TODO check this color from hint #aaa to secondary #666
  },
  labelChip: {
    marginLeft: theme.spacing(1),
    marginBottom: theme.spacing(0.25),
  },
}));

export interface DeploymentItemProps {
  id: string;
}

enum PipedVersion {
  V0 = "v0",
  V1 = "v1",
}

const NO_DESCRIPTION = "No description.";

export const DeploymentItem: FC<DeploymentItemProps> = memo(
  function DeploymentItem({ id }) {
    const classes = useStyles();
    const deployment = useAppSelector<Deployment.AsObject | undefined>(
      (state) => selectDeploymentById(state.deployments, id)
    );

    if (!deployment) {
      return null;
    }

    const pipedVersion =
      !deployment.platformProvider ||
      deployment?.deployTargetsByPluginMap?.length > 0
        ? PipedVersion.V1
        : PipedVersion.V0;

    return (
      <ListItem
        className={classes.root}
        button
        dense
        divider
        component={RouterLink}
        to={`${PAGE_PATH_DEPLOYMENTS}/${deployment.id}`}
      >
        <Box display="flex" alignItems="center">
          <DeploymentStatusIcon status={deployment.status} />
          <Typography
            variant="subtitle2"
            className={classes.statusText}
            component="span"
          >
            {DEPLOYMENT_STATE_TEXT[deployment.status]}
          </Typography>
        </Box>
        <Box
          display="flex"
          flexDirection="column"
          flex={1}
          pl={2}
          overflow="hidden"
        >
          <Box display="flex" alignItems="baseline">
            <Typography variant="h6" component="span">
              {deployment.applicationName}
            </Typography>
            <Typography
              variant="body2"
              color="textSecondary"
              className={classes.info}
            >
              {pipedVersion === PipedVersion.V0 &&
                APPLICATION_KIND_TEXT[deployment.kind]}
              {pipedVersion === PipedVersion.V1 && "APPLICATION"}
              {deployment?.labelsMap.map(([key, value], i) => (
                <Chip
                  label={key + ": " + value}
                  className={classes.labelChip}
                  key={i}
                />
              ))}
            </Typography>
          </Box>
          <Typography variant="body1" className={classes.description}>
            {deployment.summary || NO_DESCRIPTION}
          </Typography>
        </Box>
        <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
      </ListItem>
    );
  }
);
