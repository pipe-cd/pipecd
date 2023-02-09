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

export class LogBlocks extends jspb.Message {
  getLogblocksList(): Array<LogBlock>;
  setLogblocksList(value: Array<LogBlock>): LogBlocks;
  clearLogblocksList(): LogBlocks;
  addLogblocks(value?: LogBlock, index?: number): LogBlock;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): LogBlocks.AsObject;
  static toObject(includeInstance: boolean, msg: LogBlocks): LogBlocks.AsObject;
  static serializeBinaryToWriter(message: LogBlocks, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): LogBlocks;
  static deserializeBinaryFromReader(message: LogBlocks, reader: jspb.BinaryReader): LogBlocks;
}

export namespace LogBlocks {
  export type AsObject = {
    logblocksList: Array<LogBlock.AsObject>,
  }
}

export enum LogSeverity { 
  INFO = 0,
  SUCCESS = 1,
  ERROR = 2,
}
