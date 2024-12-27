import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';
import * as pkg_model_deployment_pb from 'pipecd/web/model/deployment_pb';


export class Application extends jspb.Message {
  getId(): string;
  setId(value: string): Application;

  getName(): string;
  setName(value: string): Application;

  getPipedId(): string;
  setPipedId(value: string): Application;

  getProjectId(): string;
  setProjectId(value: string): Application;

  getKind(): pkg_model_common_pb.ApplicationKind;
  setKind(value: pkg_model_common_pb.ApplicationKind): Application;

  getGitPath(): pkg_model_common_pb.ApplicationGitPath | undefined;
  setGitPath(value?: pkg_model_common_pb.ApplicationGitPath): Application;
  hasGitPath(): boolean;
  clearGitPath(): Application;

  getCloudProvider(): string;
  setCloudProvider(value: string): Application;

  getPlatformProvider(): string;
  setPlatformProvider(value: string): Application;

  getDeployTargetsList(): Array<string>;
  setDeployTargetsList(value: Array<string>): Application;
  clearDeployTargetsList(): Application;
  addDeployTargets(value: string, index?: number): Application;

  getPluginsList(): Array<string>;
  setPluginsList(value: Array<string>): Application;
  clearPluginsList(): Application;
  addPlugins(value: string, index?: number): Application;

  getDescription(): string;
  setDescription(value: string): Application;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Application;

  getMostRecentlySuccessfulDeployment(): ApplicationDeploymentReference | undefined;
  setMostRecentlySuccessfulDeployment(value?: ApplicationDeploymentReference): Application;
  hasMostRecentlySuccessfulDeployment(): boolean;
  clearMostRecentlySuccessfulDeployment(): Application;

  getMostRecentlyTriggeredDeployment(): ApplicationDeploymentReference | undefined;
  setMostRecentlyTriggeredDeployment(value?: ApplicationDeploymentReference): Application;
  hasMostRecentlyTriggeredDeployment(): boolean;
  clearMostRecentlyTriggeredDeployment(): Application;

  getSyncState(): ApplicationSyncState | undefined;
  setSyncState(value?: ApplicationSyncState): Application;
  hasSyncState(): boolean;
  clearSyncState(): Application;

  getDeploying(): boolean;
  setDeploying(value: boolean): Application;

  getDeletedAt(): number;
  setDeletedAt(value: number): Application;

  getDeleted(): boolean;
  setDeleted(value: boolean): Application;

  getDisabled(): boolean;
  setDisabled(value: boolean): Application;

  getCreatedAt(): number;
  setCreatedAt(value: number): Application;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Application;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Application.AsObject;
  static toObject(includeInstance: boolean, msg: Application): Application.AsObject;
  static serializeBinaryToWriter(message: Application, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Application;
  static deserializeBinaryFromReader(message: Application, reader: jspb.BinaryReader): Application;
}

export namespace Application {
  export type AsObject = {
    id: string,
    name: string,
    pipedId: string,
    projectId: string,
    kind: pkg_model_common_pb.ApplicationKind,
    gitPath?: pkg_model_common_pb.ApplicationGitPath.AsObject,
    cloudProvider: string,
    platformProvider: string,
    deployTargetsList: Array<string>,
    pluginsList: Array<string>,
    description: string,
    labelsMap: Array<[string, string]>,
    mostRecentlySuccessfulDeployment?: ApplicationDeploymentReference.AsObject,
    mostRecentlyTriggeredDeployment?: ApplicationDeploymentReference.AsObject,
    syncState?: ApplicationSyncState.AsObject,
    deploying: boolean,
    deletedAt: number,
    deleted: boolean,
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
  }
}

export class ApplicationSyncState extends jspb.Message {
  getStatus(): ApplicationSyncStatus;
  setStatus(value: ApplicationSyncStatus): ApplicationSyncState;

  getShortReason(): string;
  setShortReason(value: string): ApplicationSyncState;

  getReason(): string;
  setReason(value: string): ApplicationSyncState;

  getHeadDeploymentId(): string;
  setHeadDeploymentId(value: string): ApplicationSyncState;

  getTimestamp(): number;
  setTimestamp(value: number): ApplicationSyncState;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationSyncState.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationSyncState): ApplicationSyncState.AsObject;
  static serializeBinaryToWriter(message: ApplicationSyncState, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationSyncState;
  static deserializeBinaryFromReader(message: ApplicationSyncState, reader: jspb.BinaryReader): ApplicationSyncState;
}

export namespace ApplicationSyncState {
  export type AsObject = {
    status: ApplicationSyncStatus,
    shortReason: string,
    reason: string,
    headDeploymentId: string,
    timestamp: number,
  }
}

export class ApplicationDeploymentReference extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): ApplicationDeploymentReference;

  getTrigger(): pkg_model_deployment_pb.DeploymentTrigger | undefined;
  setTrigger(value?: pkg_model_deployment_pb.DeploymentTrigger): ApplicationDeploymentReference;
  hasTrigger(): boolean;
  clearTrigger(): ApplicationDeploymentReference;

  getSummary(): string;
  setSummary(value: string): ApplicationDeploymentReference;

  getVersion(): string;
  setVersion(value: string): ApplicationDeploymentReference;

  getConfigFilename(): string;
  setConfigFilename(value: string): ApplicationDeploymentReference;

  getVersionsList(): Array<pkg_model_common_pb.ArtifactVersion>;
  setVersionsList(value: Array<pkg_model_common_pb.ArtifactVersion>): ApplicationDeploymentReference;
  clearVersionsList(): ApplicationDeploymentReference;
  addVersions(value?: pkg_model_common_pb.ArtifactVersion, index?: number): pkg_model_common_pb.ArtifactVersion;

  getStartedAt(): number;
  setStartedAt(value: number): ApplicationDeploymentReference;

  getCompletedAt(): number;
  setCompletedAt(value: number): ApplicationDeploymentReference;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationDeploymentReference.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationDeploymentReference): ApplicationDeploymentReference.AsObject;
  static serializeBinaryToWriter(message: ApplicationDeploymentReference, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationDeploymentReference;
  static deserializeBinaryFromReader(message: ApplicationDeploymentReference, reader: jspb.BinaryReader): ApplicationDeploymentReference;
}

export namespace ApplicationDeploymentReference {
  export type AsObject = {
    deploymentId: string,
    trigger?: pkg_model_deployment_pb.DeploymentTrigger.AsObject,
    summary: string,
    version: string,
    configFilename: string,
    versionsList: Array<pkg_model_common_pb.ArtifactVersion.AsObject>,
    startedAt: number,
    completedAt: number,
  }
}

export enum ApplicationSyncStatus { 
  UNKNOWN = 0,
  SYNCED = 1,
  DEPLOYING = 2,
  OUT_OF_SYNC = 3,
  INVALID_CONFIG = 4,
}
