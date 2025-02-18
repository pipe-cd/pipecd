import * as jspb from 'google-protobuf'




export class DeploymentTrace extends jspb.Message {
  getId(): string;
  setId(value: string): DeploymentTrace;

  getTitle(): string;
  setTitle(value: string): DeploymentTrace;

  getCommitHash(): string;
  setCommitHash(value: string): DeploymentTrace;

  getCommitUrl(): string;
  setCommitUrl(value: string): DeploymentTrace;

  getCommitMessage(): string;
  setCommitMessage(value: string): DeploymentTrace;

  getCommitTimestamp(): number;
  setCommitTimestamp(value: number): DeploymentTrace;

  getAuthor(): string;
  setAuthor(value: string): DeploymentTrace;

  getCompletedAt(): number;
  setCompletedAt(value: number): DeploymentTrace;

  getCreatedAt(): number;
  setCreatedAt(value: number): DeploymentTrace;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): DeploymentTrace;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentTrace.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentTrace): DeploymentTrace.AsObject;
  static serializeBinaryToWriter(message: DeploymentTrace, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentTrace;
  static deserializeBinaryFromReader(message: DeploymentTrace, reader: jspb.BinaryReader): DeploymentTrace;
}

export namespace DeploymentTrace {
  export type AsObject = {
    id: string,
    title: string,
    commitHash: string,
    commitUrl: string,
    commitMessage: string,
    commitTimestamp: number,
    author: string,
    completedAt: number,
    createdAt: number,
    updatedAt: number,
  }
}

