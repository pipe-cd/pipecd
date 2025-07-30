import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";

export const useDeleteUserGroup = (): UseMutationResult<
  void,
  unknown,
  { ssoGroup: string }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      await projectAPI.deleteUserGroup(params);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
