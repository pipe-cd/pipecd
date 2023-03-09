import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';


export class PlanPreviewCommandResult extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): PlanPreviewCommandResult;

  getPipedId(): string;
  setPipedId(value: string): PlanPreviewCommandResult;

  getPipedUrl(): string;
  setPipedUrl(value: string): PlanPreviewCommandResult;

  getResultsList(): Array<ApplicationPlanPreviewResult>;
  setResultsList(value: Array<ApplicationPlanPreviewResult>): PlanPreviewCommandResult;
  clearResultsList(): PlanPreviewCommandResult;
  addResults(value?: ApplicationPlanPreviewResult, index?: number): ApplicationPlanPreviewResult;

  getError(): string;
  setError(value: string): PlanPreviewCommandResult;

  getPipedName(): string;
  setPipedName(value: string): PlanPreviewCommandResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PlanPreviewCommandResult.AsObject;
  static toObject(includeInstance: boolean, msg: PlanPreviewCommandResult): PlanPreviewCommandResult.AsObject;
  static serializeBinaryToWriter(message: PlanPreviewCommandResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PlanPreviewCommandResult;
  static deserializeBinaryFromReader(message: PlanPreviewCommandResult, reader: jspb.BinaryReader): PlanPreviewCommandResult;
}

export namespace PlanPreviewCommandResult {
  export type AsObject = {
    commandId: string,
    pipedId: string,
    pipedUrl: string,
    resultsList: Array<ApplicationPlanPreviewResult.AsObject>,
    error: string,
    pipedName: string,
  }
}

export class ApplicationPlanPreviewResult extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): ApplicationPlanPreviewResult;

  getApplicationName(): string;
  setApplicationName(value: string): ApplicationPlanPreviewResult;

  getApplicationUrl(): string;
  setApplicationUrl(value: string): ApplicationPlanPreviewResult;

  getApplicationKind(): pkg_model_common_pb.ApplicationKind;
  setApplicationKind(value: pkg_model_common_pb.ApplicationKind): ApplicationPlanPreviewResult;

  getApplicationDirectory(): string;
  setApplicationDirectory(value: string): ApplicationPlanPreviewResult;

  getPipedId(): string;
  setPipedId(value: string): ApplicationPlanPreviewResult;

  getProjectId(): string;
  setProjectId(value: string): ApplicationPlanPreviewResult;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): ApplicationPlanPreviewResult;

  getHeadBranch(): string;
  setHeadBranch(value: string): ApplicationPlanPreviewResult;

  getHeadCommit(): string;
  setHeadCommit(value: string): ApplicationPlanPreviewResult;

  getSyncStrategy(): pkg_model_common_pb.SyncStrategy;
  setSyncStrategy(value: pkg_model_common_pb.SyncStrategy): ApplicationPlanPreviewResult;

  getPlanSummary(): Uint8Array | string;
  getPlanSummary_asU8(): Uint8Array;
  getPlanSummary_asB64(): string;
  setPlanSummary(value: Uint8Array | string): ApplicationPlanPreviewResult;

  getPlanDetails(): Uint8Array | string;
  getPlanDetails_asU8(): Uint8Array;
  getPlanDetails_asB64(): string;
  setPlanDetails(value: Uint8Array | string): ApplicationPlanPreviewResult;

  getNoChange(): boolean;
  setNoChange(value: boolean): ApplicationPlanPreviewResult;

  getError(): string;
  setError(value: string): ApplicationPlanPreviewResult;

  getCreatedAt(): number;
  setCreatedAt(value: number): ApplicationPlanPreviewResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationPlanPreviewResult.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationPlanPreviewResult): ApplicationPlanPreviewResult.AsObject;
  static serializeBinaryToWriter(message: ApplicationPlanPreviewResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationPlanPreviewResult;
  static deserializeBinaryFromReader(message: ApplicationPlanPreviewResult, reader: jspb.BinaryReader): ApplicationPlanPreviewResult;
}

export namespace ApplicationPlanPreviewResult {
  export type AsObject = {
    applicationId: string,
    applicationName: string,
    applicationUrl: string,
    applicationKind: pkg_model_common_pb.ApplicationKind,
    applicationDirectory: string,
    pipedId: string,
    projectId: string,
    labelsMap: Array<[string, string]>,
    headBranch: string,
    headCommit: string,
    syncStrategy: pkg_model_common_pb.SyncStrategy,
    planSummary: Uint8Array | string,
    planDetails: Uint8Array | string,
    noChange: boolean,
    error: string,
    createdAt: number,
  }
}

