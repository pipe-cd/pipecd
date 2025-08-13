import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";

export const useToggleAvailabilityStaticAdmin = (): UseMutationResult<
  void,
  unknown,
  { isEnabled: boolean }
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params) => {
      if (params.isEnabled) {
        await projectAPI.enableStaticAdmin();
      } else {
        await projectAPI.disableStaticAdmin();
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries(["project", "detail"]);
    },
  });
};
