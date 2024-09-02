import * as jspb from 'google-protobuf'



export class DeploymentSource extends jspb.Message {
  getApplicationDirectory(): string;
  setApplicationDirectory(value: string): DeploymentSource;

  getRevision(): string;
  setRevision(value: string): DeploymentSource;

  getApplicationConfig(): PluginApplicationSpec | undefined;
  setApplicationConfig(value?: PluginApplicationSpec): DeploymentSource;
  hasApplicationConfig(): boolean;
  clearApplicationConfig(): DeploymentSource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentSource.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentSource): DeploymentSource.AsObject;
  static serializeBinaryToWriter(message: DeploymentSource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentSource;
  static deserializeBinaryFromReader(message: DeploymentSource, reader: jspb.BinaryReader): DeploymentSource;
}

export namespace DeploymentSource {
  export type AsObject = {
    applicationDirectory: string,
    revision: string,
    applicationConfig?: PluginApplicationSpec.AsObject,
  }
}

export class PluginApplicationSpec extends jspb.Message {
  getKind(): string;
  setKind(value: string): PluginApplicationSpec;

  getApiVersion(): string;
  setApiVersion(value: string): PluginApplicationSpec;

  getSpec(): Uint8Array | string;
  getSpec_asU8(): Uint8Array;
  getSpec_asB64(): string;
  setSpec(value: Uint8Array | string): PluginApplicationSpec;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PluginApplicationSpec.AsObject;
  static toObject(includeInstance: boolean, msg: PluginApplicationSpec): PluginApplicationSpec.AsObject;
  static serializeBinaryToWriter(message: PluginApplicationSpec, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PluginApplicationSpec;
  static deserializeBinaryFromReader(message: PluginApplicationSpec, reader: jspb.BinaryReader): PluginApplicationSpec;
}

export namespace PluginApplicationSpec {
  export type AsObject = {
    kind: string,
    apiVersion: string,
    spec: Uint8Array | string,
  }
}

