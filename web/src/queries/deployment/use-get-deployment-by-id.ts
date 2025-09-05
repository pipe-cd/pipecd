import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as deploymentsApi from "~/api/deployments";
import { Deployment } from "pipecd/web/model/deployment_pb";

export const useGetDeploymentById = (
  { deploymentId }: { deploymentId: string },
  queryOptions?: UseQueryOptions<Deployment.AsObject>
): UseQueryResult<Deployment.AsObject> => {
  return useQuery({
    queryKey: ["deployment", "detail", { deploymentId }],
    queryFn: async () => {
      const { deployment } = await deploymentsApi.getDeployment({
        deploymentId,
      });
      return deployment as Deployment.AsObject;
    },
    ...queryOptions,
  });
};
