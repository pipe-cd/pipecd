import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useUpdatePipedDesiredVersion = (): UseMutationResult<
  void,
  unknown,
  { version: string; pipedIds: string[] },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data) => {
      await pipedsApi.updatePipedDesiredVersion({
        version: data.version,
        pipedIdsList: data.pipedIds,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["pipeds", "list"]);
    },
  });
};
