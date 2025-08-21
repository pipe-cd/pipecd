import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useEditPiped = (): UseMutationResult<
  void,
  unknown,
  { pipedId: string; name: string; desc: string },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ pipedId, name, desc }) => {
      await pipedsApi.updatePiped({
        pipedId,
        name,
        desc,
      });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["pipeds", "list"]);
    },
  });
};
