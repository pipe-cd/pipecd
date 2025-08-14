import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useAddNewPipedKey = (): UseMutationResult<
  string,
  unknown,
  {
    pipedId: string;
  },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ pipedId }) => {
      const { key } = await pipedsApi.recreatePipedKey({ id: pipedId });
      return key;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["pipeds", "list"]);
    },
  });
};
