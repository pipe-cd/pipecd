import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';


export class Command extends jspb.Message {
  getId(): string;
  setId(value: string): Command;

  getPipedId(): string;
  setPipedId(value: string): Command;

  getApplicationId(): string;
  setApplicationId(value: string): Command;

  getDeploymentId(): string;
  setDeploymentId(value: string): Command;

  getStageId(): string;
  setStageId(value: string): Command;

  getCommander(): string;
  setCommander(value: string): Command;

  getProjectId(): string;
  setProjectId(value: string): Command;

  getStatus(): CommandStatus;
  setStatus(value: CommandStatus): Command;

  getMetadataMap(): jspb.Map<string, string>;
  clearMetadataMap(): Command;

  getHandledAt(): number;
  setHandledAt(value: number): Command;

  getType(): Command.Type;
  setType(value: Command.Type): Command;

  getSyncApplication(): Command.SyncApplication | undefined;
  setSyncApplication(value?: Command.SyncApplication): Command;
  hasSyncApplication(): boolean;
  clearSyncApplication(): Command;

  getUpdateApplicationConfig(): Command.UpdateApplicationConfig | undefined;
  setUpdateApplicationConfig(value?: Command.UpdateApplicationConfig): Command;
  hasUpdateApplicationConfig(): boolean;
  clearUpdateApplicationConfig(): Command;

  getCancelDeployment(): Command.CancelDeployment | undefined;
  setCancelDeployment(value?: Command.CancelDeployment): Command;
  hasCancelDeployment(): boolean;
  clearCancelDeployment(): Command;

  getApproveStage(): Command.ApproveStage | undefined;
  setApproveStage(value?: Command.ApproveStage): Command;
  hasApproveStage(): boolean;
  clearApproveStage(): Command;

  getBuildPlanPreview(): Command.BuildPlanPreview | undefined;
  setBuildPlanPreview(value?: Command.BuildPlanPreview): Command;
  hasBuildPlanPreview(): boolean;
  clearBuildPlanPreview(): Command;

  getChainSyncApplication(): Command.ChainSyncApplication | undefined;
  setChainSyncApplication(value?: Command.ChainSyncApplication): Command;
  hasChainSyncApplication(): boolean;
  clearChainSyncApplication(): Command;

  getSkipStage(): Command.SkipStage | undefined;
  setSkipStage(value?: Command.SkipStage): Command;
  hasSkipStage(): boolean;
  clearSkipStage(): Command;

  getRestartPiped(): Command.RestartPiped | undefined;
  setRestartPiped(value?: Command.RestartPiped): Command;
  hasRestartPiped(): boolean;
  clearRestartPiped(): Command;

  getCancelPlanPreview(): Command.CancelPlanPreview | undefined;
  setCancelPlanPreview(value?: Command.CancelPlanPreview): Command;
  hasCancelPlanPreview(): boolean;
  clearCancelPlanPreview(): Command;

  getCreatedAt(): number;
  setCreatedAt(value: number): Command;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Command;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Command.AsObject;
  static toObject(includeInstance: boolean, msg: Command): Command.AsObject;
  static serializeBinaryToWriter(message: Command, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Command;
  static deserializeBinaryFromReader(message: Command, reader: jspb.BinaryReader): Command;
}

export namespace Command {
  export type AsObject = {
    id: string,
    pipedId: string,
    applicationId: string,
    deploymentId: string,
    stageId: string,
    commander: string,
    projectId: string,
    status: CommandStatus,
    metadataMap: Array<[string, string]>,
    handledAt: number,
    type: Command.Type,
    syncApplication?: Command.SyncApplication.AsObject,
    updateApplicationConfig?: Command.UpdateApplicationConfig.AsObject,
    cancelDeployment?: Command.CancelDeployment.AsObject,
    approveStage?: Command.ApproveStage.AsObject,
    buildPlanPreview?: Command.BuildPlanPreview.AsObject,
    chainSyncApplication?: Command.ChainSyncApplication.AsObject,
    skipStage?: Command.SkipStage.AsObject,
    restartPiped?: Command.RestartPiped.AsObject,
    cancelPlanPreview?: Command.CancelPlanPreview.AsObject,
    createdAt: number,
    updatedAt: number,
  }

