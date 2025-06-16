import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useDeleteOldPipedKey = (): UseMutationResult<
  void,
  unknown,
  {
    pipedId: string;
  },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ pipedId }) => {
      await pipedsApi.deleteOldPipedKey({ pipedId });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["pipeds", "list"]);
    },
  });
};
