import {
  ApplicationLiveStateVersion,
  KubernetesApplicationLiveState,
  KubernetesResourceState,
  ResourceState,
} from "pipecd/web/model/application_live_state_pb";
import { ApplicationKind } from "pipecd/web/model/common_pb";
import {
  ApplicationLiveState,
  ApplicationLiveStateSnapshot,
} from "~/modules/applications-live-state";
import { dummyApplication, dummyApps } from "./dummy-application";
import { dummyPiped } from "./dummy-piped";
import { createRandTimes, randomUUID } from "./utils";

const resourceIds = [
  randomUUID(),
  randomUUID(),
  randomUUID(),
  randomUUID(),
  randomUUID(),
];
const resourceTimes = createRandTimes(5);

export const resourcesList: KubernetesResourceState.AsObject[] = [
  {
    id: resourceIds[0],
    ownerIdsList: [resourceIds[2]],
    parentIdsList: [],
    name: "demo-application-9504e8601a",
    apiVersion: "apps/v1",
    kind: "ReplicaSet",
    namespace: "default",
    healthStatus: KubernetesResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[0].unix(),
    updatedAt: resourceTimes[0].unix(),
  },
  {
    id: resourceIds[1],
    ownerIdsList: [resourceIds[2]],
    parentIdsList: [resourceIds[0]],
    name: "demo-application-9504e8601a-7vrdw",
    apiVersion: "v1",
    kind: "Pod",
    namespace: "default",
    healthStatus: KubernetesResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[1].unix(),
    updatedAt: resourceTimes[1].unix(),
  },
  {
    id: "f55c7891-ba25-44bb-bca4-ffbc16b0089f",
    ownerIdsList: [resourceIds[2]],
    parentIdsList: [resourceIds[0]],
    name: "demo-application-9504e8601a-vlgd5",
    apiVersion: "v1",
    kind: "Pod",
    namespace: "default",
    healthStatus: KubernetesResourceState.HealthStatus.OTHER,
    healthDescription: "",
    createdAt: resourceTimes[2].unix(),
    updatedAt: resourceTimes[2].unix(),
  },
];

export const resourcesApplicationList: ResourceState.AsObject[] = [
  {
    id: resourceIds[0],
    parentIdsList: [],
    name: "demo-application-eb8faab5-81cd-555e-97ab-e0dd1b7a63aa",
    healthStatus: ResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[0].unix(),
    updatedAt: resourceTimes[0].unix(),
    resourceType: "KUBERNETES",
    resourceMetadataMap: [
      ["Kind", "KUBERNETES"],
      ["Namespace", "default"],
      ["API Version", "apps/v1"],
    ],
    deployTarget: "local",
    pluginName: "kubernetes",
  },
  {
    id: resourceIds[1],
    parentIdsList: [resourceIds[0]],
    name: "demo-application-27428ee6-d6de-50c0-a4dc-a17904e416cd",
    healthStatus: ResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[1].unix(),
    updatedAt: resourceTimes[1].unix(),
    resourceType: "KUBERNETES",
    resourceMetadataMap: [
      ["Kind", "KUBERNETES"],
      ["Namespace", "default"],
      ["API Version", "apps/v1"],
    ],
    deployTarget: "local",
    pluginName: "kubernetes",
  },
  {
    id: resourceIds[2],
    parentIdsList: [],
    name: "demo-application-cf4728f1-ccf1-5128-ac12-c9dfa7a7c616",
    healthStatus: ResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[2].unix(),
    updatedAt: resourceTimes[2].unix(),
    resourceType: "KUBERNETES",
    resourceMetadataMap: [
      ["Kind", "KUBERNETES"],
      ["Namespace", "default"],
      ["API Version", "apps/v1"],
    ],
    deployTarget: "kubernetes",
    pluginName: "kubernetes",
  },
  {
    id: resourceIds[3],
    parentIdsList: [resourceIds[2]],
    name: "demo-application-613c9b52-9962-5ba1-9063-9c618f4c1151",
    healthStatus: ResourceState.HealthStatus.HEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[2].unix(),
    updatedAt: resourceTimes[2].unix(),
    resourceType: "KUBERNETES",
    resourceMetadataMap: [
      ["Kind", "KUBERNETES"],
      ["Namespace", "default"],
      ["API Version", "apps/v1"],
    ],
    deployTarget: "kubernetes",
    pluginName: "kubernetes",
  },
  {
    id: resourceIds[4],
    parentIdsList: [resourceIds[2]],
    name: "demo-application-4f999867-5f5d-506c-9bf3-9d37686b8ae7",
    healthStatus: ResourceState.HealthStatus.UNHEALTHY,
    healthDescription: "",
    createdAt: resourceTimes[2].unix(),
    updatedAt: resourceTimes[2].unix(),
    resourceType: "ReplicaSet",
    resourceMetadataMap: [
      ["Kind", "ReplicaSet"],
      ["Namespace", "default"],
      ["API Version", "apps/v1"],
    ],
    deployTarget: "kubernetes",
    pluginName: "kubernetes",
  },
];

