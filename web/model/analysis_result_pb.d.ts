import * as jspb from 'google-protobuf'




export class AnalysisResult extends jspb.Message {
  getStartTime(): number;
  setStartTime(value: number): AnalysisResult;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AnalysisResult.AsObject;
  static toObject(includeInstance: boolean, msg: AnalysisResult): AnalysisResult.AsObject;
  static serializeBinaryToWriter(message: AnalysisResult, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AnalysisResult;
  static deserializeBinaryFromReader(message: AnalysisResult, reader: jspb.BinaryReader): AnalysisResult;
}

export namespace AnalysisResult {
  export type AsObject = {
    startTime: number,
  }
}

