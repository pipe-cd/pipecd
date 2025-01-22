import * as jspb from 'google-protobuf'


import * as pkg_model_application_pb from 'pipecd/web/model/application_pb';
import * as pkg_model_deployment_pb from 'pipecd/web/model/deployment_pb';


export class NotificationEventDeploymentTriggered extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentTriggered;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentTriggered;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentTriggered;
  clearMentionedAccountsList(): NotificationEventDeploymentTriggered;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentTriggered;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentTriggered;
  clearMentionedGroupsList(): NotificationEventDeploymentTriggered;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentTriggered;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentTriggered.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentTriggered): NotificationEventDeploymentTriggered.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentTriggered, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentTriggered;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentTriggered, reader: jspb.BinaryReader): NotificationEventDeploymentTriggered;
}

export namespace NotificationEventDeploymentTriggered {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentPlanned extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentPlanned;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentPlanned;

  getSummary(): string;
  setSummary(value: string): NotificationEventDeploymentPlanned;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentPlanned;
  clearMentionedAccountsList(): NotificationEventDeploymentPlanned;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentPlanned;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentPlanned;
  clearMentionedGroupsList(): NotificationEventDeploymentPlanned;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentPlanned;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentPlanned.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentPlanned): NotificationEventDeploymentPlanned.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentPlanned, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentPlanned;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentPlanned, reader: jspb.BinaryReader): NotificationEventDeploymentPlanned;
}

export namespace NotificationEventDeploymentPlanned {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    summary: string,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentStarted extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentStarted;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentStarted;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentStarted;
  clearMentionedAccountsList(): NotificationEventDeploymentStarted;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentStarted;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentStarted;
  clearMentionedGroupsList(): NotificationEventDeploymentStarted;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentStarted;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentStarted.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentStarted): NotificationEventDeploymentStarted.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentStarted, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentStarted;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentStarted, reader: jspb.BinaryReader): NotificationEventDeploymentStarted;
}

export namespace NotificationEventDeploymentStarted {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentApproved extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentApproved;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentApproved;

  getApprover(): string;
  setApprover(value: string): NotificationEventDeploymentApproved;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentApproved;
  clearMentionedAccountsList(): NotificationEventDeploymentApproved;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentApproved;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentApproved;
  clearMentionedGroupsList(): NotificationEventDeploymentApproved;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentApproved;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentApproved.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentApproved): NotificationEventDeploymentApproved.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentApproved, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentApproved;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentApproved, reader: jspb.BinaryReader): NotificationEventDeploymentApproved;
}

export namespace NotificationEventDeploymentApproved {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    approver: string,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentRollingBack extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentRollingBack;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentRollingBack;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentRollingBack.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentRollingBack): NotificationEventDeploymentRollingBack.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentRollingBack, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentRollingBack;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentRollingBack, reader: jspb.BinaryReader): NotificationEventDeploymentRollingBack;
}

export namespace NotificationEventDeploymentRollingBack {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
  }
}

export class NotificationEventDeploymentSucceeded extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentSucceeded;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentSucceeded;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentSucceeded;
  clearMentionedAccountsList(): NotificationEventDeploymentSucceeded;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentSucceeded;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentSucceeded;
  clearMentionedGroupsList(): NotificationEventDeploymentSucceeded;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentSucceeded;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentSucceeded.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentSucceeded): NotificationEventDeploymentSucceeded.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentSucceeded, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentSucceeded;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentSucceeded, reader: jspb.BinaryReader): NotificationEventDeploymentSucceeded;
}

export namespace NotificationEventDeploymentSucceeded {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentFailed extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentFailed;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentFailed;

  getReason(): string;
  setReason(value: string): NotificationEventDeploymentFailed;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentFailed;
  clearMentionedAccountsList(): NotificationEventDeploymentFailed;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentFailed;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentFailed;
  clearMentionedGroupsList(): NotificationEventDeploymentFailed;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentFailed;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentFailed.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentFailed): NotificationEventDeploymentFailed.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentFailed, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentFailed;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentFailed, reader: jspb.BinaryReader): NotificationEventDeploymentFailed;
}

export namespace NotificationEventDeploymentFailed {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    reason: string,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentCancelled extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentCancelled;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentCancelled;

