import * as jspb from 'google-protobuf'




export class ApplicationGitPath extends jspb.Message {
  getRepo(): ApplicationGitRepository | undefined;
  setRepo(value?: ApplicationGitRepository): ApplicationGitPath;
  hasRepo(): boolean;
  clearRepo(): ApplicationGitPath;

  getPath(): string;
  setPath(value: string): ApplicationGitPath;

  getConfigFilename(): string;
  setConfigFilename(value: string): ApplicationGitPath;

  getUrl(): string;
  setUrl(value: string): ApplicationGitPath;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationGitPath.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationGitPath): ApplicationGitPath.AsObject;
  static serializeBinaryToWriter(message: ApplicationGitPath, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationGitPath;
  static deserializeBinaryFromReader(message: ApplicationGitPath, reader: jspb.BinaryReader): ApplicationGitPath;
}

export namespace ApplicationGitPath {
  export type AsObject = {
    repo?: ApplicationGitRepository.AsObject,
    path: string,
    configFilename: string,
    url: string,
  }
}

export class ApplicationGitRepository extends jspb.Message {
  getId(): string;
  setId(value: string): ApplicationGitRepository;

  getRemote(): string;
  setRemote(value: string): ApplicationGitRepository;

  getBranch(): string;
  setBranch(value: string): ApplicationGitRepository;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationGitRepository.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationGitRepository): ApplicationGitRepository.AsObject;
  static serializeBinaryToWriter(message: ApplicationGitRepository, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationGitRepository;
  static deserializeBinaryFromReader(message: ApplicationGitRepository, reader: jspb.BinaryReader): ApplicationGitRepository;
}

export namespace ApplicationGitRepository {
  export type AsObject = {
    id: string,
    remote: string,
    branch: string,
  }
}

export class ApplicationInfo extends jspb.Message {
  getId(): string;
  setId(value: string): ApplicationInfo;

  getName(): string;
  setName(value: string): ApplicationInfo;

  getKind(): ApplicationKind;
  setKind(value: ApplicationKind): ApplicationInfo;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): ApplicationInfo;

  getRepoId(): string;
  setRepoId(value: string): ApplicationInfo;

  getPath(): string;
  setPath(value: string): ApplicationInfo;

  getConfigFilename(): string;
  setConfigFilename(value: string): ApplicationInfo;

  getPipedId(): string;
  setPipedId(value: string): ApplicationInfo;

  getDescription(): string;
  setDescription(value: string): ApplicationInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApplicationInfo.AsObject;
  static toObject(includeInstance: boolean, msg: ApplicationInfo): ApplicationInfo.AsObject;
  static serializeBinaryToWriter(message: ApplicationInfo, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApplicationInfo;
  static deserializeBinaryFromReader(message: ApplicationInfo, reader: jspb.BinaryReader): ApplicationInfo;
}

export namespace ApplicationInfo {
  export type AsObject = {
    id: string,
    name: string,
    kind: ApplicationKind,
    labelsMap: Array<[string, string]>,
    repoId: string,
    path: string,
    configFilename: string,
    pipedId: string,
    description: string,
  }
}

export class ArtifactVersion extends jspb.Message {
  getKind(): ArtifactVersion.Kind;
  setKind(value: ArtifactVersion.Kind): ArtifactVersion;

  getVersion(): string;
  setVersion(value: string): ArtifactVersion;

  getName(): string;
  setName(value: string): ArtifactVersion;

  getUrl(): string;
  setUrl(value: string): ArtifactVersion;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ArtifactVersion.AsObject;
  static toObject(includeInstance: boolean, msg: ArtifactVersion): ArtifactVersion.AsObject;
  static serializeBinaryToWriter(message: ArtifactVersion, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ArtifactVersion;
  static deserializeBinaryFromReader(message: ArtifactVersion, reader: jspb.BinaryReader): ArtifactVersion;
}

export namespace ArtifactVersion {
  export type AsObject = {
    kind: ArtifactVersion.Kind,
    version: string,
    name: string,
    url: string,
  }

  export enum Kind { 
    UNKNOWN = 0,
    CONTAINER_IMAGE = 1,
    S3_OBJECT = 2,
    GIT_SOURCE = 3,
    TERRAFORM_MODULE = 4,
  }
}

export enum ApplicationKind { 
  KUBERNETES = 0,
  TERRAFORM = 1,
  LAMBDA = 3,
  CLOUDRUN = 4,
  ECS = 5,
  APPLICATION = 6,
}
export enum RollbackKind { 
  ROLLBACK_KUBERNETES = 0,
  ROLLBACK_TERRAFORM = 1,
  ROLLBACK_LAMBDA = 3,
  ROLLBACK_CLOUDRUN = 4,
  ROLLBACK_ECS = 5,
  ROLLBACK_CUSTOM_SYNC = 15,
}
export enum ApplicationActiveStatus { 
  ENABLED = 0,
  DISABLED = 1,
  DELETED = 2,
}
export enum SyncStrategy { 
  AUTO = 0,
  QUICK_SYNC = 1,
  PIPELINE = 2,
}
