import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useGetReleasedVersions = (
  queryOption: UseQueryOptions<string[]> = {}
): UseQueryResult<string[]> => {
  return useQuery({
    queryKey: ["pipeds", "releasedVersions"],
    queryFn: async () => {
      const { versionsList } = await pipedsApi.listReleasedVersions();
      return versionsList;
    },
    refetchOnMount: false,
    refetchOnReconnect: false,
    refetchOnWindowFocus: false,
    staleTime: 120000, // 2 minutes
    cacheTime: 300000, // 5 minutes
    ...queryOption,
  });
};