  getCommander(): string;
  setCommander(value: string): NotificationEventDeploymentCancelled;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentCancelled;
  clearMentionedAccountsList(): NotificationEventDeploymentCancelled;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentCancelled;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentCancelled;
  clearMentionedGroupsList(): NotificationEventDeploymentCancelled;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentCancelled;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentCancelled.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentCancelled): NotificationEventDeploymentCancelled.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentCancelled, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentCancelled;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentCancelled, reader: jspb.BinaryReader): NotificationEventDeploymentCancelled;
}

export namespace NotificationEventDeploymentCancelled {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    commander: string,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentWaitApproval extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): NotificationEventDeploymentWaitApproval;
  hasDeployment(): boolean;
  clearDeployment(): NotificationEventDeploymentWaitApproval;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentWaitApproval;
  clearMentionedAccountsList(): NotificationEventDeploymentWaitApproval;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentWaitApproval;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentWaitApproval;
  clearMentionedGroupsList(): NotificationEventDeploymentWaitApproval;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentWaitApproval;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentWaitApproval.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentWaitApproval): NotificationEventDeploymentWaitApproval.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentWaitApproval, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentWaitApproval;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentWaitApproval, reader: jspb.BinaryReader): NotificationEventDeploymentWaitApproval;
}

export namespace NotificationEventDeploymentWaitApproval {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventDeploymentTriggerFailed extends jspb.Message {
  getApplication(): pkg_model_application_pb.Application | undefined;
  setApplication(value?: pkg_model_application_pb.Application): NotificationEventDeploymentTriggerFailed;
  hasApplication(): boolean;
  clearApplication(): NotificationEventDeploymentTriggerFailed;

  getCommitHash(): string;
  setCommitHash(value: string): NotificationEventDeploymentTriggerFailed;

  getCommitMessage(): string;
  setCommitMessage(value: string): NotificationEventDeploymentTriggerFailed;

  getReason(): string;
  setReason(value: string): NotificationEventDeploymentTriggerFailed;

  getMentionedAccountsList(): Array<string>;
  setMentionedAccountsList(value: Array<string>): NotificationEventDeploymentTriggerFailed;
  clearMentionedAccountsList(): NotificationEventDeploymentTriggerFailed;
  addMentionedAccounts(value: string, index?: number): NotificationEventDeploymentTriggerFailed;

  getMentionedGroupsList(): Array<string>;
  setMentionedGroupsList(value: Array<string>): NotificationEventDeploymentTriggerFailed;
  clearMentionedGroupsList(): NotificationEventDeploymentTriggerFailed;
  addMentionedGroups(value: string, index?: number): NotificationEventDeploymentTriggerFailed;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventDeploymentTriggerFailed.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventDeploymentTriggerFailed): NotificationEventDeploymentTriggerFailed.AsObject;
  static serializeBinaryToWriter(message: NotificationEventDeploymentTriggerFailed, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventDeploymentTriggerFailed;
  static deserializeBinaryFromReader(message: NotificationEventDeploymentTriggerFailed, reader: jspb.BinaryReader): NotificationEventDeploymentTriggerFailed;
}

export namespace NotificationEventDeploymentTriggerFailed {
  export type AsObject = {
    application?: pkg_model_application_pb.Application.AsObject,
    commitHash: string,
    commitMessage: string,
    reason: string,
    mentionedAccountsList: Array<string>,
    mentionedGroupsList: Array<string>,
  }
}

export class NotificationEventApplicationSynced extends jspb.Message {
  getApplication(): pkg_model_application_pb.Application | undefined;
  setApplication(value?: pkg_model_application_pb.Application): NotificationEventApplicationSynced;
  hasApplication(): boolean;
  clearApplication(): NotificationEventApplicationSynced;

  getState(): pkg_model_application_pb.ApplicationSyncState | undefined;
  setState(value?: pkg_model_application_pb.ApplicationSyncState): NotificationEventApplicationSynced;
  hasState(): boolean;
  clearState(): NotificationEventApplicationSynced;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventApplicationSynced.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventApplicationSynced): NotificationEventApplicationSynced.AsObject;
  static serializeBinaryToWriter(message: NotificationEventApplicationSynced, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventApplicationSynced;
  static deserializeBinaryFromReader(message: NotificationEventApplicationSynced, reader: jspb.BinaryReader): NotificationEventApplicationSynced;
}

export namespace NotificationEventApplicationSynced {
  export type AsObject = {
    application?: pkg_model_application_pb.Application.AsObject,
    state?: pkg_model_application_pb.ApplicationSyncState.AsObject,
  }
}

export class NotificationEventApplicationOutOfSync extends jspb.Message {
  getApplication(): pkg_model_application_pb.Application | undefined;
  setApplication(value?: pkg_model_application_pb.Application): NotificationEventApplicationOutOfSync;
  hasApplication(): boolean;
  clearApplication(): NotificationEventApplicationOutOfSync;

