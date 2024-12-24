import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';


export class ApplicationLiveStateSnapshot extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): ApplicationLiveStateSnapshot;

  getPipedId(): string;
  setPipedId(value: string): ApplicationLiveStateSnapshot;

  getProjectId(): string;
  setProjectId(value: string): ApplicationLiveStateSnapshot;

  getKind(): pkg_model_common_pb.ApplicationKind;
  setKind(value: pkg_model_common_pb.ApplicationKind): ApplicationLiveStateSnapshot;

  getHealthStatus(): ApplicationLiveStateSnapshot.Status;
  setHealthStatus(value: ApplicationLiveStateSnapshot.Status): ApplicationLiveStateSnapshot;

  getKubernetes(): KubernetesApplicationLiveState | undefined;
  setKubernetes(value?: KubernetesApplicationLiveState): ApplicationLiveStateSnapshot;
  hasKubernetes(): boolean;
  clearKubernetes(): ApplicationLiveStateSnapshot;

  getTerraform(): TerraformApplicationLiveState | undefined;
  setTerraform(value?: TerraformApplicationLiveState): ApplicationLiveStateSnapshot;
  hasTerraform(): boolean;
  clearTerraform(): ApplicationLiveStateSnapshot;

  getCloudrun(): CloudRunApplicationLiveState | undefined;
  setCloudrun(value?: CloudRunApplicationLiveState): ApplicationLiveStateSnapshot;
  hasCloudrun(): boolean;
  clearCloudrun(): ApplicationLiveStateSnapshot;

  getLambda(): LambdaApplicationLiveState | undefined;
  setLambda(value?: LambdaApplicationLiveState): ApplicationLiveStateSnapshot;
  hasLambda(): boolean;
  clearLambda(): ApplicationLiveStateSnapshot;

  getEcs(): ECSApplicationLiveState | undefined;
  setEcs(value?: ECSApplicationLiveState): ApplicationLiveStateSnapshot;
  hasEcs(): boolean;
  clearEcs(): ApplicationLiveStateSnapshot;

  getVersion(): ApplicationLiveStateVersion | undefined;
  setVersion(value?: ApplicationLiveStateVersion): ApplicationLiveStateSnapshot;
  hasVersion(): boolean;
  clearVersion(): ApplicationLiveStateSnapshot;

  getApplicationLiveState(): ApplicationLiveState | undefined;
  setApplicationLiveState(value?: ApplicationLiveState): ApplicationLiveStateSnapshot;
  hasApplicationLiveState(): boolean;
  clearApplicationLiveState(): ApplicationLiveStateSnapshot;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationLiveStateSnapshot.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationLiveStateSnapshot): ApplicationLiveStateSnapshot.AsObject;
  static serializeBinaryToWriter(message: ApplicationLiveStateSnapshot, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationLiveStateSnapshot;
  static deserializeBinaryFromReader(message: ApplicationLiveStateSnapshot, reader: jspb.BinaryReader): ApplicationLiveStateSnapshot;
}

export namespace ApplicationLiveStateSnapshot {
  export type AsObject = {
    applicationId: string,
    pipedId: string,
    projectId: string,
    kind: pkg_model_common_pb.ApplicationKind,
    healthStatus: ApplicationLiveStateSnapshot.Status,
    kubernetes?: KubernetesApplicationLiveState.AsObject,
    terraform?: TerraformApplicationLiveState.AsObject,
    cloudrun?: CloudRunApplicationLiveState.AsObject,
    lambda?: LambdaApplicationLiveState.AsObject,
    ecs?: ECSApplicationLiveState.AsObject,
    version?: ApplicationLiveStateVersion.AsObject,
    applicationLiveState?: ApplicationLiveState.AsObject,
  }

  export enum Status { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

export class ApplicationLiveStateVersion extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): ApplicationLiveStateVersion;

  getIndex(): number;
  setIndex(value: number): ApplicationLiveStateVersion;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationLiveStateVersion.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationLiveStateVersion): ApplicationLiveStateVersion.AsObject;
  static serializeBinaryToWriter(message: ApplicationLiveStateVersion, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationLiveStateVersion;
  static deserializeBinaryFromReader(message: ApplicationLiveStateVersion, reader: jspb.BinaryReader): ApplicationLiveStateVersion;
}

