import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";
import { Piped } from "pipecd/web/model/piped_pb";

export const useGetPipeds = (
  options: {
    withStatus: boolean;
  },
  queryOption: UseQueryOptions<Piped.AsObject[]> = {}
): UseQueryResult<Piped.AsObject[]> => {
  return useQuery({
    queryKey: ["pipeds", "list", options],
    queryFn: async () => {
      const { pipedsList } = await pipedsApi.getPipeds({
        withStatus: options.withStatus,
      });

      return pipedsList;
    },
    ...queryOption,
  });
};
