import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as pipedsApi from "~/api/piped";
import { Piped } from "pipecd/web/model/piped_pb";
import { useMemo } from "react";

export const useGetPipedById = (
  options: {
    pipedId: string;
    withStatus?: boolean;
  },
  queryOption: UseQueryOptions<Piped.AsObject[]> = {}
): Omit<UseQueryResult<Piped.AsObject[]>, "data"> & {
  data?: Piped.AsObject;
} => {
  const query = useQuery<Piped.AsObject[]>({
    queryKey: ["pipeds", "list", { withStatus: options.withStatus }],
    queryFn: async () => {
      const { pipedsList } = await pipedsApi.getPipeds({
        withStatus: true,
      });

      return pipedsList;
    },
    ...queryOption,
  });
  const piped = useMemo(() => {
    return query.data?.find((p) => p.id === options.pipedId);
  }, [options.pipedId, query.data]);

  return {
    ...query,
    data: piped,
  };
};
