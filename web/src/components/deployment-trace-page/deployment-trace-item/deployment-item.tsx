import { Box, Chip, Typography } from "@mui/material";
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

enum PipedVersion {
  V0 = "v0",
  V1 = "v1",
}

const NO_DESCRIPTION = "No description.";

const DeploymentItem: FC<Props> = ({ deployment }) => {
  const pipedVersion =
    !deployment.platformProvider ||
    deployment?.deployTargetsByPluginMap?.length > 0
      ? PipedVersion.V1
      : PipedVersion.V0;

  return (
    <Box
      sx={(theme) => ({
        flex: 1,
        padding: theme.spacing(2),
        display: "flex",
        alignItems: "center",
        overflow: "hidden",
        columnGap: theme.spacing(2),
        [theme.breakpoints.down("md")]: {
          flexDirection: "column",
          alignItems: "flex-start",
        },
      })}
    >
      <Box
        sx={{
          display: "flex",
          alignItems: "center",
        }}
      >
        <DeploymentStatusIcon status={deployment.status} />
        <Typography
          variant="subtitle2"
          sx={{
            marginLeft: 1,
            lineHeight: "1.5rem",
            width: "100px",
          }}
          component="span"
        >
          {DEPLOYMENT_STATE_TEXT[deployment.status]}
        </Typography>
      </Box>
      <Box
        sx={{
          flex: 1,
          overflow: "hidden",
          maxWidth: "100%",
        }}
      >
        <Box
          sx={(theme) => ({
            display: "flex",
            alignItems: "baseline",
            flexWrap: "wrap",
            columnGap: theme.spacing(1),
            [theme.breakpoints.down("md")]: {
              flexDirection: "column",
            },
          })}
        >
          <Typography variant="h6" component="span">
            {deployment.applicationName}
          </Typography>
          <Typography variant="body2" color="textSecondary">
            {pipedVersion === PipedVersion.V0 &&
              APPLICATION_KIND_TEXT[deployment.kind]}
            {pipedVersion === PipedVersion.V1 && "APPLICATION"}
            {deployment?.labelsMap.map(([key, value], i) => (
              <Chip
                label={key + ": " + value}
                sx={{ ml: 1, mb: 0.25 }}
                key={i}
              />
            ))}
          </Typography>
        </Box>
        <Typography
          variant="body1"
          sx={{ ...ellipsis, color: "text.secondary" }}
        >
          {deployment.summary || NO_DESCRIPTION}
        </Typography>
      </Box>
      <div>{dayjs(deployment.createdAt * 1000).fromNow()}</div>
    </Box>
  );
};

export default DeploymentItem;
