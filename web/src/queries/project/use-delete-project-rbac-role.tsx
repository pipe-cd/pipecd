import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";

export const useDeleteProjectRBACRole = (): UseMutationResult<
  void,
  unknown,
  { name: string }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      await projectAPI.deleteRBACRole(params);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
