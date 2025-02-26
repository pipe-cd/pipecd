import { Box, Chip, makeStyles, Typography } from "@material-ui/core";
import dayjs from "dayjs";
import { FC } from "react";
import { DeploymentStatusIcon } from "~/components/deployment-status-icon";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { ellipsis } from "~/styles/text";
import { Deployment } from "~~/model/deployment_pb";

type Props = {
  deployment: Deployment.AsObject;
};

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    // backgroundColor: theme.palette.background.paper,
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
    color: theme.palette.text.hint,
  },
  labelChip: {
    marginLeft: theme.spacing(1),
    marginBottom: theme.spacing(0.25),
  },
}));

enum PipedVersion {
  V0 = "v0",
  V1 = "v1",
}

const NO_DESCRIPTION = "No description.";

const DeploymentItem: FC<Props> = ({ deployment }) => {
  const classes = useStyles();

  const pipedVersion =
    !deployment.platformProvider ||
    deployment?.deployTargetsByPluginMap?.length > 0
      ? PipedVersion.V1
      : PipedVersion.V0;

  return (
    <Box className={classes.root}>
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
    </Box>
  );
};

export default DeploymentItem;
