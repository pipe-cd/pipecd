import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import { getApplicationLiveState } from "~/api/applications";
import { ApplicationLiveStateSnapshot } from "~~/model/application_live_state_pb";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshot.AsObject
>;

export const useGetApplicationStateById = (
  applicationId: string,
  options?: UseQueryOptions<ApplicationLiveState>
): UseQueryResult<ApplicationLiveState> => {
  return useQuery({
    queryKey: ["application", "applicationLiveState", { applicationId }],
    queryFn: async () => {
      const { snapshot } = await getApplicationLiveState({
        applicationId,
      });
      return (snapshot as unknown) as ApplicationLiveState;
    },
    ...options,
  });
};
