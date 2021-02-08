import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import {
  ApplicationLiveStateVersion,
  KubernetesApplicationLiveState,
  KubernetesResourceState,
} from "pipe/pkg/app/web/model/application_live_state_pb";
import {
  ApplicationLiveState,
  ApplicationLiveStateSnapshot,
} from "../modules/applications-live-state";
import { dummyApplication } from "./dummy-application";
import { dummyEnv } from "./dummy-environment";
import faker from "faker";
import { dummyPiped } from "./dummy-piped";

const resourceIds = [
  faker.random.uuid(),
  faker.random.uuid(),
  faker.random.uuid(),
];

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
    createdAt: 1592472088,
    updatedAt: 1592472088,
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
    createdAt: 1592472088,
    updatedAt: 1592472088,
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
    createdAt: 1592472088,
    updatedAt: 1592472088,
  },
];

export const dummyApplicationLiveState: ApplicationLiveState = {
  applicationId: dummyApplication.id,
  healthStatus: ApplicationLiveStateSnapshot.Status.HEALTHY,
  envId: dummyEnv.id,
  kind: ApplicationKind.KUBERNETES,
  pipedId: dummyPiped.id,
  version: { index: 1, timestamp: 0 },
  projectId: "project-1",
  cloudrun: {},
  lambda: {},
  terraform: {},
  kubernetes: { resourcesList },
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
  snapshot.setEnvId(o.envId);
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
