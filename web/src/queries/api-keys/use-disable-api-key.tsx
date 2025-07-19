import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as APIKeysAPI from "~/api/api-keys";

export const useDisableApiKey = (): UseMutationResult<
  void,
  unknown,
  { id: string },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id }) => {
      await APIKeysAPI.disableAPIKey({ id });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["api-keys"] });
    },
  });
};
