import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import {
  Application,
  ApplicationSyncStatus,
} from "pipecd/web/model/application_pb";
import * as applicationsAPI from "~/api/applications";
import { ApplicationKind } from "pipecd/web/model/common_pb";

export interface ApplicationsFilterOptions {
  activeStatus?: string;
  kind?: string;
  syncStatus?: string;
  name?: string;
  pipedId?: string;
  // Suppose to be like ["key-1:value-1"]
  // sindresorhus/query-string doesn't support multidimensional arrays, that's why the format is a bit tricky.
  labels?: Array<string>;
}

export const useGetApplications = (
  filterValues: ApplicationsFilterOptions = {},
  queryOption: UseQueryOptions<Application.AsObject[]> = {}
): UseQueryResult<Application.AsObject[]> => {
  return useQuery({
    queryKey: ["applications", "list", filterValues],
    queryFn: async () => {
      const labels = new Array<[string, string]>();
      if (filterValues.labels) {
        for (const label of filterValues.labels) {
          const pair = label.split(":");
          if (pair.length === 2) labels.push([pair[0], pair[1]]);
        }
      }
      const req = {
        options: {
          envIdsList: [],
          kindsList: filterValues.kind
            ? [parseInt(filterValues.kind, 10) as ApplicationKind]
            : [],
          name: filterValues.name ?? "",
          pipedId: filterValues.pipedId ?? "",
          syncStatusesList: filterValues.syncStatus
            ? [parseInt(filterValues.syncStatus, 10) as ApplicationSyncStatus]
            : [],
          enabled: filterValues.activeStatus
            ? { value: filterValues.activeStatus === "enabled" }
            : undefined,
          labelsMap: labels,
        },
      };
      const { applicationsList } = await applicationsAPI.getApplications(req);
      return applicationsList.filter(
        (app) => app.deleted === false
      ) as Application.AsObject[];
    },
    ...queryOption,
  });
};