export namespace ApplicationLiveStateVersion {
  export type AsObject = {
    timestamp: number,
    index: number,
  }
}

export class KubernetesApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<KubernetesResourceState>;
  setResourcesList(value: Array<KubernetesResourceState>): KubernetesApplicationLiveState;
  clearResourcesList(): KubernetesApplicationLiveState;
  addResources(value?: KubernetesResourceState, index?: number): KubernetesResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesApplicationLiveState): KubernetesApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: KubernetesApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesApplicationLiveState;
  static deserializeBinaryFromReader(message: KubernetesApplicationLiveState, reader: jspb.BinaryReader): KubernetesApplicationLiveState;
}

export namespace KubernetesApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<KubernetesResourceState.AsObject>,
  }
}

export class TerraformApplicationLiveState extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): TerraformApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: TerraformApplicationLiveState): TerraformApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: TerraformApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): TerraformApplicationLiveState;
  static deserializeBinaryFromReader(message: TerraformApplicationLiveState, reader: jspb.BinaryReader): TerraformApplicationLiveState;
}

export namespace TerraformApplicationLiveState {
  export type AsObject = {
  }
}

export class CloudRunApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<CloudRunResourceState>;
  setResourcesList(value: Array<CloudRunResourceState>): CloudRunApplicationLiveState;
  clearResourcesList(): CloudRunApplicationLiveState;
  addResources(value?: CloudRunResourceState, index?: number): CloudRunResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloudRunApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: CloudRunApplicationLiveState): CloudRunApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: CloudRunApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloudRunApplicationLiveState;
  static deserializeBinaryFromReader(message: CloudRunApplicationLiveState, reader: jspb.BinaryReader): CloudRunApplicationLiveState;
}

export namespace CloudRunApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<CloudRunResourceState.AsObject>,
  }
}

export class ECSApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<ECSResourceState>;
  setResourcesList(value: Array<ECSResourceState>): ECSApplicationLiveState;
  clearResourcesList(): ECSApplicationLiveState;
  addResources(value?: ECSResourceState, index?: number): ECSResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ECSApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: ECSApplicationLiveState): ECSApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: ECSApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ECSApplicationLiveState;
  static deserializeBinaryFromReader(message: ECSApplicationLiveState, reader: jspb.BinaryReader): ECSApplicationLiveState;
}

export namespace ECSApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<ECSResourceState.AsObject>,
  }
}

export class LambdaApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<LambdaResourceState>;
  setResourcesList(value: Array<LambdaResourceState>): LambdaApplicationLiveState;
  clearResourcesList(): LambdaApplicationLiveState;
  addResources(value?: LambdaResourceState, index?: number): LambdaResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaApplicationLiveState): LambdaApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: LambdaApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaApplicationLiveState;
  static deserializeBinaryFromReader(message: LambdaApplicationLiveState, reader: jspb.BinaryReader): LambdaApplicationLiveState;
}

export namespace LambdaApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<LambdaResourceState.AsObject>,
  }
}

export class ApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<ResourceState>;
  setResourcesList(value: Array<ResourceState>): ApplicationLiveState;
  clearResourcesList(): ApplicationLiveState;
  addResources(value?: ResourceState, index?: number): ResourceState;

  getHealthStatus(): ApplicationLiveState.Status;
  setHealthStatus(value: ApplicationLiveState.Status): ApplicationLiveState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationLiveState): ApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: ApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationLiveState;
  static deserializeBinaryFromReader(message: ApplicationLiveState, reader: jspb.BinaryReader): ApplicationLiveState;
}

export namespace ApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<ResourceState.AsObject>,
    healthStatus: ApplicationLiveState.Status,
  }

  export enum Status { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

export class ResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): ResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): ResourceState;
  clearParentIdsList(): ResourceState;
  addParentIds(value: string, index?: number): ResourceState;

  getName(): string;
  setName(value: string): ResourceState;

  getResourceType(): string;
  setResourceType(value: string): ResourceState;

  getResourceMetadataMap(): jspb.Map<string, string>;
  clearResourceMetadataMap(): ResourceState;

  getHealthStatus(): ResourceState.HealthStatus;
  setHealthStatus(value: ResourceState.HealthStatus): ResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): ResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): ResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): ResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: ResourceState): ResourceState.AsObject;
  static serializeBinaryToWriter(message: ResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ResourceState;
  static deserializeBinaryFromReader(message: ResourceState, reader: jspb.BinaryReader): ResourceState;
}