  getState(): pkg_model_application_pb.ApplicationSyncState | undefined;
  setState(value?: pkg_model_application_pb.ApplicationSyncState): NotificationEventApplicationOutOfSync;
  hasState(): boolean;
  clearState(): NotificationEventApplicationOutOfSync;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventApplicationOutOfSync.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventApplicationOutOfSync): NotificationEventApplicationOutOfSync.AsObject;
  static serializeBinaryToWriter(message: NotificationEventApplicationOutOfSync, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventApplicationOutOfSync;
  static deserializeBinaryFromReader(message: NotificationEventApplicationOutOfSync, reader: jspb.BinaryReader): NotificationEventApplicationOutOfSync;
}

export namespace NotificationEventApplicationOutOfSync {
  export type AsObject = {
    application?: pkg_model_application_pb.Application.AsObject,
    state?: pkg_model_application_pb.ApplicationSyncState.AsObject,
  }
}

export class NotificationEventPipedStarted extends jspb.Message {
  getId(): string;
  setId(value: string): NotificationEventPipedStarted;

  getName(): string;
  setName(value: string): NotificationEventPipedStarted;

  getVersion(): string;
  setVersion(value: string): NotificationEventPipedStarted;

  getProjectId(): string;
  setProjectId(value: string): NotificationEventPipedStarted;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventPipedStarted.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventPipedStarted): NotificationEventPipedStarted.AsObject;
  static serializeBinaryToWriter(message: NotificationEventPipedStarted, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventPipedStarted;
  static deserializeBinaryFromReader(message: NotificationEventPipedStarted, reader: jspb.BinaryReader): NotificationEventPipedStarted;
}

export namespace NotificationEventPipedStarted {
  export type AsObject = {
    id: string,
    name: string,
    version: string,
    projectId: string,
  }
}

export class NotificationEventPipedStopped extends jspb.Message {
  getId(): string;
  setId(value: string): NotificationEventPipedStopped;

  getName(): string;
  setName(value: string): NotificationEventPipedStopped;

  getVersion(): string;
  setVersion(value: string): NotificationEventPipedStopped;

  getProjectId(): string;
  setProjectId(value: string): NotificationEventPipedStopped;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): NotificationEventPipedStopped.AsObject;
  static toObject(includeInstance: boolean, msg: NotificationEventPipedStopped): NotificationEventPipedStopped.AsObject;
  static serializeBinaryToWriter(message: NotificationEventPipedStopped, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): NotificationEventPipedStopped;
  static deserializeBinaryFromReader(message: NotificationEventPipedStopped, reader: jspb.BinaryReader): NotificationEventPipedStopped;
}

export namespace NotificationEventPipedStopped {
  export type AsObject = {
    id: string,
    name: string,
    version: string,
    projectId: string,
  }
}

export enum NotificationEventType { 
  EVENT_DEPLOYMENT_TRIGGERED = 0,
  EVENT_DEPLOYMENT_PLANNED = 1,
  EVENT_DEPLOYMENT_APPROVED = 2,
  EVENT_DEPLOYMENT_ROLLING_BACK = 3,
  EVENT_DEPLOYMENT_SUCCEEDED = 4,
  EVENT_DEPLOYMENT_FAILED = 5,
  EVENT_DEPLOYMENT_CANCELLED = 6,
  EVENT_DEPLOYMENT_WAIT_APPROVAL = 7,
  EVENT_DEPLOYMENT_TRIGGER_FAILED = 8,
  EVENT_DEPLOYMENT_STARTED = 9,
  EVENT_APPLICATION_SYNCED = 100,
  EVENT_APPLICATION_OUT_OF_SYNC = 101,
  EVENT_APPLICATION_HEALTHY = 200,
  EVENT_PIPED_STARTED = 300,
  EVENT_PIPED_STOPPED = 301,
}
export enum NotificationEventGroup { 
  EVENT_NONE = 0,
  EVENT_DEPLOYMENT = 1,
  EVENT_APPLICATION_SYNC = 2,
  EVENT_APPLICATION_HEALTH = 3,
  EVENT_PIPED = 4,
}
