import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';


export class Deployment extends jspb.Message {
  getId(): string;
  setId(value: string): Deployment;

  getApplicationId(): string;
  setApplicationId(value: string): Deployment;

  getApplicationName(): string;
  setApplicationName(value: string): Deployment;

  getPipedId(): string;
  setPipedId(value: string): Deployment;

  getProjectId(): string;
  setProjectId(value: string): Deployment;

  getKind(): pkg_model_common_pb.ApplicationKind;
  setKind(value: pkg_model_common_pb.ApplicationKind): Deployment;

  getGitPath(): pkg_model_common_pb.ApplicationGitPath | undefined;
  setGitPath(value?: pkg_model_common_pb.ApplicationGitPath): Deployment;
  hasGitPath(): boolean;
  clearGitPath(): Deployment;

  getCloudProvider(): string;
  setCloudProvider(value: string): Deployment;

  getPlatformProvider(): string;
  setPlatformProvider(value: string): Deployment;

  getDeployTargetsList(): Array<string>;
  setDeployTargetsList(value: Array<string>): Deployment;
  clearDeployTargetsList(): Deployment;
  addDeployTargets(value: string, index?: number): Deployment;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Deployment;

  getTrigger(): DeploymentTrigger | undefined;
  setTrigger(value?: DeploymentTrigger): Deployment;
  hasTrigger(): boolean;
  clearTrigger(): Deployment;

  getSummary(): string;
  setSummary(value: string): Deployment;

  getVersion(): string;
  setVersion(value: string): Deployment;

  getVersionsList(): Array<pkg_model_common_pb.ArtifactVersion>;
  setVersionsList(value: Array<pkg_model_common_pb.ArtifactVersion>): Deployment;
  clearVersionsList(): Deployment;
  addVersions(value?: pkg_model_common_pb.ArtifactVersion, index?: number): pkg_model_common_pb.ArtifactVersion;

  getRunningCommitHash(): string;
  setRunningCommitHash(value: string): Deployment;

  getRunningConfigFilename(): string;
  setRunningConfigFilename(value: string): Deployment;

  getStatus(): DeploymentStatus;
  setStatus(value: DeploymentStatus): Deployment;

  getStatusReason(): string;
  setStatusReason(value: string): Deployment;

  getStagesList(): Array<PipelineStage>;
  setStagesList(value: Array<PipelineStage>): Deployment;
  clearStagesList(): Deployment;
  addStages(value?: PipelineStage, index?: number): PipelineStage;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): Deployment;

  getPluginMetadataMap(): jspb.Map<string, KeyValues>;
  clearPluginMetadataMap(): Deployment;

  getDeploymentChainId(): string;
  setDeploymentChainId(value: string): Deployment;

  getDeploymentChainBlockIndex(): number;
  setDeploymentChainBlockIndex(value: number): Deployment;

  getCompletedAt(): number;
  setCompletedAt(value: number): Deployment;

  getCreatedAt(): number;
  setCreatedAt(value: number): Deployment;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Deployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Deployment.AsObject;
  static toObject(includeInstance: boolean, msg: Deployment): Deployment.AsObject;
  static serializeBinaryToWriter(message: Deployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Deployment;
  static deserializeBinaryFromReader(message: Deployment, reader: jspb.BinaryReader): Deployment;
}

export namespace Deployment {
  export type AsObject = {
    id: string,
    applicationId: string,
    applicationName: string,
    pipedId: string,
    projectId: string,
    kind: pkg_model_common_pb.ApplicationKind,
    gitPath?: pkg_model_common_pb.ApplicationGitPath.AsObject,
    cloudProvider: string,
    platformProvider: string,
    deployTargetsList: Array<string>,
    labelsMap: Array<[string, string]>,
    trigger?: DeploymentTrigger.AsObject,
    summary: string,
    version: string,
    versionsList: Array<pkg_model_common_pb.ArtifactVersion.AsObject>,
    runningCommitHash: string,
    runningConfigFilename: string,
    status: DeploymentStatus,
    statusReason: string,
    stagesList: Array<PipelineStage.AsObject>,
    metadataMap: Array<[string, string]>,
    pluginMetadataMap: Array<[string, KeyValues.AsObject]>,
    deploymentChainId: string,
    deploymentChainBlockIndex: number,
    completedAt: number,
    createdAt: number,
    updatedAt: number,
  }
}

export class DeploymentTrigger extends jspb.Message {
  getCommit(): Commit | undefined;
  setCommit(value?: Commit): DeploymentTrigger;
  hasCommit(): boolean;
  clearCommit(): DeploymentTrigger;

  getCommander(): string;
  setCommander(value: string): DeploymentTrigger;

  getTimestamp(): number;
  setTimestamp(value: number): DeploymentTrigger;

  getSyncStrategy(): pkg_model_common_pb.SyncStrategy;
  setSyncStrategy(value: pkg_model_common_pb.SyncStrategy): DeploymentTrigger;

