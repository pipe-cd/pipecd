import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";
import { ProjectRBACPolicy } from "pipecd/web/model/project_pb";

export const useAddProjectRBACRole = (): UseMutationResult<
  void,
  unknown,
  { name: string; policies: ProjectRBACPolicy[] }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      await projectAPI.addRBACRole(params);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