export namespace ResourceState {
  export type AsObject = {
    id: string,
    parentIdsList: Array<string>,
    name: string,
    resourceType: string,
    resourceMetadataMap: Array<[string, string]>,
    healthStatus: ResourceState.HealthStatus,
    healthDescription: string,
    createdAt: number,
    updatedAt: number,
  }

  export enum HealthStatus { 
    UNKNOWN = 0,
    HEALTHY = 1,
    UNHEALTHY = 2,
  }
}

export class KubernetesResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): KubernetesResourceState;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): KubernetesResourceState;
  clearOwnerIdsList(): KubernetesResourceState;
  addOwnerIds(value: string, index?: number): KubernetesResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): KubernetesResourceState;
  clearParentIdsList(): KubernetesResourceState;
  addParentIds(value: string, index?: number): KubernetesResourceState;

  getName(): string;
  setName(value: string): KubernetesResourceState;

  getApiVersion(): string;
  setApiVersion(value: string): KubernetesResourceState;

  getKind(): string;
  setKind(value: string): KubernetesResourceState;

  getNamespace(): string;
  setNamespace(value: string): KubernetesResourceState;

  getHealthStatus(): KubernetesResourceState.HealthStatus;
  setHealthStatus(value: KubernetesResourceState.HealthStatus): KubernetesResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): KubernetesResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): KubernetesResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): KubernetesResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesResourceState): KubernetesResourceState.AsObject;
  static serializeBinaryToWriter(message: KubernetesResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesResourceState;
  static deserializeBinaryFromReader(message: KubernetesResourceState, reader: jspb.BinaryReader): KubernetesResourceState;
}

export namespace KubernetesResourceState {
  export type AsObject = {
    id: string,
    ownerIdsList: Array<string>,
    parentIdsList: Array<string>,
    name: string,
    apiVersion: string,
    kind: string,
    namespace: string,
    healthStatus: KubernetesResourceState.HealthStatus,
    healthDescription: string,
    createdAt: number,
    updatedAt: number,
  }

  export enum HealthStatus { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

export class KubernetesResourceStateEvent extends jspb.Message {
  getId(): string;
  setId(value: string): KubernetesResourceStateEvent;

  getApplicationId(): string;
  setApplicationId(value: string): KubernetesResourceStateEvent;

  getType(): KubernetesResourceStateEvent.Type;
  setType(value: KubernetesResourceStateEvent.Type): KubernetesResourceStateEvent;

  getState(): KubernetesResourceState | undefined;
  setState(value?: KubernetesResourceState): KubernetesResourceStateEvent;
  hasState(): boolean;
  clearState(): KubernetesResourceStateEvent;

  getSnapshotVersion(): ApplicationLiveStateVersion | undefined;
  setSnapshotVersion(value?: ApplicationLiveStateVersion): KubernetesResourceStateEvent;
  hasSnapshotVersion(): boolean;
  clearSnapshotVersion(): KubernetesResourceStateEvent;

  getCreatedAt(): number;
  setCreatedAt(value: number): KubernetesResourceStateEvent;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KubernetesResourceStateEvent.AsObject;
  static toObject(includeInstance: boolean, msg: KubernetesResourceStateEvent): KubernetesResourceStateEvent.AsObject;
  static serializeBinaryToWriter(message: KubernetesResourceStateEvent, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KubernetesResourceStateEvent;
  static deserializeBinaryFromReader(message: KubernetesResourceStateEvent, reader: jspb.BinaryReader): KubernetesResourceStateEvent;
}

export namespace KubernetesResourceStateEvent {
  export type AsObject = {
    id: string,
    applicationId: string,
    type: KubernetesResourceStateEvent.Type,
    state?: KubernetesResourceState.AsObject,
    snapshotVersion?: ApplicationLiveStateVersion.AsObject,
    createdAt: number,
  }

