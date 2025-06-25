import { Box } from "@mui/material";
import { FC, memo, useEffect, useMemo, useState } from "react";
import { useParams } from "react-router-dom";
import { DeploymentDetail } from "./deployment-detail";
import { LogViewer } from "./log-viewer";
import { Pipeline } from "./pipeline";
import { Deployment } from "pipecd/web/model/deployment_pb";
import { useGetDeploymentById } from "~/queries/deployment/use-get-deployment-by-id";
import { useGetStageLogs } from "~/queries/stage-logs/use-get-stage-logs";
import { isDeploymentRunning } from "~/utils/is-deployment-running";
import { findDefaultActiveStageInDeployment } from "~/utils/find-default-active-stage-in-deployment";

const DEPLOYMENT_FETCH_INTERVAL = 4000;
const LOG_FETCH_INTERVAL = 2000;

export type ActiveStageInfo = {
  stageId: string;
  deploymentId: string;
  name: string;
} | null;

export const DeploymentDetailPage: FC = memo(function DeploymentDetailPage() {
  const [activeStageInfo, setActiveStageInfo] = useState<ActiveStageInfo>(null);
  const { deploymentId } = useParams<{ deploymentId: string }>();

  const { data: deployment } = useGetDeploymentById(
    { deploymentId: deploymentId ?? "" },
    {
      enabled: !!deploymentId,
      refetchInterval: (data) => {
        return isDeploymentRunning(data?.status)
          ? DEPLOYMENT_FETCH_INTERVAL
          : false;
      },
    }
  );

  const { data: stageLogs } = useGetStageLogs(
    {
      deployment: deployment ?? ({} as Deployment.AsObject),
      offsetIndex: 0,
      retriedCount: 0,
      stageId: activeStageInfo?.stageId ?? "",
    },
    {
      refetchInterval: () => {
        return !!activeStageInfo && isDeploymentRunning(deployment?.status)
          ? LOG_FETCH_INTERVAL
          : false;
      },
      enabled: !!deployment && !!activeStageInfo?.stageId,
      placeholderData: {
        stageId: activeStageInfo?.stageId ?? "",
        deploymentId: deployment?.id ?? "",
        logBlocks: [],
      },
    }
  );

  useEffect(() => {
    const defaultActiveStage = findDefaultActiveStageInDeployment(deployment);
    if (defaultActiveStage && deployment) {
      setActiveStageInfo({
        deploymentId: deployment.id ?? "",
        stageId: defaultActiveStage.id,
        name: defaultActiveStage.name,
      });
    }
  }, [deployment]);

  // NOTE: Clear active stage when leave detail page
  useEffect(
    () => () => {
      setActiveStageInfo(null);
    },
    []
  );

  const activeStage = useMemo(() => {
    return (
      deployment?.stagesList.find((s) => s.id === activeStageInfo?.stageId) ??
      null
    );
  }, [deployment, activeStageInfo]);

  return (
    <Box
      sx={{
        display: "flex",
        flexDirection: "column",
        alignItems: "stretch",
        flex: 1,
        overflow: "auto",
      }}
    >
      <Box
        sx={{
          flex: 1,
        }}
      >
        <DeploymentDetail
          deploymentId={deploymentId ?? ""}
          deployment={deployment}
        />
        <Pipeline
          deployment={deployment}
          activeStageInfo={activeStageInfo}
          changeActiveStage={setActiveStageInfo}
        />
      </Box>
      <LogViewer
        stageLog={stageLogs}
        changeActiveStage={setActiveStageInfo}
        activeStage={activeStage}
      />
    </Box>
  );
});
