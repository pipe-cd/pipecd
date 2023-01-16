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

  getEcs(): EcsApplicationLiveState | undefined;
  setEcs(value?: EcsApplicationLiveState): ApplicationLiveStateSnapshot;
  hasEcs(): boolean;
  clearEcs(): ApplicationLiveStateSnapshot;

  getVersion(): ApplicationLiveStateVersion | undefined;
  setVersion(value?: ApplicationLiveStateVersion): ApplicationLiveStateSnapshot;
  hasVersion(): boolean;
  clearVersion(): ApplicationLiveStateSnapshot;

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
    ecs?: EcsApplicationLiveState.AsObject,
    version?: ApplicationLiveStateVersion.AsObject,
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

export class LambdaApplicationLiveState extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LambdaApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: LambdaApplicationLiveState): LambdaApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: LambdaApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LambdaApplicationLiveState;
  static deserializeBinaryFromReader(message: LambdaApplicationLiveState, reader: jspb.BinaryReader): LambdaApplicationLiveState;
}

export namespace LambdaApplicationLiveState {
  export type AsObject = {
  }
}

export class EcsApplicationLiveState extends jspb.Message {
  getResourcesList(): Array<EcsResourceState>;
  setResourcesList(value: Array<EcsResourceState>): EcsApplicationLiveState;
  clearResourcesList(): EcsApplicationLiveState;
  addResources(value?: EcsResourceState, index?: number): EcsResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EcsApplicationLiveState.AsObject;
  static toObject(includeInstance: boolean, msg: EcsApplicationLiveState): EcsApplicationLiveState.AsObject;
  static serializeBinaryToWriter(message: EcsApplicationLiveState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EcsApplicationLiveState;
  static deserializeBinaryFromReader(message: EcsApplicationLiveState, reader: jspb.BinaryReader): EcsApplicationLiveState;
}

export namespace EcsApplicationLiveState {
  export type AsObject = {
    resourcesList: Array<EcsResourceState.AsObject>,
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

export class EcsResourceState extends jspb.Message {
  getId(): string;
  setId(value: string): EcsResourceState;

  getOwnerIdsList(): Array<string>;
  setOwnerIdsList(value: Array<string>): EcsResourceState;
  clearOwnerIdsList(): EcsResourceState;
  addOwnerIds(value: string, index?: number): EcsResourceState;

  getParentIdsList(): Array<string>;
  setParentIdsList(value: Array<string>): EcsResourceState;
  clearParentIdsList(): EcsResourceState;
  addParentIds(value: string, index?: number): EcsResourceState;

  getName(): string;
  setName(value: string): EcsResourceState;

  getApiVersion(): string;
  setApiVersion(value: string): EcsResourceState;

  getKind(): string;
  setKind(value: string): EcsResourceState;

  getNamespace(): string;
  setNamespace(value: string): EcsResourceState;

  getHealthStatus(): EcsResourceState.HealthStatus;
  setHealthStatus(value: EcsResourceState.HealthStatus): EcsResourceState;

  getHealthDescription(): string;
  setHealthDescription(value: string): EcsResourceState;

  getCreatedAt(): number;
  setCreatedAt(value: number): EcsResourceState;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): EcsResourceState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EcsResourceState.AsObject;
  static toObject(includeInstance: boolean, msg: EcsResourceState): EcsResourceState.AsObject;
  static serializeBinaryToWriter(message: EcsResourceState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EcsResourceState;
  static deserializeBinaryFromReader(message: EcsResourceState, reader: jspb.BinaryReader): EcsResourceState;
}

export namespace EcsResourceState {
  export type AsObject = {
    id: string,
    ownerIdsList: Array<string>,
    parentIdsList: Array<string>,
    name: string,
    apiVersion: string,
    kind: string,
    namespace: string,
    healthStatus: EcsResourceState.HealthStatus,
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