  export enum Type { 
    ADD_OR_UPDATED = 0,
    DELETED = 2,
  }
}

export class CloudRunResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): CloudRunResourceState;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): CloudRunResourceState;
  clearOwnerIdsList(): CloudRunResourceState;
  addOwnerIds(value: string, index?: number): CloudRunResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): CloudRunResourceState;
  clearParentIdsList(): CloudRunResourceState;
  addParentIds(value: string, index?: number): CloudRunResourceState;

  getName(): string;
  setName(value: string): CloudRunResourceState;

  getApiVersion(): string;
  setApiVersion(value: string): CloudRunResourceState;

  getKind(): string;
  setKind(value: string): CloudRunResourceState;

  getNamespace(): string;
  setNamespace(value: string): CloudRunResourceState;

  getHealthStatus(): CloudRunResourceState.HealthStatus;
  setHealthStatus(value: CloudRunResourceState.HealthStatus): CloudRunResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): CloudRunResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): CloudRunResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): CloudRunResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CloudRunResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: CloudRunResourceState): CloudRunResourceState.AsObject;
  static serializeBinaryToWriter(message: CloudRunResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CloudRunResourceState;
  static deserializeBinaryFromReader(message: CloudRunResourceState, reader: jspb.BinaryReader): CloudRunResourceState;
}

export namespace CloudRunResourceState {
  export type AsObject = {
    id: string,
    ownerIdsList: Array<string>,
    parentIdsList: Array<string>,
    name: string,
    apiVersion: string,
    kind: string,
    namespace: string,
    healthStatus: CloudRunResourceState.HealthStatus,
    healthDescription: string,
    createdAt: number,
    updatedAt: number,
  }

  export enum HealthStatus { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

export class ECSResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): ECSResourceState;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): ECSResourceState;
  clearOwnerIdsList(): ECSResourceState;
  addOwnerIds(value: string, index?: number): ECSResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): ECSResourceState;
  clearParentIdsList(): ECSResourceState;
  addParentIds(value: string, index?: number): ECSResourceState;

  getName(): string;
  setName(value: string): ECSResourceState;

  getKind(): string;
  setKind(value: string): ECSResourceState;

  getHealthStatus(): ECSResourceState.HealthStatus;
  setHealthStatus(value: ECSResourceState.HealthStatus): ECSResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): ECSResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): ECSResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): ECSResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ECSResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: ECSResourceState): ECSResourceState.AsObject;
  static serializeBinaryToWriter(message: ECSResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ECSResourceState;
  static deserializeBinaryFromReader(message: ECSResourceState, reader: jspb.BinaryReader): ECSResourceState;
}

export namespace ECSResourceState {
  export type AsObject = {
    id: string,
    ownerIdsList: Array<string>,
    parentIdsList: Array<string>,
    name: string,
    kind: string,
    healthStatus: ECSResourceState.HealthStatus,
    healthDescription: string,
    createdAt: number,
    updatedAt: number,
  }

  export enum HealthStatus { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

export class LambdaResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): LambdaResourceState;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): LambdaResourceState;
  clearOwnerIdsList(): LambdaResourceState;
  addOwnerIds(value: string, index?: number): LambdaResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): LambdaResourceState;
  clearParentIdsList(): LambdaResourceState;
  addParentIds(value: string, index?: number): LambdaResourceState;

  getName(): string;
  setName(value: string): LambdaResourceState;

  getKind(): string;
  setKind(value: string): LambdaResourceState;

  getHealthStatus(): LambdaResourceState.HealthStatus;
  setHealthStatus(value: LambdaResourceState.HealthStatus): LambdaResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): LambdaResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): LambdaResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): LambdaResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaResourceState): LambdaResourceState.AsObject;
  static serializeBinaryToWriter(message: LambdaResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaResourceState;
  static deserializeBinaryFromReader(message: LambdaResourceState, reader: jspb.BinaryReader): LambdaResourceState;
}

export namespace LambdaResourceState {
  export type AsObject = {
    id: string,
    ownerIdsList: Array<string>,
    parentIdsList: Array<string>,
    name: string,
    kind: string,
    healthStatus: LambdaResourceState.HealthStatus,
    healthDescription: string,
    createdAt: number,
    updatedAt: number,
  }

  export enum HealthStatus { 
    UNKNOWN = 0,
    HEALTHY = 1,
    OTHER = 2,
  }
}

