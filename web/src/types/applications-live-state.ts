import {
  ApplicationLiveStateSnapshot,
  KubernetesResourceState,
} from "pipecd/web/model/application_live_state_pb";

export type ApplicationLiveState = Required<
  ApplicationLiveStateSnapshot.AsObject
>;

export const HealthStatus = KubernetesResourceState.HealthStatus;
export type HealthStatus = KubernetesResourceState.HealthStatus;

export {
  ApplicationLiveStateSnapshot,
  KubernetesResourceState,
  CloudRunResourceState,
  ECSResourceState,
  LambdaResourceState,
} from "pipecd/web/model/application_live_state_pb";
