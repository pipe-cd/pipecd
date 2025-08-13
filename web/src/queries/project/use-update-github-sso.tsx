import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";
import { GitHubSSO } from "./use-get-project";

export const useUpdateGithubSso = (): UseMutationResult<
  void,
  unknown,
  Partial<GitHubSSO> & { clientId: string; clientSecret: string }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      await projectAPI.updateGitHubSSO(params);
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
