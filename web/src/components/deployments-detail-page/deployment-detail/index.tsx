import {
  Box,
  Chip,
  CircularProgress,
  Link,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import CancelIcon from "@material-ui/icons/Cancel";
import OpenInNewIcon from "@material-ui/icons/OpenInNew";
import dayjs from "dayjs";
import { FC, memo } from "react";
import { Link as RouterLink } from "react-router-dom";
import { CopyIconButton } from "~/components/copy-icon-button";
import { DeploymentStatusIcon } from "~/components/deployment-status-icon";
import { DetailTableRow } from "~/components/detail-table-row";
import { SplitButton } from "~/components/split-button";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import {
  cancelDeployment,
  Deployment,
  isDeploymentRunning,
  selectById as selectDeploymentById,
  selectDeploymentIsCanceling,
} from "~/modules/deployments";
import { selectPipedById } from "~/modules/pipeds";
import { fetchStageLog } from "~/modules/stage-logs";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    position: "relative",
  },
  textMargin: {
    marginLeft: theme.spacing(1),
  },
  age: {
    color: theme.palette.text.secondary,
    marginLeft: theme.spacing(1),
  },
  content: {
    flex: 1,
  },
  actionButtons: {
    color: theme.palette.error.main,
    position: "absolute",
    top: theme.spacing(2),
    right: theme.spacing(2),
  },
  statusReason: {
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
  linkIcon: {
    fontSize: 16,
    verticalAlign: "text-bottom",
    marginLeft: theme.spacing(0.5),
  },
  labelChip: {
    marginLeft: theme.spacing(1),
    marginBottom: theme.spacing(0.25),
  },
}));

export interface DeploymentDetailProps {
  deploymentId: string;
}

const CANCEL_OPTIONS = [
  "Cancel",
  "Cancel with Rollback",
  "Cancel without Rollback",
];
const LOG_FETCH_INTERVAL = 2000;

export const DeploymentDetail: FC<DeploymentDetailProps> = memo(
  function DeploymentDetail({ deploymentId }) {
    const classes = useStyles();
    const dispatch = useAppDispatch();

    const deployment = useAppSelector<Deployment.AsObject | undefined>(
      (state) => selectDeploymentById(state.deployments, deploymentId)
    );
    const activeStage = useAppSelector((state) => state.activeStage);
    const piped = useAppSelector(selectPipedById(deployment?.pipedId));
    const isCanceling = useAppSelector(
      selectDeploymentIsCanceling(deploymentId)
    );

    useInterval(
      () => {
        if (activeStage) {
          dispatch(
            fetchStageLog({
              deploymentId: activeStage.deploymentId,
              stageId: activeStage.stageId,
              offsetIndex: 0,
              retriedCount: 0,
            })
          );
        }
      },
      activeStage && isDeploymentRunning(deployment?.status)
        ? LOG_FETCH_INTERVAL
        : null
    );

    if (!deployment || !piped) {
      return (
        <Box
          flex={1}
          display="flex"
          alignItems="center"
          justifyContent="center"
        >
          <CircularProgress />
        </Box>
      );
    }

    return (
      <Paper square elevation={1} className={classes.root}>
        <Box display="flex" flexDirection="column">
          <div className={classes.content}>
            <Box display="flex" alignItems="center">
              <DeploymentStatusIcon status={deployment.status} />
              <Typography className={classes.textMargin} variant="h6">
                {DEPLOYMENT_STATE_TEXT[deployment.status]}
              </Typography>
              <Typography variant="body1" className={classes.age}>
                {dayjs(deployment.createdAt * 1000).fromNow()}
              </Typography>
              {deployment.labelsMap.map(([key, value], i) => (
                <Chip
                  label={key + ": " + value}
                  className={classes.labelChip}
                  variant="outlined"
                  key={i}
                />
              ))}
            </Box>
            <Typography
              variant="body2"
              color="textSecondary"
              className={classes.statusReason}
            >
              {deployment.statusReason}
            </Typography>
          </div>
          <Box display="flex">
            <div className={classes.content}>
              <table>
                <tbody>
                  <DetailTableRow
                    label="Deployment ID"
                    value={
                      <>
                        {deploymentId}
                        <CopyIconButton
                          name="Deployment ID"
                          value={deploymentId}
                          size="small"
                        />
                      </>
                    }
                  />
                  <DetailTableRow
                    label="Application"
                    value={
                      <Link
                        variant="body2"
                        component={RouterLink}
                        to={`${PAGE_PATH_APPLICATIONS}/${deployment.applicationId}`}
                      >
                        {deployment.applicationName}
                      </Link>
                    }
                  />
                  <DetailTableRow label="Piped" value={piped.name} />
                  <DetailTableRow
                    label="Platform Provider"
                    value={deployment.platformProvider}
                  />
                  <DetailTableRow label="Summary" value={deployment.summary} />
                </tbody>
              </table>
            </div>
            <div className={classes.content}>
              <table>
                <tbody>
                  {deployment.trigger?.commit && (
                    <DetailTableRow
                      label="Commit"
                      value={
                        <Box display="flex">
                          <Typography variant="body2">
                            {deployment.trigger.commit.message}
                            <span className={classes.textMargin}>
                              (
                              <Link
                                variant="body2"
                                href={deployment.trigger.commit.url}
                                target="_blank"
                                rel="noreferrer"
                              >
                                {`${deployment.trigger.commit.hash.slice(
                                  0,
                                  7
                                )}`}
                                <OpenInNewIcon className={classes.linkIcon} />
                              </Link>
                              )
                            </span>
                          </Typography>
                        </Box>
                      }
                    />
                  )}
                  <DetailTableRow
                    label="Triggered by"
                    value={
                      deployment.trigger?.commander ||
                      deployment.trigger?.commit?.author ||
                      ""
                    }
                  />
                </tbody>
              </table>
            </div>
            {isDeploymentRunning(deployment.status) && (
              <SplitButton
                className={classes.actionButtons}
                options={CANCEL_OPTIONS}
                label="select merge strategy"
                onClick={(index) => {
                  dispatch(
                    cancelDeployment({
                      deploymentId,
                      forceRollback: index === 1,
                      forceNoRollback: index === 2,
                    })
                  );
                }}
                startIcon={<CancelIcon />}
                loading={isCanceling}
                disabled={isCanceling}
              />
            )}
          </Box>
        </Box>
      </Paper>
    );
  }
);