  export class SyncApplication extends jspb.Message {
    getApplicationId(): string;
    setApplicationId(value: string): SyncApplication;

    getSyncStrategy(): pkg_model_common_pb.SyncStrategy;
    setSyncStrategy(value: pkg_model_common_pb.SyncStrategy): SyncApplication;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SyncApplication.AsObject;
    static toObject(includeInstance: boolean, msg: SyncApplication): SyncApplication.AsObject;
    static serializeBinaryToWriter(message: SyncApplication, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SyncApplication;
    static deserializeBinaryFromReader(message: SyncApplication, reader: jspb.BinaryReader): SyncApplication;
  }

  export namespace SyncApplication {
    export type AsObject = {
      applicationId: string,
      syncStrategy: pkg_model_common_pb.SyncStrategy,
    }
  }


  export class UpdateApplicationConfig extends jspb.Message {
    getApplicationId(): string;
    setApplicationId(value: string): UpdateApplicationConfig;

    getConfigPath(): string;
    setConfigPath(value: string): UpdateApplicationConfig;

    getConfig(): string;
    setConfig(value: string): UpdateApplicationConfig;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdateApplicationConfig.AsObject;
    static toObject(includeInstance: boolean, msg: UpdateApplicationConfig): UpdateApplicationConfig.AsObject;
    static serializeBinaryToWriter(message: UpdateApplicationConfig, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdateApplicationConfig;
    static deserializeBinaryFromReader(message: UpdateApplicationConfig, reader: jspb.BinaryReader): UpdateApplicationConfig;
  }

  export namespace UpdateApplicationConfig {
    export type AsObject = {
      applicationId: string,
      configPath: string,
      config: string,
    }
  }


  export class CancelDeployment extends jspb.Message {
    getDeploymentId(): string;
    setDeploymentId(value: string): CancelDeployment;

    getForceRollback(): boolean;
    setForceRollback(value: boolean): CancelDeployment;

    getForceNoRollback(): boolean;
    setForceNoRollback(value: boolean): CancelDeployment;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CancelDeployment.AsObject;
    static toObject(includeInstance: boolean, msg: CancelDeployment): CancelDeployment.AsObject;
    static serializeBinaryToWriter(message: CancelDeployment, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CancelDeployment;
    static deserializeBinaryFromReader(message: CancelDeployment, reader: jspb.BinaryReader): CancelDeployment;
  }

  export namespace CancelDeployment {
    export type AsObject = {
      deploymentId: string,
      forceRollback: boolean,
      forceNoRollback: boolean,
    }
  }


  export class ApproveStage extends jspb.Message {
    getDeploymentId(): string;
    setDeploymentId(value: string): ApproveStage;

    getStageId(): string;
    setStageId(value: string): ApproveStage;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ApproveStage.AsObject;
    static toObject(includeInstance: boolean, msg: ApproveStage): ApproveStage.AsObject;
    static serializeBinaryToWriter(message: ApproveStage, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ApproveStage;
    static deserializeBinaryFromReader(message: ApproveStage, reader: jspb.BinaryReader): ApproveStage;
  }

  export namespace ApproveStage {
    export type AsObject = {
      deploymentId: string,
      stageId: string,
    }
  }


  export class BuildPlanPreview extends jspb.Message {
    getRepositoryId(): string;
    setRepositoryId(value: string): BuildPlanPreview;

    getHeadBranch(): string;
    setHeadBranch(value: string): BuildPlanPreview;

    getHeadCommit(): string;
    setHeadCommit(value: string): BuildPlanPreview;

    getBaseBranch(): string;
    setBaseBranch(value: string): BuildPlanPreview;

    getTimeout(): number;
    setTimeout(value: number): BuildPlanPreview;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BuildPlanPreview.AsObject;
    static toObject(includeInstance: boolean, msg: BuildPlanPreview): BuildPlanPreview.AsObject;
    static serializeBinaryToWriter(message: BuildPlanPreview, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BuildPlanPreview;
    static deserializeBinaryFromReader(message: BuildPlanPreview, reader: jspb.BinaryReader): BuildPlanPreview;
  }

