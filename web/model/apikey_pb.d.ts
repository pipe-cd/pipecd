import * as jspb from 'google-protobuf'




export class APIKey extends jspb.Message {
  getId(): string;
  setId(value: string): APIKey;

  getName(): string;
  setName(value: string): APIKey;

  getKeyHash(): string;
  setKeyHash(value: string): APIKey;

  getProjectId(): string;
  setProjectId(value: string): APIKey;

  getRole(): APIKey.Role;
  setRole(value: APIKey.Role): APIKey;

  getCreator(): string;
  setCreator(value: string): APIKey;

  getLastUsedAt(): number;
  setLastUsedAt(value: number): APIKey;

  getDisabled(): boolean;
  setDisabled(value: boolean): APIKey;

  getCreatedAt(): number;
  setCreatedAt(value: number): APIKey;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): APIKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): APIKey.AsObject;
  static toObject(includeInstance: boolean, msg: APIKey): APIKey.AsObject;
  static serializeBinaryToWriter(message: APIKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): APIKey;
  static deserializeBinaryFromReader(message: APIKey, reader: jspb.BinaryReader): APIKey;
}

export namespace APIKey {
  export type AsObject = {
    id: string,
    name: string,
    keyHash: string,
    projectId: string,
    role: APIKey.Role,
    creator: string,
    lastUsedAt: number,
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
  }

  export enum Role { 
    READ_ONLY = 0,
    READ_WRITE = 1,
  }
}

