import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
} from "@mui/material";
import { FC, memo, useCallback, useState } from "react";
import {
  METADATA_APPROVED_BY,
  METADATA_SKIPPED_BY,
  METADATA_STAGE_DISPLAY_KEY,
} from "~/constants/metadata-keys";
import { ApprovalStage } from "./approval-stage";
import { PipelineStage } from "./pipeline-stage";
import { ManualOperation } from "~~/model/deployment_pb";
import { useApproveStage } from "~/queries/deployment/use-approve-stage";
import { isDeploymentRunning } from "~/utils/is-deployment-running";
import { Deployment, StageStatus } from "~~/model/deployment_pb";
import { Stage } from "~/types/deployment";
import { ActiveStageInfo } from "..";

enum PIPED_VERSION {
  V0 = "v0",
  V1 = "v1",
}

const WAIT_APPROVAL_NAME = "WAIT_APPROVAL";
const STAGE_HEIGHT = 56;
const APPROVED_STAGE_HEIGHT = 66;

const isStartedStage = (stage: Stage): boolean => {
  return stage.status !== StageStatus.STAGE_NOT_STARTED_YET;
};

/**
 * ## For piped v0
 * ### Visibility of stages
 * - field `visible` = true
 * ### Order of stages
 * - stages with requiresList.length === 0 will be in the first column
 * - stages with requiresList includes id of previous stages will be in the next columns
 *  */
const createStagesPipedV0 = (allStage: Stage[]): Stage[][] => {
  const visibleStages = allStage.filter((stage) => stage.visible);

  const stages: Stage[][] = [];
  stages[0] = visibleStages.filter((stage) => stage.requiresList.length === 0);

  let index = 0;
  while (stages[index].length > 0) {
    const previousIds = stages[index].map((stage) => stage.id);
    index++;
    stages[index] = visibleStages.filter((stage) =>
      stage.requiresList.some(
        (id) =>
          previousIds.includes(id) &&
          // prevent self-requirement
          stage.id !== id
      )
    );
  }

  return stages;
};
/**
 * ## For piped v1
 * ### Visibility of stages
 * - field `rollback` = false or stage is started
 * ### Order of stages (temporary solution)
 * - stages with requiresList.length === 0 will be in the first column
 * - stages with requiresList includes id of previous stages will be in the next columns
 */
export const createStagesPipedV1 = (allStages: Stage[]): Stage[][] => {
  const visibleStages = allStages.filter(
    (stage) => !stage.rollback || isStartedStage(stage)
  );
  const stages: Stage[][] = [];
  stages[0] = visibleStages.filter((stage) => stage.requiresList.length === 0);

  let index = 0;
  while (stages[index].length > 0) {
    const previousIds = stages[index].map((stage) => stage.id);
    index++;
    stages[index] = visibleStages.filter((stage) =>
      stage.requiresList.some(
        (id) =>
          previousIds.includes(id) &&
          // prevent self-requirement
          stage.id !== id
      )
    );
  }

  return stages;
};

const createStagesForRendering = (
  deployment: Deployment.AsObject | undefined
): Stage[][] => {
  if (!deployment) {
    return [];
  }

  const pipedVersion = deployment.deployTargetsByPluginMap.length
    ? PIPED_VERSION.V1
    : PIPED_VERSION.V0;

  if (pipedVersion === PIPED_VERSION.V0) {
    return createStagesPipedV0(deployment.stagesList);
  }
  if (pipedVersion === PIPED_VERSION.V1) {
    return createStagesPipedV1(deployment.stagesList);
  }

  return [];
};

const LARGE_STAGE_NAMES = ["WAIT_APPROVAL", "K8S_TRAFFIC_ROUTING"];

export interface PipelineProps {
  deployment: Deployment.AsObject | undefined;
  activeStageInfo?: ActiveStageInfo;
  changeActiveStage: (info: ActiveStageInfo) => void;
}

// deprecated. Use findDisplayMetadataText for pipedv1.
const findApprover = (
  metadata: Array<[string, string]>
): string | undefined => {
  const res = metadata.find(([key]) => key === METADATA_APPROVED_BY);

  if (res) {
    return `Approved by: ${res[1]}`;
  }

  return undefined;
};

// deprecated. Use findDisplayMetadataText for pipedv1.
const findSkipper = (metadata: Array<[string, string]>): string | undefined => {
  const res = metadata.find(([key]) => key === METADATA_SKIPPED_BY);

  if (res) {
    return `Skipped by: ${res[1]}`;
  }

  return undefined;
};

const findDisplayMetadataText = (
  metadata: Array<[string, string]>
): string | undefined => {
  const res = metadata.find(([key]) => key === METADATA_STAGE_DISPLAY_KEY);
  if (res) {
    return res[1];
  }
  return undefined;
};

