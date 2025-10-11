import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";

export const useGetBreakingChanges = (
  options: { projectId: string },
  queryOption: UseQueryOptions<string> = {}
): UseQueryResult<string> => {
  return useQuery({
    queryKey: ["pipeds", "breakingChanges", options],
    queryFn: async () => {
      const { notes } = await pipedsApi.listBreakingChanges({
        projectId: options.projectId,
      });
      return notes;
    },
    retry: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
    refetchOnWindowFocus: false,
    staleTime: 120000, // 2 minutes
    cacheTime: 300000, // 5 minutes
    ...queryOption,
  });
};
