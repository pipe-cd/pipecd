import * as jspb from 'google-protobuf'


import * as pkg_model_deployment_pb from 'pipecd/web/model/deployment_pb';


export class InsightSample extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): InsightSample;

  getDataPoint(): InsightDataPoint | undefined;
  setDataPoint(value?: InsightDataPoint): InsightSample;
  hasDataPoint(): boolean;
  clearDataPoint(): InsightSample;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightSample.AsObject;
  static toObject(includeInstance: boolean, msg: InsightSample): InsightSample.AsObject;
  static serializeBinaryToWriter(message: InsightSample, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightSample;
  static deserializeBinaryFromReader(message: InsightSample, reader: jspb.BinaryReader): InsightSample;
}

export namespace InsightSample {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
    dataPoint?: InsightDataPoint.AsObject,
  }
}

export class InsightSampleStream extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): InsightSampleStream;

  getDataPointsList(): Array<InsightDataPoint>;
  setDataPointsList(value: Array<InsightDataPoint>): InsightSampleStream;
  clearDataPointsList(): InsightSampleStream;
  addDataPoints(value?: InsightDataPoint, index?: number): InsightDataPoint;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightSampleStream.AsObject;
  static toObject(includeInstance: boolean, msg: InsightSampleStream): InsightSampleStream.AsObject;
  static serializeBinaryToWriter(message: InsightSampleStream, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightSampleStream;
  static deserializeBinaryFromReader(message: InsightSampleStream, reader: jspb.BinaryReader): InsightSampleStream;
}

export namespace InsightSampleStream {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
    dataPointsList: Array<InsightDataPoint.AsObject>,
  }
}

export class InsightDataPoint extends jspb.Message {
  getTimestamp(): number;
  setTimestamp(value: number): InsightDataPoint;

  getValue(): number;
  setValue(value: number): InsightDataPoint;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDataPoint.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDataPoint): InsightDataPoint.AsObject;
  static serializeBinaryToWriter(message: InsightDataPoint, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDataPoint;
  static deserializeBinaryFromReader(message: InsightDataPoint, reader: jspb.BinaryReader): InsightDataPoint;
}

export namespace InsightDataPoint {
  export type AsObject = {
    timestamp: number,
    value: number,
  }
}

export class InsightApplicationCount extends jspb.Message {
  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): InsightApplicationCount;

  getCount(): number;
  setCount(value: number): InsightApplicationCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightApplicationCount.AsObject;
  static toObject(includeInstance: boolean, msg: InsightApplicationCount): InsightApplicationCount.AsObject;
  static serializeBinaryToWriter(message: InsightApplicationCount, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightApplicationCount;
  static deserializeBinaryFromReader(message: InsightApplicationCount, reader: jspb.BinaryReader): InsightApplicationCount;
}

export namespace InsightApplicationCount {
  export type AsObject = {
    labelsMap: Array<[string, string]>,
    count: number,
  }
}

export class InsightDeploymentChunk extends jspb.Message {
  getVersion(): InsightDeploymentVersion;
  setVersion(value: InsightDeploymentVersion): InsightDeploymentChunk;

  getFrom(): number;
  setFrom(value: number): InsightDeploymentChunk;

  getTo(): number;
  setTo(value: number): InsightDeploymentChunk;

  getDeploymentsList(): Array<InsightDeployment>;
  setDeploymentsList(value: Array<InsightDeployment>): InsightDeploymentChunk;
  clearDeploymentsList(): InsightDeploymentChunk;
  addDeployments(value?: InsightDeployment, index?: number): InsightDeployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeploymentChunk.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeploymentChunk): InsightDeploymentChunk.AsObject;
  static serializeBinaryToWriter(message: InsightDeploymentChunk, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeploymentChunk;
  static deserializeBinaryFromReader(message: InsightDeploymentChunk, reader: jspb.BinaryReader): InsightDeploymentChunk;
}

export namespace InsightDeploymentChunk {
  export type AsObject = {
    version: InsightDeploymentVersion,
    from: number,
    to: number,
    deploymentsList: Array<InsightDeployment.AsObject>,
  }
}

