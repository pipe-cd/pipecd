import {
  useMutation,
  UseMutationResult,
  useQueryClient,
} from "@tanstack/react-query";
import * as APIKeysAPI from "~/api/api-keys";
import { GenerateAPIKeyResponse } from "~~/api_client/service_pb";
import { APIKey } from "pipecd/web/model/apikey_pb";

export const useGenerateApiKey = (): UseMutationResult<
  GenerateAPIKeyResponse.AsObject,
  unknown,
  { name: string; role: APIKey.Role },
  unknown
> => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (values: { name: string; role: APIKey.Role }) => {
      return APIKeysAPI.generateAPIKey(values);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["api-keys"] });
    },
  });
};
