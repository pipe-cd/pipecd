import * as jspb from 'google-protobuf'


import * as pkg_model_role_pb from 'pipecd/web/model/role_pb';


export class User extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): User;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): User;

  getRole(): pkg_model_role_pb.Role | undefined;
  setRole(value?: pkg_model_role_pb.Role): User;
  hasRole(): boolean;
  clearRole(): User;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): User.AsObject;
  static toObject(includeInstance: boolean, msg: User): User.AsObject;
  static serializeBinaryToWriter(message: User, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): User;
  static deserializeBinaryFromReader(message: User, reader: jspb.BinaryReader): User;
}

export namespace User {
  export type AsObject = {
    username: string,
    avatarUrl: string,
    role?: pkg_model_role_pb.Role.AsObject,
  }
}