  getStrategySummary(): string;
  setStrategySummary(value: string): DeploymentTrigger;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentTrigger.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentTrigger): DeploymentTrigger.AsObject;
  static serializeBinaryToWriter(message: DeploymentTrigger, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentTrigger;
  static deserializeBinaryFromReader(message: DeploymentTrigger, reader: jspb.BinaryReader): DeploymentTrigger;
}

export namespace DeploymentTrigger {
  export type AsObject = {
    commit?: Commit.AsObject,
    commander: string,
    timestamp: number,
    syncStrategy: pkg_model_common_pb.SyncStrategy,
    strategySummary: string,
  }
}

export class PipelineStage extends jspb.Message {
  getId(): string;
  setId(value: string): PipelineStage;

  getName(): string;
  setName(value: string): PipelineStage;

  getDesc(): string;
  setDesc(value: string): PipelineStage;

  getIndex(): number;
  setIndex(value: number): PipelineStage;

  getPredefined(): boolean;
  setPredefined(value: boolean): PipelineStage;

  getRequiresList(): Array<string>;
  setRequiresList(value: Array<string>): PipelineStage;
  clearRequiresList(): PipelineStage;
  addRequires(value: string, index?: number): PipelineStage;

  getVisible(): boolean;
  setVisible(value: boolean): PipelineStage;

  getStatus(): StageStatus;
  setStatus(value: StageStatus): PipelineStage;

  getStatusReason(): string;
  setStatusReason(value: string): PipelineStage;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): PipelineStage;

  getRetriedCount(): number;
  setRetriedCount(value: number): PipelineStage;

  getRollback(): boolean;
  setRollback(value: boolean): PipelineStage;

  getCompletedAt(): number;
  setCompletedAt(value: number): PipelineStage;

  getCreatedAt(): number;
  setCreatedAt(value: number): PipelineStage;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): PipelineStage;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PipelineStage.AsObject;
  static toObject(includeInstance: boolean, msg: PipelineStage): PipelineStage.AsObject;
  static serializeBinaryToWriter(message: PipelineStage, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PipelineStage;
  static deserializeBinaryFromReader(message: PipelineStage, reader: jspb.BinaryReader): PipelineStage;
}

export namespace PipelineStage {
  export type AsObject = {
    id: string,
    name: string,
    desc: string,
    index: number,
    predefined: boolean,
    requiresList: Array<string>,
    visible: boolean,
    status: StageStatus,
    statusReason: string,
    metadataMap: Array<[string, string]>,
    retriedCount: number,
    rollback: boolean,
    completedAt: number,
    createdAt: number,
    updatedAt: number,
  }
}

export class Commit extends jspb.Message {
  getHash(): string;
  setHash(value: string): Commit;

  getMessage(): string;
  setMessage(value: string): Commit;

  getAuthor(): string;
  setAuthor(value: string): Commit;

  getBranch(): string;
  setBranch(value: string): Commit;

  getPullRequest(): number;
  setPullRequest(value: number): Commit;

  getUrl(): string;
  setUrl(value: string): Commit;

  getCreatedAt(): number;
  setCreatedAt(value: number): Commit;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Commit.AsObject;
  static toObject(includeInstance: boolean, msg: Commit): Commit.AsObject;
  static serializeBinaryToWriter(message: Commit, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Commit;
  static deserializeBinaryFromReader(message: Commit, reader: jspb.BinaryReader): Commit;
}

export namespace Commit {
  export type AsObject = {
    hash: string,
    message: string,
    author: string,
    branch: string,
    pullRequest: number,
    url: string,
    createdAt: number,
  }
}

export class KeyValues extends jspb.Message {
  getKeyvaluesMap(): jspb.Map<string, string>;
  clearKeyvaluesMap(): KeyValues;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): KeyValues.AsObject;
  static toObject(includeInstance: boolean, msg: KeyValues): KeyValues.AsObject;
  static serializeBinaryToWriter(message: KeyValues, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): KeyValues;
  static deserializeBinaryFromReader(message: KeyValues, reader: jspb.BinaryReader): KeyValues;
}

export namespace KeyValues {
  export type AsObject = {
    keyvaluesMap: Array<[string, string]>,
  }
}

export enum DeploymentStatus { 
  DEPLOYMENT_PENDING = 0,
  DEPLOYMENT_PLANNED = 1,
  DEPLOYMENT_RUNNING = 2,
  DEPLOYMENT_ROLLING_BACK = 3,
  DEPLOYMENT_SUCCESS = 4,
  DEPLOYMENT_FAILURE = 5,
  DEPLOYMENT_CANCELLED = 6,
}
export enum StageStatus { 
  STAGE_NOT_STARTED_YET = 0,
  STAGE_RUNNING = 1,
  STAGE_SUCCESS = 2,
  STAGE_FAILURE = 3,
  STAGE_CANCELLED = 4,
  STAGE_SKIPPED = 5,
  STAGE_EXITED = 6,
}
export enum TriggerKind { 
  ON_COMMIT = 0,
  ON_COMMAND = 1,
  ON_OUT_OF_SYNC = 2,
  ON_CHAIN = 3,
}
