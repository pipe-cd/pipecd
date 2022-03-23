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

export class InsightDeploymentSubset extends jspb.Message {
  getId(): string;
  setId(value: string): InsightDeploymentSubset;

  getCreatedAt(): number;
  setCreatedAt(value: number): InsightDeploymentSubset;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): InsightDeploymentSubset;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeploymentSubset.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeploymentSubset): InsightDeploymentSubset.AsObject;
  static serializeBinaryToWriter(message: InsightDeploymentSubset, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeploymentSubset;
  static deserializeBinaryFromReader(message: InsightDeploymentSubset, reader: jspb.BinaryReader): InsightDeploymentSubset;
}

export namespace InsightDeploymentSubset {
  export type AsObject = {
    id: string,
    createdAt: number,
    updatedAt: number,
  }
}

export class InsightDailyDeployment extends jspb.Message {
  getDate(): number;
  setDate(value: number): InsightDailyDeployment;

  getCreatedAt(): number;
  setCreatedAt(value: number): InsightDailyDeployment;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): InsightDailyDeployment;

  getDailyDeploymentsList(): Array<InsightDeploymentSubset>;
  setDailyDeploymentsList(value: Array<InsightDeploymentSubset>): InsightDailyDeployment;
  clearDailyDeploymentsList(): InsightDailyDeployment;
  addDailyDeployments(value?: InsightDeploymentSubset, index?: number): InsightDeploymentSubset;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDailyDeployment.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDailyDeployment): InsightDailyDeployment.AsObject;
  static serializeBinaryToWriter(message: InsightDailyDeployment, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDailyDeployment;
  static deserializeBinaryFromReader(message: InsightDailyDeployment, reader: jspb.BinaryReader): InsightDailyDeployment;
}

export namespace InsightDailyDeployment {
  export type AsObject = {
    date: number,
    createdAt: number,
    updatedAt: number,
    dailyDeploymentsList: Array<InsightDeploymentSubset.AsObject>,
  }
}

export class InsightDeploymentChunk extends jspb.Message {
  getDateRange(): InsightChunkDateRange | undefined;
  setDateRange(value?: InsightChunkDateRange): InsightDeploymentChunk;
  hasDateRange(): boolean;
  clearDateRange(): InsightDeploymentChunk;

  getDeploymentsList(): Array<InsightDailyDeployment>;
  setDeploymentsList(value: Array<InsightDailyDeployment>): InsightDeploymentChunk;
  clearDeploymentsList(): InsightDeploymentChunk;
  addDeployments(value?: InsightDailyDeployment, index?: number): InsightDailyDeployment;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeploymentChunk.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeploymentChunk): InsightDeploymentChunk.AsObject;
  static serializeBinaryToWriter(message: InsightDeploymentChunk, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeploymentChunk;
  static deserializeBinaryFromReader(message: InsightDeploymentChunk, reader: jspb.BinaryReader): InsightDeploymentChunk;
}

export namespace InsightDeploymentChunk {
  export type AsObject = {
    dateRange?: InsightChunkDateRange.AsObject,
    deploymentsList: Array<InsightDailyDeployment.AsObject>,
  }
}

export class InsightDeploymentChunkMetaData extends jspb.Message {
  getDataList(): Array<InsightDeploymentChunkMetaData.InsightChunkData>;
  setDataList(value: Array<InsightDeploymentChunkMetaData.InsightChunkData>): InsightDeploymentChunkMetaData;
  clearDataList(): InsightDeploymentChunkMetaData;
  addData(value?: InsightDeploymentChunkMetaData.InsightChunkData, index?: number): InsightDeploymentChunkMetaData.InsightChunkData;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightDeploymentChunkMetaData.AsObject;
  static toObject(includeInstance: boolean, msg: InsightDeploymentChunkMetaData): InsightDeploymentChunkMetaData.AsObject;
  static serializeBinaryToWriter(message: InsightDeploymentChunkMetaData, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightDeploymentChunkMetaData;
  static deserializeBinaryFromReader(message: InsightDeploymentChunkMetaData, reader: jspb.BinaryReader): InsightDeploymentChunkMetaData;
}

export namespace InsightDeploymentChunkMetaData {
  export type AsObject = {
    dataList: Array<InsightDeploymentChunkMetaData.InsightChunkData.AsObject>,
  }

  export class InsightChunkData extends jspb.Message {
    getDateRange(): InsightChunkDateRange | undefined;
    setDateRange(value?: InsightChunkDateRange): InsightChunkData;
    hasDateRange(): boolean;
    clearDateRange(): InsightChunkData;

    getChunkKey(): string;
    setChunkKey(value: string): InsightChunkData;

    getChunkSize(): number;
    setChunkSize(value: number): InsightChunkData;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InsightChunkData.AsObject;
    static toObject(includeInstance: boolean, msg: InsightChunkData): InsightChunkData.AsObject;
    static serializeBinaryToWriter(message: InsightChunkData, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InsightChunkData;
    static deserializeBinaryFromReader(message: InsightChunkData, reader: jspb.BinaryReader): InsightChunkData;
  }

  export namespace InsightChunkData {
    export type AsObject = {
      dateRange?: InsightChunkDateRange.AsObject,
      chunkKey: string,
      chunkSize: number,
    }
  }

}

export class InsightChunkDateRange extends jspb.Message {
  getFrom(): number;
  setFrom(value: number): InsightChunkDateRange;

  getTo(): number;
  setTo(value: number): InsightChunkDateRange;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): InsightChunkDateRange.AsObject;
  static toObject(includeInstance: boolean, msg: InsightChunkDateRange): InsightChunkDateRange.AsObject;
  static serializeBinaryToWriter(message: InsightChunkDateRange, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): InsightChunkDateRange;
  static deserializeBinaryFromReader(message: InsightChunkDateRange, reader: jspb.BinaryReader): InsightChunkDateRange;
}

export namespace InsightChunkDateRange {
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
