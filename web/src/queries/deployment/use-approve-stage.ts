import {
  useQueryClient,
  useMutation,
  UseMutationResult,
} from "@tanstack/react-query";
import * as deploymentsApi from "~/api/deployments";
import { useCommand } from "~/contexts/command-context";

export const useApproveStage = (): UseMutationResult<
  string,
  unknown,
  { deploymentId: string; stageId: string },
  unknown
> => {
  const queryClient = useQueryClient();
  const { addCommand } = useCommand();

  return useMutation({
    mutationFn: async (payload) => {
      const { commandId } = await deploymentsApi.approveStage(payload);
      return commandId;
    },
    onSuccess: (commandId) => {
      addCommand(commandId);
      queryClient.invalidateQueries(["deployment"]);
    },
  });
};
