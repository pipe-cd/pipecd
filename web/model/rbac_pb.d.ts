import * as jspb from 'google-protobuf'

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as pkg_model_project_pb from 'pipecd/web/model/project_pb';


export class RBAC extends jspb.Message {
  getResource(): pkg_model_project_pb.ProjectRBACResource.ResourceType;
  setResource(value: pkg_model_project_pb.ProjectRBACResource.ResourceType): RBAC;

  getAction(): pkg_model_project_pb.ProjectRBACPolicy.Action;
  setAction(value: pkg_model_project_pb.ProjectRBACPolicy.Action): RBAC;

  getIgnored(): boolean;
  setIgnored(value: boolean): RBAC;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RBAC.AsObject;
  static toObject(includeInstance: boolean, msg: RBAC): RBAC.AsObject;
  static serializeBinaryToWriter(message: RBAC, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RBAC;
  static deserializeBinaryFromReader(message: RBAC, reader: jspb.BinaryReader): RBAC;
}

export namespace RBAC {
  export type AsObject = {
    resource: pkg_model_project_pb.ProjectRBACResource.ResourceType,
    action: pkg_model_project_pb.ProjectRBACPolicy.Action,
    ignored: boolean,
  }
}

