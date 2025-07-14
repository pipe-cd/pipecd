import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { Application } from "pipecd/web/model/application_pb";
import * as applicationsAPI from "~/api/applications";

export const useGetApplicationDetail = (
  applicationId: Application.AsObject["id"],
  queryOption: UseQueryOptions<Application.AsObject> = {}
): UseQueryResult<Application.AsObject> => {
  return useQuery({
    queryKey: ["applications", "detail", { applicationId }],
    queryFn: async () => {
      const { application } = await applicationsAPI.getApplication({
        applicationId,
      });
      return application as Application.AsObject;
    },
    ...queryOption,
  });
};
