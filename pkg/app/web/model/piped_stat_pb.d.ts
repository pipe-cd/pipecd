import * as jspb from 'google-protobuf'




export class PipedStat extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): PipedStat;

  getMetrics(): Uint8Array | string;
  getMetrics_asU8(): Uint8Array;
  getMetrics_asB64(): string;
  setMetrics(value: Uint8Array | string): PipedStat;

  getTimestamp(): number;
  setTimestamp(value: number): PipedStat;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PipedStat.AsObject;
  static toObject(includeInstance: boolean, msg: PipedStat): PipedStat.AsObject;
  static serializeBinaryToWriter(message: PipedStat, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PipedStat;
  static deserializeBinaryFromReader(message: PipedStat, reader: jspb.BinaryReader): PipedStat;
}

export namespace PipedStat {
  export type AsObject = {
    pipedId: string,
    metrics: Uint8Array | string,
    timestamp: number,
  }
}

