import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as deploymentsApi from "~/api/deployments";
import { useCommand } from "~/contexts/command-context";

export const useCancelDeployment = (): UseMutationResult<
  string,
  unknown,
  {
    deploymentId: string;
    forceRollback: boolean;
    forceNoRollback: boolean;
  },
  unknown
> => {
  const queryClient = useQueryClient();
  const { addCommand } = useCommand();

  return useMutation({
    mutationFn: async ({
      deploymentId,
      forceRollback,
      forceNoRollback,
    }: {
      deploymentId: string;
      forceRollback: boolean;
      forceNoRollback: boolean;
    }) => {
      const { commandId } = await deploymentsApi.cancelDeployment({
        deploymentId,
        forceRollback,
        forceNoRollback,
      });
      return commandId;
    },
    onSuccess: (commandId) => {
      addCommand(commandId);
      queryClient.invalidateQueries(["deployment"]);
    },
  });
};
// await thunkAPI.dispatch(fetchCommand(commandId));
