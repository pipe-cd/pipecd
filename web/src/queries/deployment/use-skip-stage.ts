import {
  useQueryClient,
  useMutation,
  UseMutationResult,
} from "@tanstack/react-query";
import * as deploymentsApi from "~/api/deployments";
import { useCommand } from "~/contexts/command-context";

export const useSkipStage = (): UseMutationResult<
  string,
  unknown,
  { deploymentId: string; stageId: string },
  unknown
> => {
  const queryClient = useQueryClient();
  const { addCommand } = useCommand();

  return useMutation({
    mutationFn: async (payload) => {
      const { commandId } = await deploymentsApi.skipStage(payload);
      return commandId;
    },
    onSuccess: (commandId) => {
      addCommand(commandId);
      queryClient.invalidateQueries(["deployment"]);
    },
  });
};
