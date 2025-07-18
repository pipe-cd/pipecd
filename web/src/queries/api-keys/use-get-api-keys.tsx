import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { APIKey } from "pipecd/web/model/apikey_pb";
import * as APIKeysAPI from "~/api/api-keys";

export const useGetApiKeys = (
  options: { enabled: boolean },
  queryOption: UseQueryOptions<APIKey.AsObject[]> = {}
): UseQueryResult<APIKey.AsObject[]> => {
  return useQuery({
    queryKey: ["api-keys", "list"],
    queryFn: async () => {
      const res = await APIKeysAPI.getAPIKeys({ options });
      return res.keysList;
    },
    ...queryOption,
  });
};
