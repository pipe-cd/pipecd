import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";

export const useUpdateStaticAdmin = (): UseMutationResult<
  void,
  unknown,
  { username?: string; password?: string }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      await projectAPI.updateStaticAdmin(params);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
