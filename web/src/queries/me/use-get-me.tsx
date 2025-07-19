import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { getMe } from "~/api/me";
import { GetMeResponse } from "~~/api_client/service_pb";

export const useGetMe = (
  option: UseQueryOptions<GetMeResponse.AsObject> = {}
): UseQueryResult<GetMeResponse.AsObject> => {
  return useQuery({
    queryKey: ["me"],
    queryFn: () => getMe(),
    ...option,
  });
};
