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
    ...queryOption,
  });
};
