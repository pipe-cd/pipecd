import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as applicationsAPI from "~/api/applications";
import { ApplicationInfo } from "pipecd/web/model/common_pb";

export const useGetUnregisteredApplications = (
  queryOption: UseQueryOptions<ApplicationInfo.AsObject[]> = {}
): UseQueryResult<ApplicationInfo.AsObject[]> => {
  return useQuery({
    queryKey: ["applications", "unregistered"],
    queryFn: async () => {
      const {
        applicationsList,
      } = await applicationsAPI.getUnregisteredApplications();
      return applicationsList as ApplicationInfo.AsObject[];
    },
    ...queryOption,
  });
};
