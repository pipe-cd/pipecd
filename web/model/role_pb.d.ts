import * as jspb from 'google-protobuf'

import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';
import * as pkg_model_project_pb from 'pipecd/web/model/project_pb';


export class Role extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): Role;

  getProjectRole(): Role.ProjectRole;
  setProjectRole(value: Role.ProjectRole): Role;

  getProjectPoliciesList(): Array<pkg_model_project_pb.ProjectRBACPolicy>;
  setProjectPoliciesList(value: Array<pkg_model_project_pb.ProjectRBACPolicy>): Role;
  clearProjectPoliciesList(): Role;
  addProjectPolicies(value?: pkg_model_project_pb.ProjectRBACPolicy, index?: number): pkg_model_project_pb.ProjectRBACPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Role.AsObject;
  static toObject(includeInstance: boolean, msg: Role): Role.AsObject;
  static serializeBinaryToWriter(message: Role, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Role;
  static deserializeBinaryFromReader(message: Role, reader: jspb.BinaryReader): Role;
}

export namespace Role {
  export type AsObject = {
    projectId: string,
    projectRole: Role.ProjectRole,
    projectPoliciesList: Array<pkg_model_project_pb.ProjectRBACPolicy.AsObject>,
  }

  export enum ProjectRole { 
    VIEWER = 0,
    EDITOR = 1,
    ADMIN = 2,
  }
}

