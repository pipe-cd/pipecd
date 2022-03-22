import * as jspb from 'google-protobuf'




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

export class DeploymentSubset extends jspb.Message {
  getCreatedAt(): number;
  setCreatedAt(value: number): DeploymentSubset;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): DeploymentSubset;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentSubset.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentSubset): DeploymentSubset.AsObject;
  static serializeBinaryToWriter(message: DeploymentSubset, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentSubset;
  static deserializeBinaryFromReader(message: DeploymentSubset, reader: jspb.BinaryReader): DeploymentSubset;
}

export namespace DeploymentSubset {
  export type AsObject = {
    createdAt: number,
    updatedAt: number,
  }
}

export class DailyDeployment extends jspb.Message {
  getDate(): number;
  setDate(value: number): DailyDeployment;

  getCreatedAt(): number;
  setCreatedAt(value: number): DailyDeployment;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): DailyDeployment;

  getDailyDeploymentsList(): Array<DeploymentSubset>;
  setDailyDeploymentsList(value: Array<DeploymentSubset>): DailyDeployment;
  clearDailyDeploymentsList(): DailyDeployment;
  addDailyDeployments(value?: DeploymentSubset, index?: number): DeploymentSubset;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DailyDeployment.AsObject;
  static toObject(includeInstance: boolean, msg: DailyDeployment): DailyDeployment.AsObject;
  static serializeBinaryToWriter(message: DailyDeployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DailyDeployment;
  static deserializeBinaryFromReader(message: DailyDeployment, reader: jspb.BinaryReader): DailyDeployment;
}

export namespace DailyDeployment {
  export type AsObject = {
    date: number,
    createdAt: number,
    updatedAt: number,
    dailyDeploymentsList: Array<DeploymentSubset.AsObject>,
  }
}

export class DeploymentChunk extends jspb.Message {
  getDateRange(): ChunkDateRange | undefined;
  setDateRange(value?: ChunkDateRange): DeploymentChunk;
  hasDateRange(): boolean;
  clearDateRange(): DeploymentChunk;

  getDeploymentsList(): Array<DailyDeployment>;
  setDeploymentsList(value: Array<DailyDeployment>): DeploymentChunk;
  clearDeploymentsList(): DeploymentChunk;
  addDeployments(value?: DailyDeployment, index?: number): DailyDeployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentChunk.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentChunk): DeploymentChunk.AsObject;
  static serializeBinaryToWriter(message: DeploymentChunk, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentChunk;
  static deserializeBinaryFromReader(message: DeploymentChunk, reader: jspb.BinaryReader): DeploymentChunk;
}

export namespace DeploymentChunk {
  export type AsObject = {
    dateRange?: ChunkDateRange.AsObject,
    deploymentsList: Array<DailyDeployment.AsObject>,
  }
}

export class DeploymentChunkMetaData extends jspb.Message {
  getDataList(): Array<DeploymentChunkMetaData.ChunkData>;
  setDataList(value: Array<DeploymentChunkMetaData.ChunkData>): DeploymentChunkMetaData;
  clearDataList(): DeploymentChunkMetaData;
  addData(value?: DeploymentChunkMetaData.ChunkData, index?: number): DeploymentChunkMetaData.ChunkData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentChunkMetaData.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentChunkMetaData): DeploymentChunkMetaData.AsObject;
  static serializeBinaryToWriter(message: DeploymentChunkMetaData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentChunkMetaData;
  static deserializeBinaryFromReader(message: DeploymentChunkMetaData, reader: jspb.BinaryReader): DeploymentChunkMetaData;
}

export namespace DeploymentChunkMetaData {
  export type AsObject = {
    dataList: Array<DeploymentChunkMetaData.ChunkData.AsObject>,
  }

  export class ChunkData extends jspb.Message {
    getDateRange(): ChunkDateRange | undefined;
    setDateRange(value?: ChunkDateRange): ChunkData;
    hasDateRange(): boolean;
    clearDateRange(): ChunkData;

    getChunkKey(): string;
    setChunkKey(value: string): ChunkData;

    getChunkSize(): number;
    setChunkSize(value: number): ChunkData;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChunkData.AsObject;
    static toObject(includeInstance: boolean, msg: ChunkData): ChunkData.AsObject;
    static serializeBinaryToWriter(message: ChunkData, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChunkData;
    static deserializeBinaryFromReader(message: ChunkData, reader: jspb.BinaryReader): ChunkData;
  }

  export namespace ChunkData {
    export type AsObject = {
      dateRange?: ChunkDateRange.AsObject,
      chunkKey: string,
      chunkSize: number,
    }
  }

}

export class ChunkDateRange extends jspb.Message {
  getFrom(): number;
  setFrom(value: number): ChunkDateRange;

  getTo(): number;
  setTo(value: number): ChunkDateRange;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChunkDateRange.AsObject;
  static toObject(includeInstance: boolean, msg: ChunkDateRange): ChunkDateRange.AsObject;
  static serializeBinaryToWriter(message: ChunkDateRange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChunkDateRange;
  static deserializeBinaryFromReader(message: ChunkDateRange, reader: jspb.BinaryReader): ChunkDateRange;
}

export namespace ChunkDateRange {
  export type AsObject = {
    from: number,
    to: number,
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