export const dummyApplicationLiveState: ApplicationLiveState = {
  applicationId: dummyApplication.id,
  healthStatus: ApplicationLiveStateSnapshot.Status.HEALTHY,
  kind: ApplicationKind.KUBERNETES,
  pipedId: dummyPiped.id,
  version: { index: 1, timestamp: 0 },
  projectId: "project-1",
  cloudrun: { resourcesList: [] },
  ecs: { resourcesList: [] },
  lambda: { resourcesList: [] },
  terraform: {},
  kubernetes: { resourcesList },
  applicationLiveState: { resourcesList: [], healthStatus: 0 },
};

export const dummyLiveStates: Record<ApplicationKind, ApplicationLiveState> = {
  [ApplicationKind.KUBERNETES]: {
    ...dummyApplicationLiveState,
    applicationId: dummyApps[ApplicationKind.KUBERNETES].id,
    kind: ApplicationKind.KUBERNETES,
    kubernetes: { resourcesList },
  },
  [ApplicationKind.TERRAFORM]: {
    ...dummyApplicationLiveState,
    applicationId: dummyApps[ApplicationKind.TERRAFORM].id,
    kind: ApplicationKind.TERRAFORM,
  },
  [ApplicationKind.LAMBDA]: {
    ...dummyApplicationLiveState,
    applicationId: dummyApps[ApplicationKind.LAMBDA].id,
    kind: ApplicationKind.LAMBDA,
  },
  [ApplicationKind.CLOUDRUN]: {
    ...dummyApplicationLiveState,
    applicationId: dummyApps[ApplicationKind.CLOUDRUN].id,
    kind: ApplicationKind.CLOUDRUN,
  },
  [ApplicationKind.ECS]: {
    ...dummyApplicationLiveState,
    applicationId: dummyApps[ApplicationKind.ECS].id,
    kind: ApplicationKind.ECS,
  },
};

function createKubernetesResourceStateFromObject(
  o: KubernetesResourceState.AsObject
): KubernetesResourceState {
  const state = new KubernetesResourceState();
  state.setId(o.id);
  state.setName(o.name);
  state.setApiVersion(o.apiVersion);
  state.setKind(o.kind);
  state.setNamespace(o.namespace);
  state.setHealthStatus(o.healthStatus);
  state.setHealthDescription(o.healthDescription);
  state.setCreatedAt(o.createdAt);
  state.setUpdatedAt(o.updatedAt);
  state.setOwnerIdsList(o.ownerIdsList);
  state.setParentIdsList(o.parentIdsList);
  return state;
}

function createKubernetesApplicationLiveStateFromObject(
  o: KubernetesApplicationLiveState.AsObject
): KubernetesApplicationLiveState {
  const state = new KubernetesApplicationLiveState();
  state.setResourcesList(
    o.resourcesList.map(createKubernetesResourceStateFromObject)
  );
  return state;
}

export function createLiveStateSnapshotFromObject(
  o: ApplicationLiveState
): ApplicationLiveStateSnapshot {
  const snapshot = new ApplicationLiveStateSnapshot();
  snapshot.setApplicationId(o.applicationId);
  snapshot.setHealthStatus(o.healthStatus);
  snapshot.setKind(o.kind);
  snapshot.setPipedId(o.pipedId);
  if (o.version) {
    const version = new ApplicationLiveStateVersion();
    version.setIndex(o.version.index);
    version.setTimestamp(o.version.timestamp);
    snapshot.setVersion(version);
  }
  snapshot.setProjectId(o.projectId);
  if (o.kubernetes) {
    snapshot.setKubernetes(
      createKubernetesApplicationLiveStateFromObject(o.kubernetes)
    );
  }
  return snapshot;
}
