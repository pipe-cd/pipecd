import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as applicationsAPI from "~/api/applications";

export const useEnableApplication = (): UseMutationResult<
  void,
  unknown,
  { applicationId: string },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload) => {
      await applicationsAPI.enableApplication(payload);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["applications"] });
    },
  });
};
