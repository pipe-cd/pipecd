import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';


export class Piped extends jspb.Message {
  getId(): string;
  setId(value: string): Piped;

  getName(): string;
  setName(value: string): Piped;

  getDesc(): string;
  setDesc(value: string): Piped;

  getKeyHash(): string;
  setKeyHash(value: string): Piped;

  getProjectId(): string;
  setProjectId(value: string): Piped;

  getVersion(): string;
  setVersion(value: string): Piped;

  getConfig(): string;
  setConfig(value: string): Piped;

  getStartedAt(): number;
  setStartedAt(value: number): Piped;

  getCloudProvidersList(): Array<Piped.CloudProvider>;
  setCloudProvidersList(value: Array<Piped.CloudProvider>): Piped;
  clearCloudProvidersList(): Piped;
  addCloudProviders(value?: Piped.CloudProvider, index?: number): Piped.CloudProvider;

  getRepositoriesList(): Array<pkg_model_common_pb.ApplicationGitRepository>;
  setRepositoriesList(value: Array<pkg_model_common_pb.ApplicationGitRepository>): Piped;
  clearRepositoriesList(): Piped;
  addRepositories(value?: pkg_model_common_pb.ApplicationGitRepository, index?: number): pkg_model_common_pb.ApplicationGitRepository;

  getStatus(): Piped.ConnectionStatus;
  setStatus(value: Piped.ConnectionStatus): Piped;

  getSecretEncryption(): Piped.SecretEncryption | undefined;
  setSecretEncryption(value?: Piped.SecretEncryption): Piped;
  hasSecretEncryption(): boolean;
  clearSecretEncryption(): Piped;

  getKeysList(): Array<PipedKey>;
  setKeysList(value: Array<PipedKey>): Piped;
  clearKeysList(): Piped;
  addKeys(value?: PipedKey, index?: number): PipedKey;

  getDesiredVersion(): string;
  setDesiredVersion(value: string): Piped;

  getDisabled(): boolean;
  setDisabled(value: boolean): Piped;

  getCreatedAt(): number;
  setCreatedAt(value: number): Piped;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Piped;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Piped.AsObject;
  static toObject(includeInstance: boolean, msg: Piped): Piped.AsObject;
  static serializeBinaryToWriter(message: Piped, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Piped;
  static deserializeBinaryFromReader(message: Piped, reader: jspb.BinaryReader): Piped;
}

export namespace Piped {
  export type AsObject = {
    id: string,
    name: string,
    desc: string,
    keyHash: string,
    projectId: string,
    version: string,
    config: string,
    startedAt: number,
    cloudProvidersList: Array<Piped.CloudProvider.AsObject>,
    repositoriesList: Array<pkg_model_common_pb.ApplicationGitRepository.AsObject>,
    status: Piped.ConnectionStatus,
    secretEncryption?: Piped.SecretEncryption.AsObject,
    keysList: Array<PipedKey.AsObject>,
    desiredVersion: string,
    disabled: boolean,
    createdAt: number,
    updatedAt: number,
  }

  export class CloudProvider extends jspb.Message {
    getName(): string;
    setName(value: string): CloudProvider;

    getType(): string;
    setType(value: string): CloudProvider;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CloudProvider.AsObject;
    static toObject(includeInstance: boolean, msg: CloudProvider): CloudProvider.AsObject;
    static serializeBinaryToWriter(message: CloudProvider, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CloudProvider;
    static deserializeBinaryFromReader(message: CloudProvider, reader: jspb.BinaryReader): CloudProvider;
  }

  export namespace CloudProvider {
    export type AsObject = {
      name: string,
      type: string,
    }
  }


  export class SecretEncryption extends jspb.Message {
    getType(): string;
    setType(value: string): SecretEncryption;

    getPublicKey(): string;
    setPublicKey(value: string): SecretEncryption;

    getEncryptServiceAccount(): string;
    setEncryptServiceAccount(value: string): SecretEncryption;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SecretEncryption.AsObject;
    static toObject(includeInstance: boolean, msg: SecretEncryption): SecretEncryption.AsObject;
    static serializeBinaryToWriter(message: SecretEncryption, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SecretEncryption;
    static deserializeBinaryFromReader(message: SecretEncryption, reader: jspb.BinaryReader): SecretEncryption;
  }

  export namespace SecretEncryption {
    export type AsObject = {
      type: string,
      publicKey: string,
      encryptServiceAccount: string,
    }
  }


  export enum ConnectionStatus { 
    UNKNOWN = 0,
    ONLINE = 1,
    OFFLINE = 2,
  }
}

export class PipedKey extends jspb.Message {
  getHash(): string;
  setHash(value: string): PipedKey;

  getCreator(): string;
  setCreator(value: string): PipedKey;

  getCreatedAt(): number;
  setCreatedAt(value: number): PipedKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PipedKey.AsObject;
  static toObject(includeInstance: boolean, msg: PipedKey): PipedKey.AsObject;
  static serializeBinaryToWriter(message: PipedKey, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PipedKey;
  static deserializeBinaryFromReader(message: PipedKey, reader: jspb.BinaryReader): PipedKey;
}

export namespace PipedKey {
  export type AsObject = {
    hash: string,
    creator: string,
    createdAt: number,
  }
}