  export namespace BuildPlanPreview {
    export type AsObject = {
      repositoryId: string,
      headBranch: string,
      headCommit: string,
      baseBranch: string,
      timeout: number,
    }
  }


  export class CancelPlanPreview extends jspb.Message {
    getRepositoryId(): string;
    setRepositoryId(value: string): CancelPlanPreview;

    getHeadBranch(): string;
    setHeadBranch(value: string): CancelPlanPreview;

    getHeadCommit(): string;
    setHeadCommit(value: string): CancelPlanPreview;

    getBaseBranch(): string;
    setBaseBranch(value: string): CancelPlanPreview;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CancelPlanPreview.AsObject;
    static toObject(includeInstance: boolean, msg: CancelPlanPreview): CancelPlanPreview.AsObject;
    static serializeBinaryToWriter(message: CancelPlanPreview, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CancelPlanPreview;
    static deserializeBinaryFromReader(message: CancelPlanPreview, reader: jspb.BinaryReader): CancelPlanPreview;
  }

  export namespace CancelPlanPreview {
    export type AsObject = {
      repositoryId: string,
      headBranch: string,
      headCommit: string,
      baseBranch: string,
    }
  }


  export class ChainSyncApplication extends jspb.Message {
    getDeploymentChainId(): string;
    setDeploymentChainId(value: string): ChainSyncApplication;

    getBlockIndex(): number;
    setBlockIndex(value: number): ChainSyncApplication;

    getApplicationId(): string;
    setApplicationId(value: string): ChainSyncApplication;

    getSyncStrategy(): pkg_model_common_pb.SyncStrategy;
    setSyncStrategy(value: pkg_model_common_pb.SyncStrategy): ChainSyncApplication;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChainSyncApplication.AsObject;
    static toObject(includeInstance: boolean, msg: ChainSyncApplication): ChainSyncApplication.AsObject;
    static serializeBinaryToWriter(message: ChainSyncApplication, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChainSyncApplication;
    static deserializeBinaryFromReader(message: ChainSyncApplication, reader: jspb.BinaryReader): ChainSyncApplication;
  }

  export namespace ChainSyncApplication {
    export type AsObject = {
      deploymentChainId: string,
      blockIndex: number,
      applicationId: string,
      syncStrategy: pkg_model_common_pb.SyncStrategy,
    }
  }


  export class SkipStage extends jspb.Message {
    getDeploymentId(): string;
    setDeploymentId(value: string): SkipStage;

    getStageId(): string;
    setStageId(value: string): SkipStage;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SkipStage.AsObject;
    static toObject(includeInstance: boolean, msg: SkipStage): SkipStage.AsObject;
    static serializeBinaryToWriter(message: SkipStage, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SkipStage;
    static deserializeBinaryFromReader(message: SkipStage, reader: jspb.BinaryReader): SkipStage;
  }

  export namespace SkipStage {
    export type AsObject = {
      deploymentId: string,
      stageId: string,
    }
  }


  export class RestartPiped extends jspb.Message {
    getPipedId(): string;
    setPipedId(value: string): RestartPiped;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RestartPiped.AsObject;
    static toObject(includeInstance: boolean, msg: RestartPiped): RestartPiped.AsObject;
    static serializeBinaryToWriter(message: RestartPiped, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RestartPiped;
    static deserializeBinaryFromReader(message: RestartPiped, reader: jspb.BinaryReader): RestartPiped;
  }

  export namespace RestartPiped {
    export type AsObject = {
      pipedId: string,
    }
  }


  export enum Type { 
    SYNC_APPLICATION = 0,
    UPDATE_APPLICATION_CONFIG = 1,
    CANCEL_DEPLOYMENT = 2,
    APPROVE_STAGE = 3,
    BUILD_PLAN_PREVIEW = 4,
    CHAIN_SYNC_APPLICATION = 5,
    SKIP_STAGE = 6,
    RESTART_PIPED = 7,
    CANCEL_PLAN_PREVIEW = 8,
  }
}

export enum CommandStatus { 
  COMMAND_NOT_HANDLED_YET = 0,
  COMMAND_SUCCEEDED = 1,
  COMMAND_FAILED = 2,
  COMMAND_TIMEOUT = 3,
}