export class InsightDeployment extends jspb.Message {
  getId(): string;
  setId(value: string): InsightDeployment;

  getAppId(): string;
  setAppId(value: string): InsightDeployment;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): InsightDeployment;

  getStartedAt(): number;
  setStartedAt(value: number): InsightDeployment;

  getCompletedAt(): number;
  setCompletedAt(value: number): InsightDeployment;

  getRollbackStartedAt(): number;
  setRollbackStartedAt(value: number): InsightDeployment;

  getCompleteStatus(): pkg_model_deployment_pb.DeploymentStatus;
  setCompleteStatus(value: pkg_model_deployment_pb.DeploymentStatus): InsightDeployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeployment.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeployment): InsightDeployment.AsObject;
  static serializeBinaryToWriter(message: InsightDeployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeployment;
  static deserializeBinaryFromReader(message: InsightDeployment, reader: jspb.BinaryReader): InsightDeployment;
}

export namespace InsightDeployment {
  export type AsObject = {
    id: string,
    appId: string,
    labelsMap: Array<[string, string]>,
    startedAt: number,
    completedAt: number,
    rollbackStartedAt: number,
    completeStatus: pkg_model_deployment_pb.DeploymentStatus,
  }
}

export class InsightDeploymentChunkMetadata extends jspb.Message {
  getChunksList(): Array<InsightDeploymentChunkMetadata.ChunkMeta>;
  setChunksList(value: Array<InsightDeploymentChunkMetadata.ChunkMeta>): InsightDeploymentChunkMetadata;
  clearChunksList(): InsightDeploymentChunkMetadata;
  addChunks(value?: InsightDeploymentChunkMetadata.ChunkMeta, index?: number): InsightDeploymentChunkMetadata.ChunkMeta;

  getCreatedAt(): number;
  setCreatedAt(value: number): InsightDeploymentChunkMetadata;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): InsightDeploymentChunkMetadata;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeploymentChunkMetadata.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeploymentChunkMetadata): InsightDeploymentChunkMetadata.AsObject;
  static serializeBinaryToWriter(message: InsightDeploymentChunkMetadata, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeploymentChunkMetadata;
  static deserializeBinaryFromReader(message: InsightDeploymentChunkMetadata, reader: jspb.BinaryReader): InsightDeploymentChunkMetadata;
}

export namespace InsightDeploymentChunkMetadata {
  export type AsObject = {
    chunksList: Array<InsightDeploymentChunkMetadata.ChunkMeta.AsObject>,
    createdAt: number,
    updatedAt: number,
  }

  export class ChunkMeta extends jspb.Message {
    getFrom(): number;
    setFrom(value: number): ChunkMeta;

    getTo(): number;
    setTo(value: number): ChunkMeta;

    getName(): string;
    setName(value: string): ChunkMeta;

    getSize(): number;
    setSize(value: number): ChunkMeta;

    getCount(): number;
    setCount(value: number): ChunkMeta;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChunkMeta.AsObject;
    static toObject(includeInstance: boolean, msg: ChunkMeta): ChunkMeta.AsObject;
    static serializeBinaryToWriter(message: ChunkMeta, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChunkMeta;
    static deserializeBinaryFromReader(message: ChunkMeta, reader: jspb.BinaryReader): ChunkMeta;
  }

  export namespace ChunkMeta {
    export type AsObject = {
      from: number,
      to: number,
      name: string,
      size: number,
      count: number,
    }
  }

}

export enum InsightResultType { 
  MATRIX = 0,
  VECTOR = 1,
}
export enum InsightMetricsKind { 
  DEPLOYMENT_FREQUENCY = 0,
  CHANGE_FAILURE_RATE = 1,
  MTTR = 2,
  LEAD_TIME = 3,
  APPLICATIONS_COUNT = 4,
}
export enum InsightApplicationCountLabelKey { 
  KIND = 0,
  ACTIVE_STATUS = 1,
}
export enum InsightDeploymentVersion { 
  V0 = 0,
}