export const Pipeline: FC<PipelineProps> = memo(function Pipeline({
  deployment,
  activeStageInfo,
  changeActiveStage: changeActiveStage,
}) {
  const [approveTargetId, setApproveTargetId] = useState<string | null>(null);
  const { mutate: approveStage } = useApproveStage();

  const isOpenApproveDialog = Boolean(approveTargetId);
  const deploymentId = deployment?.id ?? "";

  const stages = createStagesForRendering(deployment);
  const isRunning = isDeploymentRunning(deployment?.status);

  const handleOnClickStage = useCallback(
    (stageId: string, stageName: string) => {
      changeActiveStage({
        deploymentId,
        stageId,
        name: stageName,
      });
    },
    [changeActiveStage, deploymentId]
  );

  const handleApprove = (): void => {
    if (approveTargetId) {
      approveStage({ deploymentId, stageId: approveTargetId });
      setApproveTargetId(null);
    }
  };

  return (
    <Box
      sx={{
        textAlign: "center",
        overflow: "scroll",

        "&::-webkit-scrollbar": {
          height: "7px",
        },

        "&::-webkit-scrollbar-thumb": {
          borderRadius: 8,
          backgroundColor: "rgba(0,0,0,0.3)",
        },
      }}
    >
      <Box
        sx={{
          display: "inline-flex",
        }}
      >
        {stages.map((stageColumn, columnIndex) => {
          let isPrevStageLarge = false;
          return (
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
              }}
              key={`pipeline-${columnIndex}`}
            >
              {stageColumn.map((stage, stageIndex) => {
                const displayMetadataText = findDisplayMetadataText(
                  stage.metadataMap
                );
                // TODO: remove approver and skipper. they should be included in findDisplayMetadataText for compatibility.
                const approver = findApprover(stage.metadataMap);
                const skipper = findSkipper(stage.metadataMap);
                const isActive = activeStageInfo
                  ? activeStageInfo.deploymentId === deploymentId &&
                    activeStageInfo.stageId === stage.id
                  : false;
                const showLine = columnIndex > 0;
                const showStraightLine = showLine && stageIndex === 0;
                const showCurvedLine = showLine && stageIndex > 0;
                // TODO: remove approver. use displayMetadataText instead.
                const isCurvedLineExtend =
                  showCurvedLine && (Boolean(approver) || isPrevStageLarge);

                const stageComp = (
                  <Box
                    key={stage.id}
                    sx={(theme) => ({
                      display: "flex",
                      padding: theme.spacing(2),
                      ...(showStraightLine && {
                        position: "relative",
                        "&::before": {
                          content: '""',
                          position: "absolute",
                          top: "48%",
                          left: theme.spacing(-2),
                          borderTop: `2px solid ${theme.palette.divider}`,
                          width: theme.spacing(4),
                          height: 1,
                        },
                      }),

                      ...(showCurvedLine && {
                        position: "relative",
                        "&::before": {
                          content: '""',
                          position: "absolute",
                          bottom: "50%",
                          left: 0,
                          borderLeft: `2px solid ${theme.palette.divider}`,
                          borderBottom: `2px solid ${theme.palette.divider}`,
                          width: theme.spacing(2),
                          height: `calc(${
                            isCurvedLineExtend
                              ? APPROVED_STAGE_HEIGHT
                              : STAGE_HEIGHT
                          }px + ${theme.spacing(4)})`,
                        },
                      }),
                    })}
                  >
                    {/* TODO: Remove stageName condition after finishing deployments which are made 
                         while the server does not inject availableOperation */}
                    {(stage.name === WAIT_APPROVAL_NAME ||
                      stage.availableOperation ===
                        ManualOperation.MANUAL_OPERATION_APPROVE) &&
                    stage.status === StageStatus.STAGE_RUNNING ? (
                      <ApprovalStage
                        id={stage.id}
                        name={stage.name}
                        onClick={() => {
                          setApproveTargetId(stage.id);
                        }}
                        active={isActive}
                      />
                    ) : (
                      <PipelineStage
                        id={stage.id}
                        name={stage.name}
                        status={stage.status}
                        metadata={stage.metadataMap}
                        onClick={handleOnClickStage}
                        active={isActive}
                        isDeploymentRunning={isRunning}
                        // TODO: use only displayMetadataText
                        displayMetadataText={
                          displayMetadataText || approver || skipper
                        }
                      />
                    )}
                  </Box>
                );
                isPrevStageLarge = LARGE_STAGE_NAMES.includes(stage.name);
                return stageComp;
              })}
            </Box>
          );
        })}

        <Dialog
          open={isOpenApproveDialog}
          onClose={() => setApproveTargetId(null)}
        >
          <DialogTitle>Approve stage</DialogTitle>
          <DialogContent>
            <DialogContentText>
              {`To continue deploying, click "APPROVE".`}
            </DialogContentText>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setApproveTargetId(null)}>CANCEL</Button>
            <Button color="primary" onClick={handleApprove}>
              APPROVE
            </Button>
          </DialogActions>
        </Dialog>
      </Box>
    </Box>
  );
});
