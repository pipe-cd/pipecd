import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useAddPiped = (): UseMutationResult<
  { id: string; key: string },
  unknown,
  { name: string; desc: string },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (data) => {
      const res = await pipedsApi.registerPiped({
        desc: data.desc,
        name: data.name,
      });
      return res;
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["pipeds", "list"]);
    },
  });
};
