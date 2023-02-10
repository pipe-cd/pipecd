import * as jspb from 'google-protobuf'




export class LogBlock extends jspb.Message {
  getIndex(): number;
  setIndex(value: number): LogBlock;

  getLog(): string;
  setLog(value: string): LogBlock;

  getSeverity(): LogSeverity;
  setSeverity(value: LogSeverity): LogBlock;

  getCreatedAt(): number;
  setCreatedAt(value: number): LogBlock;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogBlock.AsObject;
  static toObject(includeInstance: boolean, msg: LogBlock): LogBlock.AsObject;
  static serializeBinaryToWriter(message: LogBlock, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogBlock;
  static deserializeBinaryFromReader(message: LogBlock, reader: jspb.BinaryReader): LogBlock;
}

export namespace LogBlock {
  export type AsObject = {
    index: number,
    log: string,
    severity: LogSeverity,
    createdAt: number,
  }
}

export enum LogSeverity { 
  INFO = 0,
  SUCCESS = 1,
  ERROR = 2,
}
