import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as applicationsAPI from "~/api/applications";

export const useDisableApplication = (): UseMutationResult<
  void,
  unknown,
  { applicationId: string },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload) => {
      await applicationsAPI.disableApplication(payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["applications"] });
    },
  });
};
