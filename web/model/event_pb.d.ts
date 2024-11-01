import * as jspb from 'google-protobuf'




export class Event extends jspb.Message {
  getId(): string;
  setId(value: string): Event;

  getName(): string;
  setName(value: string): Event;

  getData(): string;
  setData(value: string): Event;

  getProjectId(): string;
  setProjectId(value: string): Event;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): Event;

  getEventKey(): string;
  setEventKey(value: string): Event;

  getStatus(): EventStatus;
  setStatus(value: EventStatus): Event;

  getStatusDescription(): string;
  setStatusDescription(value: string): Event;

  getContextsMap(): jspb.Map<string, string>;
  clearContextsMap(): Event;

  getHandledAt(): number;
  setHandledAt(value: number): Event;

  getCreatedAt(): number;
  setCreatedAt(value: number): Event;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Event;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Event.AsObject;
  static toObject(includeInstance: boolean, msg: Event): Event.AsObject;
  static serializeBinaryToWriter(message: Event, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Event;
  static deserializeBinaryFromReader(message: Event, reader: jspb.BinaryReader): Event;
}

export namespace Event {
  export type AsObject = {
    id: string,
    name: string,
    data: string,
    projectId: string,
    labelsMap: Array<[string, string]>,
    eventKey: string,
    status: EventStatus,
    statusDescription: string,
    contextsMap: Array<[string, string]>,
    handledAt: number,
    createdAt: number,
    updatedAt: number,
  }
}

export enum EventStatus { 
  EVENT_NOT_HANDLED = 0,
  EVENT_SUCCESS = 1,
  EVENT_FAILURE = 2,
  EVENT_OUTDATED = 3,
}
