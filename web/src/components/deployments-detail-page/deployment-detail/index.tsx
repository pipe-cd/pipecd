import {
  Box,
  Chip,
  CircularProgress,
  Link,
  Paper,
  Typography,
} from "@mui/material";
import CancelIcon from "@mui/icons-material/Cancel";
import OpenInNewIcon from "@mui/icons-material/OpenInNew";
import dayjs from "dayjs";
import { FC, memo, useMemo, useState } from "react";
import { Link as RouterLink } from "react-router-dom";
import { CopyIconButton } from "~/components/copy-icon-button";
import { DeploymentStatusIcon } from "~/components/deployment-status-icon";
import { DetailTableRow } from "~/components/detail-table-row";
import { SplitButton } from "~/components/split-button";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { isDeploymentRunning } from "~/utils/is-deployment-running";
import { Deployment } from "pipecd/web/model/deployment_pb";
import { useGetPipedById } from "~/queries/pipeds/use-get-piped-by-id";
import { useCancelDeployment } from "~/queries/deployment/use-cancel-deployment";
import { useCommand } from "~/contexts/command-context";
import { Command, CommandStatus } from "~~/model/command_pb";

enum PIPED_VERSION {
  V0 = "v0",
  V1 = "v1",
}

export interface DeploymentDetailProps {
  deploymentId: string;
  deployment?: Deployment.AsObject;
}

const CANCEL_OPTIONS = [
  "Cancel",
  "Cancel with Rollback",
  "Cancel without Rollback",
];

export const DeploymentDetail: FC<DeploymentDetailProps> = memo(
  function DeploymentDetail({ deploymentId, deployment }) {
    const [commandId, setCommandId] = useState<string>();
    const { fetchedCommands, commandIds } = useCommand();

    const {
      mutate: cancelDeployment,
      isLoading: isCancelInitLoading,
    } = useCancelDeployment();

    const handleCancelDeployment = async (payload: {
      deploymentId: string;
      forceRollback: boolean;
      forceNoRollback: boolean;
    }): Promise<void> => {
      cancelDeployment(payload, {
        onSuccess: (commandId) => {
          setCommandId(commandId);
        },
      });
    };

    const { data: piped } = useGetPipedById(
      { withStatus: true, pipedId: deployment?.pipedId ?? "" },
      { enabled: !!deployment?.pipedId }
    );

    const isCanceling = useMemo(() => {
      const isCancelCommandRunning = (deploymentId: string): boolean => {
        const deploymentCommand = Object.values(fetchedCommands).find(
          (item) =>
            item.deploymentId === deploymentId &&
            item.type === Command.Type.CANCEL_DEPLOYMENT &&
            item.status === CommandStatus.COMMAND_NOT_HANDLED_YET
        );
        if (deploymentCommand) {
          return true;
        }
        return false;
      };

      const isCancelCommandInit = commandIds?.has(commandId ?? "");

      return (
        isCancelInitLoading ||
        isCancelCommandInit ||
        isCancelCommandRunning(deploymentId)
      );
    }, [
      commandId,
      commandIds,
      deploymentId,
      fetchedCommands,
      isCancelInitLoading,
    ]);

    const pipedVersion = useMemo(() => {
      if (deployment?.deployTargetsByPluginMap?.length) return PIPED_VERSION.V1;

      return PIPED_VERSION.V0;
    }, [deployment?.deployTargetsByPluginMap?.length]);

    if (!deployment || !piped) {
      return (
        <Box
          sx={{
            flex: 1,
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
          }}
        >
          <CircularProgress />
        </Box>
      );
    }

    return (
      <Paper
        square
        elevation={1}
        sx={{
          padding: 2,
          position: "relative",
        }}
      >
        <Box
          sx={{
            display: "flex",
            flexDirection: "column",
          }}
        >
          <Box
            sx={{
              flex: 1,
            }}
          >
            <Box
              sx={{
                display: "flex",
                alignItems: "center",
              }}
            >
              <DeploymentStatusIcon status={deployment.status} />
              <Typography
                variant="h6"
                sx={{
                  ml: 1,
                }}
              >
                {DEPLOYMENT_STATE_TEXT[deployment.status]}
              </Typography>
              <Typography
                variant="body1"
                sx={{
                  color: "text.secondary",
                  marginLeft: 1,
                }}
              >
                {dayjs(deployment.createdAt * 1000).fromNow()}
              </Typography>
              {deployment.labelsMap.map(([key, value], i) => (
                <Chip
                  label={key + ": " + value}
                  sx={{
                    marginLeft: 1,
                    marginBottom: 0.25,
                  }}
                  variant="outlined"
                  key={i}
                />
              ))}
            </Box>
            <Typography
              variant="body2"
              color="textSecondary"
              sx={{
                pt: 1,
                pb: 1,
              }}
            >
              {deployment.statusReason}
            </Typography>
          </Box>
          <Box
            sx={{
              display: "flex",
            }}
          >
            <Box
              sx={{
                flex: 1,
              }}
            >
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
                  {pipedVersion === PIPED_VERSION.V0 && (
                    <DetailTableRow
                      label="Platform Provider"
                      value={deployment.platformProvider}
                    />
                  )}
                  {pipedVersion === PIPED_VERSION.V1 && (
                    <DetailTableRow
                      label="Deploy Targets"
                      value={deployment?.deployTargetsByPluginMap
                        ?.map(([pluginName, { deployTargetsList }]) =>
                          deployTargetsList.map(
                            (deployTarget) => `${deployTarget} - ${pluginName}`
                          )
                        )
                        .join(", ")}
                    />
                  )}
                  <DetailTableRow label="Summary" value={deployment.summary} />
                </tbody>
              </table>
            </Box>
            <Box
              sx={{
                flex: 1,
              }}
            >
              <table>
                <tbody>
                  {deployment.trigger?.commit && (
                    <DetailTableRow
                      label="Commit"
                      value={
                        <Box
                          sx={{
                            display: "flex",
                          }}
                        >
                          <Typography variant="body2">
                            {deployment.trigger.commit.message}
                            <Typography
                              component={"span"}
                              sx={{
                                ml: 1,
                              }}
                            >
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
                                <OpenInNewIcon
                                  sx={{
                                    fontSize: 16,
                                    verticalAlign: "text-bottom",
                                    marginLeft: 0.5,
                                  }}
                                />
                              </Link>
                              )
                            </Typography>
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
            </Box>
            {isDeploymentRunning(deployment.status) && (
              <Box
                sx={(theme) => ({
                  color: "error.main",
                  position: "absolute",
                  top: theme.spacing(2),
                  right: theme.spacing(2),
                })}
              >
                <SplitButton
                  options={CANCEL_OPTIONS}
                  label="select merge strategy"
                  onClick={(index) => {
                    handleCancelDeployment({
                      deploymentId,
                      forceRollback: index === 1,
                      forceNoRollback: index === 2,
                    });
                  }}
                  startIcon={<CancelIcon />}
                  loading={isCanceling}
                  disabled={isCanceling}
                />
              </Box>
            )}
          </Box>
        </Box>
      </Paper>
    );
  }
);
