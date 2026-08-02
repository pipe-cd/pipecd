import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";

export const useToggleProjectAvailability = (): UseMutationResult<
  void,
  unknown,
  { enable: boolean }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ enable }) => {
      if (enable) {
        await projectAPI.enableProject();
      } else {
        await projectAPI.disableProject();
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
