import * as jspb from 'google-protobuf'




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

export enum InsightMetricsKind { 
  DEPLOYMENT_FREQUENCY = 0,
  CHANGE_FAILURE_RATE = 1,
  MTTR = 2,
  LEAD_TIME = 3,
  APPLICATIONS_COUNT = 4,
}
export enum InsightResultType { 
  MATRIX = 0,
  VECTOR = 1,
}
export enum InsightResolution { 
  DAILY = 0,
  MONTHLY = 1,
}
export enum InsightApplicationCountLabelKey { 
  KIND = 0,
  ACTIVE_STATUS = 1,
}
