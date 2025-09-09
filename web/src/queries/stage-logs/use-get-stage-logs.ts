import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { getStageLog } from "~/api/stage-log";
import { StageLog } from "~/types/stage-log";
import { Deployment, StageStatus } from "~~/model/deployment_pb";
import { StatusCode } from "grpc-web";

export const useGetStageLogs = (
  options: {
    deployment?: Deployment.AsObject;
    offsetIndex: number;
    retriedCount: number;
    stageId: string;
  },
  queryOption: UseQueryOptions<StageLog> = {}
): UseQueryResult<StageLog> => {
  return useQuery({
    queryKey: [
      "stage-logs",
      {
        deploymentId: options.deployment?.id ?? "",
        stageId: options.stageId,
        offsetIndex: options.offsetIndex,
        retriedCount: options.retriedCount,
      },
    ],
    queryFn: async () => {
      const { deployment, offsetIndex, retriedCount, stageId } = options || {};
      const deploymentId = deployment?.id ?? "";

      const initialLogs: StageLog = {
        stageId,
        deploymentId,
        logBlocks: [],
      };

      if (!deployment) {
        throw new Error(`Deployment: ${deploymentId} is not exists in state.`);
      }

      const stage = deployment.stagesList.find((stage) => stage.id === stageId);

      if (!stage) {
        throw new Error(
          `Stage (ID: ${stageId}) is not found in application state.`
        );
      }

      if (stage.status === StageStatus.STAGE_NOT_STARTED_YET) {
        return initialLogs;
      }

      const response = await getStageLog({
        deploymentId,
        offsetIndex,
        retriedCount,
        stageId,
      }).catch((e: { code: number }) => {
        // If status is running and error code is NOT_FOUND, it is maybe first state of deployment log.
        // So we ignore this error and then return initialLogs below code.
        if (
          e.code === StatusCode.NOT_FOUND &&
          stage.status === StageStatus.STAGE_RUNNING
        ) {
          return;
        }

        throw e;
      });

      if (!response) {
        return initialLogs;
      }

      return {
        stageId,
        deploymentId,
        logBlocks: response.blocksList,
      };
    },

    ...queryOption,
  });
};
