import * as jspb from 'google-protobuf'



export class DeploymentSource extends jspb.Message {
  getApplicationDirectory(): string;
  setApplicationDirectory(value: string): DeploymentSource;

  getRevision(): string;
  setRevision(value: string): DeploymentSource;

  getApplicationConfig(): Uint8Array | string;
  getApplicationConfig_asU8(): Uint8Array;
  getApplicationConfig_asB64(): string;
  setApplicationConfig(value: Uint8Array | string): DeploymentSource;

  getApplicationConfigFilename(): string;
  setApplicationConfigFilename(value: string): DeploymentSource;

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
    applicationConfig: Uint8Array | string,
    applicationConfigFilename: string,
  }
}

