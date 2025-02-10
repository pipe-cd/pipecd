import * as jspb from 'google-protobuf'


import * as pkg_model_common_pb from 'pipecd/web/model/common_pb';
import * as pkg_model_insight_pb from 'pipecd/web/model/insight_pb';
import * as pkg_model_application_pb from 'pipecd/web/model/application_pb';
import * as pkg_model_application_live_state_pb from 'pipecd/web/model/application_live_state_pb';
import * as pkg_model_command_pb from 'pipecd/web/model/command_pb';
import * as pkg_model_deployment_pb from 'pipecd/web/model/deployment_pb';
import * as pkg_model_deployment_chain_pb from 'pipecd/web/model/deployment_chain_pb';
import * as pkg_model_logblock_pb from 'pipecd/web/model/logblock_pb';
import * as pkg_model_piped_pb from 'pipecd/web/model/piped_pb';
import * as pkg_model_rbac_pb from 'pipecd/web/model/rbac_pb';
import * as pkg_model_project_pb from 'pipecd/web/model/project_pb';
import * as pkg_model_apikey_pb from 'pipecd/web/model/apikey_pb';
import * as pkg_model_event_pb from 'pipecd/web/model/event_pb';
import * as google_protobuf_wrappers_pb from 'google-protobuf/google/protobuf/wrappers_pb';
import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';


export class RegisterPipedRequest extends jspb.Message {
  getName(): string;
  setName(value: string): RegisterPipedRequest;

  getDesc(): string;
  setDesc(value: string): RegisterPipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterPipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterPipedRequest): RegisterPipedRequest.AsObject;
  static serializeBinaryToWriter(message: RegisterPipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterPipedRequest;
  static deserializeBinaryFromReader(message: RegisterPipedRequest, reader: jspb.BinaryReader): RegisterPipedRequest;
}

export namespace RegisterPipedRequest {
  export type AsObject = {
    name: string,
    desc: string,
  }
}

export class RegisterPipedResponse extends jspb.Message {
  getId(): string;
  setId(value: string): RegisterPipedResponse;

  getKey(): string;
  setKey(value: string): RegisterPipedResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RegisterPipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RegisterPipedResponse): RegisterPipedResponse.AsObject;
  static serializeBinaryToWriter(message: RegisterPipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RegisterPipedResponse;
  static deserializeBinaryFromReader(message: RegisterPipedResponse, reader: jspb.BinaryReader): RegisterPipedResponse;
}

export namespace RegisterPipedResponse {
  export type AsObject = {
    id: string,
    key: string,
  }
}

export class UpdatePipedRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): UpdatePipedRequest;

  getName(): string;
  setName(value: string): UpdatePipedRequest;

  getDesc(): string;
  setDesc(value: string): UpdatePipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePipedRequest): UpdatePipedRequest.AsObject;
  static serializeBinaryToWriter(message: UpdatePipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePipedRequest;
  static deserializeBinaryFromReader(message: UpdatePipedRequest, reader: jspb.BinaryReader): UpdatePipedRequest;
}

export namespace UpdatePipedRequest {
  export type AsObject = {
    pipedId: string,
    name: string,
    desc: string,
  }
}

export class UpdatePipedResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePipedResponse): UpdatePipedResponse.AsObject;
  static serializeBinaryToWriter(message: UpdatePipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePipedResponse;
  static deserializeBinaryFromReader(message: UpdatePipedResponse, reader: jspb.BinaryReader): UpdatePipedResponse;
}

export namespace UpdatePipedResponse {
  export type AsObject = {
  }
}

export class RecreatePipedKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): RecreatePipedKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RecreatePipedKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RecreatePipedKeyRequest): RecreatePipedKeyRequest.AsObject;
  static serializeBinaryToWriter(message: RecreatePipedKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RecreatePipedKeyRequest;
  static deserializeBinaryFromReader(message: RecreatePipedKeyRequest, reader: jspb.BinaryReader): RecreatePipedKeyRequest;
}

export namespace RecreatePipedKeyRequest {
  export type AsObject = {
    id: string,
  }
}

export class RecreatePipedKeyResponse extends jspb.Message {
  getKey(): string;
  setKey(value: string): RecreatePipedKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RecreatePipedKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RecreatePipedKeyResponse): RecreatePipedKeyResponse.AsObject;
  static serializeBinaryToWriter(message: RecreatePipedKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RecreatePipedKeyResponse;
  static deserializeBinaryFromReader(message: RecreatePipedKeyResponse, reader: jspb.BinaryReader): RecreatePipedKeyResponse;
}

export namespace RecreatePipedKeyResponse {
  export type AsObject = {
    key: string,
  }
}

export class DeleteOldPipedKeysRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): DeleteOldPipedKeysRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOldPipedKeysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOldPipedKeysRequest): DeleteOldPipedKeysRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteOldPipedKeysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOldPipedKeysRequest;
  static deserializeBinaryFromReader(message: DeleteOldPipedKeysRequest, reader: jspb.BinaryReader): DeleteOldPipedKeysRequest;
}

export namespace DeleteOldPipedKeysRequest {
  export type AsObject = {
    pipedId: string,
  }
}

export class DeleteOldPipedKeysResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteOldPipedKeysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteOldPipedKeysResponse): DeleteOldPipedKeysResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteOldPipedKeysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteOldPipedKeysResponse;
  static deserializeBinaryFromReader(message: DeleteOldPipedKeysResponse, reader: jspb.BinaryReader): DeleteOldPipedKeysResponse;
}

export namespace DeleteOldPipedKeysResponse {
  export type AsObject = {
  }
}

export class EnablePipedRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): EnablePipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnablePipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnablePipedRequest): EnablePipedRequest.AsObject;
  static serializeBinaryToWriter(message: EnablePipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnablePipedRequest;
  static deserializeBinaryFromReader(message: EnablePipedRequest, reader: jspb.BinaryReader): EnablePipedRequest;
}

export namespace EnablePipedRequest {
  export type AsObject = {
    pipedId: string,
  }
}

export class EnablePipedResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnablePipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnablePipedResponse): EnablePipedResponse.AsObject;
  static serializeBinaryToWriter(message: EnablePipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnablePipedResponse;
  static deserializeBinaryFromReader(message: EnablePipedResponse, reader: jspb.BinaryReader): EnablePipedResponse;
}

export namespace EnablePipedResponse {
  export type AsObject = {
  }
}

export class DisablePipedRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): DisablePipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisablePipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisablePipedRequest): DisablePipedRequest.AsObject;
  static serializeBinaryToWriter(message: DisablePipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisablePipedRequest;
  static deserializeBinaryFromReader(message: DisablePipedRequest, reader: jspb.BinaryReader): DisablePipedRequest;
}

export namespace DisablePipedRequest {
  export type AsObject = {
    pipedId: string,
  }
}

export class DisablePipedResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisablePipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisablePipedResponse): DisablePipedResponse.AsObject;
  static serializeBinaryToWriter(message: DisablePipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisablePipedResponse;
  static deserializeBinaryFromReader(message: DisablePipedResponse, reader: jspb.BinaryReader): DisablePipedResponse;
}

export namespace DisablePipedResponse {
  export type AsObject = {
  }
}

export class ListPipedsRequest extends jspb.Message {
  getWithStatus(): boolean;
  setWithStatus(value: boolean): ListPipedsRequest;

  getOptions(): ListPipedsRequest.Options | undefined;
  setOptions(value?: ListPipedsRequest.Options): ListPipedsRequest;
  hasOptions(): boolean;
  clearOptions(): ListPipedsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPipedsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListPipedsRequest): ListPipedsRequest.AsObject;
  static serializeBinaryToWriter(message: ListPipedsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPipedsRequest;
  static deserializeBinaryFromReader(message: ListPipedsRequest, reader: jspb.BinaryReader): ListPipedsRequest;
}

export namespace ListPipedsRequest {
  export type AsObject = {
    withStatus: boolean,
    options?: ListPipedsRequest.Options.AsObject,
  }

  export class Options extends jspb.Message {
    getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): Options;
    hasEnabled(): boolean;
    clearEnabled(): Options;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
      enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

}

export class ListPipedsResponse extends jspb.Message {
  getPipedsList(): Array<pkg_model_piped_pb.Piped>;
  setPipedsList(value: Array<pkg_model_piped_pb.Piped>): ListPipedsResponse;
  clearPipedsList(): ListPipedsResponse;
  addPipeds(value?: pkg_model_piped_pb.Piped, index?: number): pkg_model_piped_pb.Piped;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListPipedsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListPipedsResponse): ListPipedsResponse.AsObject;
  static serializeBinaryToWriter(message: ListPipedsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListPipedsResponse;
  static deserializeBinaryFromReader(message: ListPipedsResponse, reader: jspb.BinaryReader): ListPipedsResponse;
}

export namespace ListPipedsResponse {
  export type AsObject = {
    pipedsList: Array<pkg_model_piped_pb.Piped.AsObject>,
  }
}

export class GetPipedRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): GetPipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetPipedRequest): GetPipedRequest.AsObject;
  static serializeBinaryToWriter(message: GetPipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPipedRequest;
  static deserializeBinaryFromReader(message: GetPipedRequest, reader: jspb.BinaryReader): GetPipedRequest;
}

export namespace GetPipedRequest {
  export type AsObject = {
    pipedId: string,
  }
}

export class GetPipedResponse extends jspb.Message {
  getPiped(): pkg_model_piped_pb.Piped | undefined;
  setPiped(value?: pkg_model_piped_pb.Piped): GetPipedResponse;
  hasPiped(): boolean;
  clearPiped(): GetPipedResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetPipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetPipedResponse): GetPipedResponse.AsObject;
  static serializeBinaryToWriter(message: GetPipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetPipedResponse;
  static deserializeBinaryFromReader(message: GetPipedResponse, reader: jspb.BinaryReader): GetPipedResponse;
}

export namespace GetPipedResponse {
  export type AsObject = {
    piped?: pkg_model_piped_pb.Piped.AsObject,
  }
}

export class UpdatePipedDesiredVersionRequest extends jspb.Message {
  getVersion(): string;
  setVersion(value: string): UpdatePipedDesiredVersionRequest;

  getPipedIdsList(): Array<string>;
  setPipedIdsList(value: Array<string>): UpdatePipedDesiredVersionRequest;
  clearPipedIdsList(): UpdatePipedDesiredVersionRequest;
  addPipedIds(value: string, index?: number): UpdatePipedDesiredVersionRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePipedDesiredVersionRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePipedDesiredVersionRequest): UpdatePipedDesiredVersionRequest.AsObject;
  static serializeBinaryToWriter(message: UpdatePipedDesiredVersionRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePipedDesiredVersionRequest;
  static deserializeBinaryFromReader(message: UpdatePipedDesiredVersionRequest, reader: jspb.BinaryReader): UpdatePipedDesiredVersionRequest;
}

export namespace UpdatePipedDesiredVersionRequest {
  export type AsObject = {
    version: string,
    pipedIdsList: Array<string>,
  }
}

export class UpdatePipedDesiredVersionResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdatePipedDesiredVersionResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdatePipedDesiredVersionResponse): UpdatePipedDesiredVersionResponse.AsObject;
  static serializeBinaryToWriter(message: UpdatePipedDesiredVersionResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdatePipedDesiredVersionResponse;
  static deserializeBinaryFromReader(message: UpdatePipedDesiredVersionResponse, reader: jspb.BinaryReader): UpdatePipedDesiredVersionResponse;
}

export namespace UpdatePipedDesiredVersionResponse {
  export type AsObject = {
  }
}

export class RestartPipedRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): RestartPipedRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RestartPipedRequest.AsObject;
  static toObject(includeInstance: boolean, msg: RestartPipedRequest): RestartPipedRequest.AsObject;
  static serializeBinaryToWriter(message: RestartPipedRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RestartPipedRequest;
  static deserializeBinaryFromReader(message: RestartPipedRequest, reader: jspb.BinaryReader): RestartPipedRequest;
}

export namespace RestartPipedRequest {
  export type AsObject = {
    pipedId: string,
  }
}

export class RestartPipedResponse extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): RestartPipedResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): RestartPipedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: RestartPipedResponse): RestartPipedResponse.AsObject;
  static serializeBinaryToWriter(message: RestartPipedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): RestartPipedResponse;
  static deserializeBinaryFromReader(message: RestartPipedResponse, reader: jspb.BinaryReader): RestartPipedResponse;
}

export namespace RestartPipedResponse {
  export type AsObject = {
    commandId: string,
  }
}

export class ListReleasedVersionsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListReleasedVersionsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListReleasedVersionsRequest): ListReleasedVersionsRequest.AsObject;
  static serializeBinaryToWriter(message: ListReleasedVersionsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListReleasedVersionsRequest;
  static deserializeBinaryFromReader(message: ListReleasedVersionsRequest, reader: jspb.BinaryReader): ListReleasedVersionsRequest;
}

export namespace ListReleasedVersionsRequest {
  export type AsObject = {
  }
}

export class ListReleasedVersionsResponse extends jspb.Message {
  getVersionsList(): Array<string>;
  setVersionsList(value: Array<string>): ListReleasedVersionsResponse;
  clearVersionsList(): ListReleasedVersionsResponse;
  addVersions(value: string, index?: number): ListReleasedVersionsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListReleasedVersionsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListReleasedVersionsResponse): ListReleasedVersionsResponse.AsObject;
  static serializeBinaryToWriter(message: ListReleasedVersionsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListReleasedVersionsResponse;
  static deserializeBinaryFromReader(message: ListReleasedVersionsResponse, reader: jspb.BinaryReader): ListReleasedVersionsResponse;
}

export namespace ListReleasedVersionsResponse {
  export type AsObject = {
    versionsList: Array<string>,
  }
}

export class ListDeprecatedNotesRequest extends jspb.Message {
  getProjectId(): string;
  setProjectId(value: string): ListDeprecatedNotesRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeprecatedNotesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeprecatedNotesRequest): ListDeprecatedNotesRequest.AsObject;
  static serializeBinaryToWriter(message: ListDeprecatedNotesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeprecatedNotesRequest;
  static deserializeBinaryFromReader(message: ListDeprecatedNotesRequest, reader: jspb.BinaryReader): ListDeprecatedNotesRequest;
}

export namespace ListDeprecatedNotesRequest {
  export type AsObject = {
    projectId: string,
  }
}

export class ListDeprecatedNotesResponse extends jspb.Message {
  getNotes(): string;
  setNotes(value: string): ListDeprecatedNotesResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeprecatedNotesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeprecatedNotesResponse): ListDeprecatedNotesResponse.AsObject;
  static serializeBinaryToWriter(message: ListDeprecatedNotesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeprecatedNotesResponse;
  static deserializeBinaryFromReader(message: ListDeprecatedNotesResponse, reader: jspb.BinaryReader): ListDeprecatedNotesResponse;
}

export namespace ListDeprecatedNotesResponse {
  export type AsObject = {
    notes: string,
  }
}

export class AddApplicationRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddApplicationRequest;

  getPipedId(): string;
  setPipedId(value: string): AddApplicationRequest;

  getGitPath(): pkg_model_common_pb.ApplicationGitPath | undefined;
  setGitPath(value?: pkg_model_common_pb.ApplicationGitPath): AddApplicationRequest;
  hasGitPath(): boolean;
  clearGitPath(): AddApplicationRequest;

  getKind(): pkg_model_common_pb.ApplicationKind;
  setKind(value: pkg_model_common_pb.ApplicationKind): AddApplicationRequest;

  getPlatformProvider(): string;
  setPlatformProvider(value: string): AddApplicationRequest;

  getDeployTargetsByPluginMap(): jspb.Map<string, pkg_model_deployment_pb.DeployTargets>;
  clearDeployTargetsByPluginMap(): AddApplicationRequest;

  getDescription(): string;
  setDescription(value: string): AddApplicationRequest;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): AddApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddApplicationRequest): AddApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: AddApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddApplicationRequest;
  static deserializeBinaryFromReader(message: AddApplicationRequest, reader: jspb.BinaryReader): AddApplicationRequest;
}

export namespace AddApplicationRequest {
  export type AsObject = {
    name: string,
    pipedId: string,
    gitPath?: pkg_model_common_pb.ApplicationGitPath.AsObject,
    kind: pkg_model_common_pb.ApplicationKind,
    platformProvider: string,
    deployTargetsByPluginMap: Array<[string, pkg_model_deployment_pb.DeployTargets.AsObject]>,
    description: string,
    labelsMap: Array<[string, string]>,
  }
}

export class AddApplicationResponse extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): AddApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddApplicationResponse): AddApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: AddApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddApplicationResponse;
  static deserializeBinaryFromReader(message: AddApplicationResponse, reader: jspb.BinaryReader): AddApplicationResponse;
}

export namespace AddApplicationResponse {
  export type AsObject = {
    applicationId: string,
  }
}

export class UpdateApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): UpdateApplicationRequest;

  getName(): string;
  setName(value: string): UpdateApplicationRequest;

  getPipedId(): string;
  setPipedId(value: string): UpdateApplicationRequest;

  getKind(): pkg_model_common_pb.ApplicationKind;
  setKind(value: pkg_model_common_pb.ApplicationKind): UpdateApplicationRequest;

  getPlatformProvider(): string;
  setPlatformProvider(value: string): UpdateApplicationRequest;

  getDeployTargetsByPluginMap(): jspb.Map<string, pkg_model_deployment_pb.DeployTargets>;
  clearDeployTargetsByPluginMap(): UpdateApplicationRequest;

  getConfigFilename(): string;
  setConfigFilename(value: string): UpdateApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateApplicationRequest): UpdateApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateApplicationRequest;
  static deserializeBinaryFromReader(message: UpdateApplicationRequest, reader: jspb.BinaryReader): UpdateApplicationRequest;
}

export namespace UpdateApplicationRequest {
  export type AsObject = {
    applicationId: string,
    name: string,
    pipedId: string,
    kind: pkg_model_common_pb.ApplicationKind,
    platformProvider: string,
    deployTargetsByPluginMap: Array<[string, pkg_model_deployment_pb.DeployTargets.AsObject]>,
    configFilename: string,
  }
}

export class UpdateApplicationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateApplicationResponse): UpdateApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateApplicationResponse;
  static deserializeBinaryFromReader(message: UpdateApplicationResponse, reader: jspb.BinaryReader): UpdateApplicationResponse;
}

export namespace UpdateApplicationResponse {
  export type AsObject = {
  }
}

export class EnableApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): EnableApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableApplicationRequest): EnableApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: EnableApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableApplicationRequest;
  static deserializeBinaryFromReader(message: EnableApplicationRequest, reader: jspb.BinaryReader): EnableApplicationRequest;
}

export namespace EnableApplicationRequest {
  export type AsObject = {
    applicationId: string,
  }
}

export class EnableApplicationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableApplicationResponse): EnableApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: EnableApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableApplicationResponse;
  static deserializeBinaryFromReader(message: EnableApplicationResponse, reader: jspb.BinaryReader): EnableApplicationResponse;
}

export namespace EnableApplicationResponse {
  export type AsObject = {
  }
}

export class DisableApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): DisableApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableApplicationRequest): DisableApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: DisableApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableApplicationRequest;
  static deserializeBinaryFromReader(message: DisableApplicationRequest, reader: jspb.BinaryReader): DisableApplicationRequest;
}

export namespace DisableApplicationRequest {
  export type AsObject = {
    applicationId: string,
  }
}

export class DisableApplicationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableApplicationResponse): DisableApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: DisableApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableApplicationResponse;
  static deserializeBinaryFromReader(message: DisableApplicationResponse, reader: jspb.BinaryReader): DisableApplicationResponse;
}

export namespace DisableApplicationResponse {
  export type AsObject = {
  }
}

export class DeleteApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): DeleteApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteApplicationRequest): DeleteApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteApplicationRequest;
  static deserializeBinaryFromReader(message: DeleteApplicationRequest, reader: jspb.BinaryReader): DeleteApplicationRequest;
}

export namespace DeleteApplicationRequest {
  export type AsObject = {
    applicationId: string,
  }
}

export class DeleteApplicationResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteApplicationResponse): DeleteApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteApplicationResponse;
  static deserializeBinaryFromReader(message: DeleteApplicationResponse, reader: jspb.BinaryReader): DeleteApplicationResponse;
}

export namespace DeleteApplicationResponse {
  export type AsObject = {
  }
}

export class ListApplicationsRequest extends jspb.Message {
  getOptions(): ListApplicationsRequest.Options | undefined;
  setOptions(value?: ListApplicationsRequest.Options): ListApplicationsRequest;
  hasOptions(): boolean;
  clearOptions(): ListApplicationsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListApplicationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListApplicationsRequest): ListApplicationsRequest.AsObject;
  static serializeBinaryToWriter(message: ListApplicationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListApplicationsRequest;
  static deserializeBinaryFromReader(message: ListApplicationsRequest, reader: jspb.BinaryReader): ListApplicationsRequest;
}

export namespace ListApplicationsRequest {
  export type AsObject = {
    options?: ListApplicationsRequest.Options.AsObject,
  }

  export class Options extends jspb.Message {
    getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): Options;
    hasEnabled(): boolean;
    clearEnabled(): Options;

    getKindsList(): Array<pkg_model_common_pb.ApplicationKind>;
    setKindsList(value: Array<pkg_model_common_pb.ApplicationKind>): Options;
    clearKindsList(): Options;
    addKinds(value: pkg_model_common_pb.ApplicationKind, index?: number): Options;

    getSyncStatusesList(): Array<pkg_model_application_pb.ApplicationSyncStatus>;
    setSyncStatusesList(value: Array<pkg_model_application_pb.ApplicationSyncStatus>): Options;
    clearSyncStatusesList(): Options;
    addSyncStatuses(value: pkg_model_application_pb.ApplicationSyncStatus, index?: number): Options;

    getName(): string;
    setName(value: string): Options;

    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): Options;

    getPipedId(): string;
    setPipedId(value: string): Options;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
      enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
      kindsList: Array<pkg_model_common_pb.ApplicationKind>,
      syncStatusesList: Array<pkg_model_application_pb.ApplicationSyncStatus>,
      name: string,
      labelsMap: Array<[string, string]>,
      pipedId: string,
    }
  }

}

export class ListApplicationsResponse extends jspb.Message {
  getApplicationsList(): Array<pkg_model_application_pb.Application>;
  setApplicationsList(value: Array<pkg_model_application_pb.Application>): ListApplicationsResponse;
  clearApplicationsList(): ListApplicationsResponse;
  addApplications(value?: pkg_model_application_pb.Application, index?: number): pkg_model_application_pb.Application;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListApplicationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListApplicationsResponse): ListApplicationsResponse.AsObject;
  static serializeBinaryToWriter(message: ListApplicationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListApplicationsResponse;
  static deserializeBinaryFromReader(message: ListApplicationsResponse, reader: jspb.BinaryReader): ListApplicationsResponse;
}

export namespace ListApplicationsResponse {
  export type AsObject = {
    applicationsList: Array<pkg_model_application_pb.Application.AsObject>,
  }
}

export class SyncApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): SyncApplicationRequest;

  getSyncStrategy(): pkg_model_common_pb.SyncStrategy;
  setSyncStrategy(value: pkg_model_common_pb.SyncStrategy): SyncApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SyncApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SyncApplicationRequest): SyncApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: SyncApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SyncApplicationRequest;
  static deserializeBinaryFromReader(message: SyncApplicationRequest, reader: jspb.BinaryReader): SyncApplicationRequest;
}

export namespace SyncApplicationRequest {
  export type AsObject = {
    applicationId: string,
    syncStrategy: pkg_model_common_pb.SyncStrategy,
  }
}

export class SyncApplicationResponse extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): SyncApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SyncApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SyncApplicationResponse): SyncApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: SyncApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SyncApplicationResponse;
  static deserializeBinaryFromReader(message: SyncApplicationResponse, reader: jspb.BinaryReader): SyncApplicationResponse;
}

export namespace SyncApplicationResponse {
  export type AsObject = {
    commandId: string,
  }
}

export class GetApplicationRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): GetApplicationRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationRequest): GetApplicationRequest.AsObject;
  static serializeBinaryToWriter(message: GetApplicationRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationRequest;
  static deserializeBinaryFromReader(message: GetApplicationRequest, reader: jspb.BinaryReader): GetApplicationRequest;
}

export namespace GetApplicationRequest {
  export type AsObject = {
    applicationId: string,
  }
}

export class GetApplicationResponse extends jspb.Message {
  getApplication(): pkg_model_application_pb.Application | undefined;
  setApplication(value?: pkg_model_application_pb.Application): GetApplicationResponse;
  hasApplication(): boolean;
  clearApplication(): GetApplicationResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationResponse): GetApplicationResponse.AsObject;
  static serializeBinaryToWriter(message: GetApplicationResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationResponse;
  static deserializeBinaryFromReader(message: GetApplicationResponse, reader: jspb.BinaryReader): GetApplicationResponse;
}

export namespace GetApplicationResponse {
  export type AsObject = {
    application?: pkg_model_application_pb.Application.AsObject,
  }
}

export class GenerateApplicationSealedSecretRequest extends jspb.Message {
  getPipedId(): string;
  setPipedId(value: string): GenerateApplicationSealedSecretRequest;

  getData(): string;
  setData(value: string): GenerateApplicationSealedSecretRequest;

  getBase64Encoding(): boolean;
  setBase64Encoding(value: boolean): GenerateApplicationSealedSecretRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateApplicationSealedSecretRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateApplicationSealedSecretRequest): GenerateApplicationSealedSecretRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateApplicationSealedSecretRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateApplicationSealedSecretRequest;
  static deserializeBinaryFromReader(message: GenerateApplicationSealedSecretRequest, reader: jspb.BinaryReader): GenerateApplicationSealedSecretRequest;
}

export namespace GenerateApplicationSealedSecretRequest {
  export type AsObject = {
    pipedId: string,
    data: string,
    base64Encoding: boolean,
  }
}

export class GenerateApplicationSealedSecretResponse extends jspb.Message {
  getData(): string;
  setData(value: string): GenerateApplicationSealedSecretResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateApplicationSealedSecretResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateApplicationSealedSecretResponse): GenerateApplicationSealedSecretResponse.AsObject;
  static serializeBinaryToWriter(message: GenerateApplicationSealedSecretResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateApplicationSealedSecretResponse;
  static deserializeBinaryFromReader(message: GenerateApplicationSealedSecretResponse, reader: jspb.BinaryReader): GenerateApplicationSealedSecretResponse;
}

export namespace GenerateApplicationSealedSecretResponse {
  export type AsObject = {
    data: string,
  }
}

export class ListUnregisteredApplicationsRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUnregisteredApplicationsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListUnregisteredApplicationsRequest): ListUnregisteredApplicationsRequest.AsObject;
  static serializeBinaryToWriter(message: ListUnregisteredApplicationsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUnregisteredApplicationsRequest;
  static deserializeBinaryFromReader(message: ListUnregisteredApplicationsRequest, reader: jspb.BinaryReader): ListUnregisteredApplicationsRequest;
}

export namespace ListUnregisteredApplicationsRequest {
  export type AsObject = {
  }
}

export class ListUnregisteredApplicationsResponse extends jspb.Message {
  getApplicationsList(): Array<pkg_model_common_pb.ApplicationInfo>;
  setApplicationsList(value: Array<pkg_model_common_pb.ApplicationInfo>): ListUnregisteredApplicationsResponse;
  clearApplicationsList(): ListUnregisteredApplicationsResponse;
  addApplications(value?: pkg_model_common_pb.ApplicationInfo, index?: number): pkg_model_common_pb.ApplicationInfo;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListUnregisteredApplicationsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListUnregisteredApplicationsResponse): ListUnregisteredApplicationsResponse.AsObject;
  static serializeBinaryToWriter(message: ListUnregisteredApplicationsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListUnregisteredApplicationsResponse;
  static deserializeBinaryFromReader(message: ListUnregisteredApplicationsResponse, reader: jspb.BinaryReader): ListUnregisteredApplicationsResponse;
}

export namespace ListUnregisteredApplicationsResponse {
  export type AsObject = {
    applicationsList: Array<pkg_model_common_pb.ApplicationInfo.AsObject>,
  }
}

export class ListDeploymentsRequest extends jspb.Message {
  getOptions(): ListDeploymentsRequest.Options | undefined;
  setOptions(value?: ListDeploymentsRequest.Options): ListDeploymentsRequest;
  hasOptions(): boolean;
  clearOptions(): ListDeploymentsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListDeploymentsRequest;

  getCursor(): string;
  setCursor(value: string): ListDeploymentsRequest;

  getPageMinUpdatedAt(): number;
  setPageMinUpdatedAt(value: number): ListDeploymentsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentsRequest): ListDeploymentsRequest.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentsRequest;
  static deserializeBinaryFromReader(message: ListDeploymentsRequest, reader: jspb.BinaryReader): ListDeploymentsRequest;
}

export namespace ListDeploymentsRequest {
  export type AsObject = {
    options?: ListDeploymentsRequest.Options.AsObject,
    pageSize: number,
    cursor: string,
    pageMinUpdatedAt: number,
  }

  export class Options extends jspb.Message {
    getStatusesList(): Array<pkg_model_deployment_pb.DeploymentStatus>;
    setStatusesList(value: Array<pkg_model_deployment_pb.DeploymentStatus>): Options;
    clearStatusesList(): Options;
    addStatuses(value: pkg_model_deployment_pb.DeploymentStatus, index?: number): Options;

    getKindsList(): Array<pkg_model_common_pb.ApplicationKind>;
    setKindsList(value: Array<pkg_model_common_pb.ApplicationKind>): Options;
    clearKindsList(): Options;
    addKinds(value: pkg_model_common_pb.ApplicationKind, index?: number): Options;

    getApplicationIdsList(): Array<string>;
    setApplicationIdsList(value: Array<string>): Options;
    clearApplicationIdsList(): Options;
    addApplicationIds(value: string, index?: number): Options;

    getApplicationName(): string;
    setApplicationName(value: string): Options;

    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): Options;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
      statusesList: Array<pkg_model_deployment_pb.DeploymentStatus>,
      kindsList: Array<pkg_model_common_pb.ApplicationKind>,
      applicationIdsList: Array<string>,
      applicationName: string,
      labelsMap: Array<[string, string]>,
    }
  }

}

export class ListDeploymentsResponse extends jspb.Message {
  getDeploymentsList(): Array<pkg_model_deployment_pb.Deployment>;
  setDeploymentsList(value: Array<pkg_model_deployment_pb.Deployment>): ListDeploymentsResponse;
  clearDeploymentsList(): ListDeploymentsResponse;
  addDeployments(value?: pkg_model_deployment_pb.Deployment, index?: number): pkg_model_deployment_pb.Deployment;

  getCursor(): string;
  setCursor(value: string): ListDeploymentsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentsResponse): ListDeploymentsResponse.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentsResponse;
  static deserializeBinaryFromReader(message: ListDeploymentsResponse, reader: jspb.BinaryReader): ListDeploymentsResponse;
}

export namespace ListDeploymentsResponse {
  export type AsObject = {
    deploymentsList: Array<pkg_model_deployment_pb.Deployment.AsObject>,
    cursor: string,
  }
}

export class GetDeploymentRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): GetDeploymentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDeploymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDeploymentRequest): GetDeploymentRequest.AsObject;
  static serializeBinaryToWriter(message: GetDeploymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDeploymentRequest;
  static deserializeBinaryFromReader(message: GetDeploymentRequest, reader: jspb.BinaryReader): GetDeploymentRequest;
}

export namespace GetDeploymentRequest {
  export type AsObject = {
    deploymentId: string,
  }
}

export class GetDeploymentResponse extends jspb.Message {
  getDeployment(): pkg_model_deployment_pb.Deployment | undefined;
  setDeployment(value?: pkg_model_deployment_pb.Deployment): GetDeploymentResponse;
  hasDeployment(): boolean;
  clearDeployment(): GetDeploymentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDeploymentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDeploymentResponse): GetDeploymentResponse.AsObject;
  static serializeBinaryToWriter(message: GetDeploymentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDeploymentResponse;
  static deserializeBinaryFromReader(message: GetDeploymentResponse, reader: jspb.BinaryReader): GetDeploymentResponse;
}

export namespace GetDeploymentResponse {
  export type AsObject = {
    deployment?: pkg_model_deployment_pb.Deployment.AsObject,
  }
}

export class GetStageLogRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): GetStageLogRequest;

  getStageId(): string;
  setStageId(value: string): GetStageLogRequest;

  getRetriedCount(): number;
  setRetriedCount(value: number): GetStageLogRequest;

  getOffsetIndex(): number;
  setOffsetIndex(value: number): GetStageLogRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetStageLogRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetStageLogRequest): GetStageLogRequest.AsObject;
  static serializeBinaryToWriter(message: GetStageLogRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetStageLogRequest;
  static deserializeBinaryFromReader(message: GetStageLogRequest, reader: jspb.BinaryReader): GetStageLogRequest;
}

export namespace GetStageLogRequest {
  export type AsObject = {
    deploymentId: string,
    stageId: string,
    retriedCount: number,
    offsetIndex: number,
  }
}

export class GetStageLogResponse extends jspb.Message {
  getBlocksList(): Array<pkg_model_logblock_pb.LogBlock>;
  setBlocksList(value: Array<pkg_model_logblock_pb.LogBlock>): GetStageLogResponse;
  clearBlocksList(): GetStageLogResponse;
  addBlocks(value?: pkg_model_logblock_pb.LogBlock, index?: number): pkg_model_logblock_pb.LogBlock;

  getCompleted(): boolean;
  setCompleted(value: boolean): GetStageLogResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetStageLogResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetStageLogResponse): GetStageLogResponse.AsObject;
  static serializeBinaryToWriter(message: GetStageLogResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetStageLogResponse;
  static deserializeBinaryFromReader(message: GetStageLogResponse, reader: jspb.BinaryReader): GetStageLogResponse;
}

export namespace GetStageLogResponse {
  export type AsObject = {
    blocksList: Array<pkg_model_logblock_pb.LogBlock.AsObject>,
    completed: boolean,
  }
}

export class CancelDeploymentRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): CancelDeploymentRequest;

  getForceRollback(): boolean;
  setForceRollback(value: boolean): CancelDeploymentRequest;

  getForceNoRollback(): boolean;
  setForceNoRollback(value: boolean): CancelDeploymentRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelDeploymentRequest.AsObject;
  static toObject(includeInstance: boolean, msg: CancelDeploymentRequest): CancelDeploymentRequest.AsObject;
  static serializeBinaryToWriter(message: CancelDeploymentRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelDeploymentRequest;
  static deserializeBinaryFromReader(message: CancelDeploymentRequest, reader: jspb.BinaryReader): CancelDeploymentRequest;
}

export namespace CancelDeploymentRequest {
  export type AsObject = {
    deploymentId: string,
    forceRollback: boolean,
    forceNoRollback: boolean,
  }
}

export class CancelDeploymentResponse extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): CancelDeploymentResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): CancelDeploymentResponse.AsObject;
  static toObject(includeInstance: boolean, msg: CancelDeploymentResponse): CancelDeploymentResponse.AsObject;
  static serializeBinaryToWriter(message: CancelDeploymentResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): CancelDeploymentResponse;
  static deserializeBinaryFromReader(message: CancelDeploymentResponse, reader: jspb.BinaryReader): CancelDeploymentResponse;
}

export namespace CancelDeploymentResponse {
  export type AsObject = {
    commandId: string,
  }
}

export class SkipStageRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): SkipStageRequest;

  getStageId(): string;
  setStageId(value: string): SkipStageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SkipStageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: SkipStageRequest): SkipStageRequest.AsObject;
  static serializeBinaryToWriter(message: SkipStageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SkipStageRequest;
  static deserializeBinaryFromReader(message: SkipStageRequest, reader: jspb.BinaryReader): SkipStageRequest;
}

export namespace SkipStageRequest {
  export type AsObject = {
    deploymentId: string,
    stageId: string,
  }
}

export class SkipStageResponse extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): SkipStageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): SkipStageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: SkipStageResponse): SkipStageResponse.AsObject;
  static serializeBinaryToWriter(message: SkipStageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): SkipStageResponse;
  static deserializeBinaryFromReader(message: SkipStageResponse, reader: jspb.BinaryReader): SkipStageResponse;
}

export namespace SkipStageResponse {
  export type AsObject = {
    commandId: string,
  }
}

export class ApproveStageRequest extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): ApproveStageRequest;

  getStageId(): string;
  setStageId(value: string): ApproveStageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApproveStageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ApproveStageRequest): ApproveStageRequest.AsObject;
  static serializeBinaryToWriter(message: ApproveStageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApproveStageRequest;
  static deserializeBinaryFromReader(message: ApproveStageRequest, reader: jspb.BinaryReader): ApproveStageRequest;
}

export namespace ApproveStageRequest {
  export type AsObject = {
    deploymentId: string,
    stageId: string,
  }
}

export class ApproveStageResponse extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): ApproveStageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ApproveStageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ApproveStageResponse): ApproveStageResponse.AsObject;
  static serializeBinaryToWriter(message: ApproveStageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ApproveStageResponse;
  static deserializeBinaryFromReader(message: ApproveStageResponse, reader: jspb.BinaryReader): ApproveStageResponse;
}

export namespace ApproveStageResponse {
  export type AsObject = {
    commandId: string,
  }
}

export class GetApplicationLiveStateRequest extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): GetApplicationLiveStateRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationLiveStateRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationLiveStateRequest): GetApplicationLiveStateRequest.AsObject;
  static serializeBinaryToWriter(message: GetApplicationLiveStateRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationLiveStateRequest;
  static deserializeBinaryFromReader(message: GetApplicationLiveStateRequest, reader: jspb.BinaryReader): GetApplicationLiveStateRequest;
}

export namespace GetApplicationLiveStateRequest {
  export type AsObject = {
    applicationId: string,
  }
}

export class GetApplicationLiveStateResponse extends jspb.Message {
  getSnapshot(): pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot | undefined;
  setSnapshot(value?: pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot): GetApplicationLiveStateResponse;
  hasSnapshot(): boolean;
  clearSnapshot(): GetApplicationLiveStateResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetApplicationLiveStateResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetApplicationLiveStateResponse): GetApplicationLiveStateResponse.AsObject;
  static serializeBinaryToWriter(message: GetApplicationLiveStateResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetApplicationLiveStateResponse;
  static deserializeBinaryFromReader(message: GetApplicationLiveStateResponse, reader: jspb.BinaryReader): GetApplicationLiveStateResponse;
}

export namespace GetApplicationLiveStateResponse {
  export type AsObject = {
    snapshot?: pkg_model_application_live_state_pb.ApplicationLiveStateSnapshot.AsObject,
  }
}

export class GetProjectRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectRequest): GetProjectRequest.AsObject;
  static serializeBinaryToWriter(message: GetProjectRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectRequest;
  static deserializeBinaryFromReader(message: GetProjectRequest, reader: jspb.BinaryReader): GetProjectRequest;
}

export namespace GetProjectRequest {
  export type AsObject = {
  }
}

export class GetProjectResponse extends jspb.Message {
  getProject(): pkg_model_project_pb.Project | undefined;
  setProject(value?: pkg_model_project_pb.Project): GetProjectResponse;
  hasProject(): boolean;
  clearProject(): GetProjectResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetProjectResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetProjectResponse): GetProjectResponse.AsObject;
  static serializeBinaryToWriter(message: GetProjectResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetProjectResponse;
  static deserializeBinaryFromReader(message: GetProjectResponse, reader: jspb.BinaryReader): GetProjectResponse;
}

export namespace GetProjectResponse {
  export type AsObject = {
    project?: pkg_model_project_pb.Project.AsObject,
  }
}

export class UpdateProjectStaticAdminRequest extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): UpdateProjectStaticAdminRequest;

  getPassword(): string;
  setPassword(value: string): UpdateProjectStaticAdminRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectStaticAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectStaticAdminRequest): UpdateProjectStaticAdminRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectStaticAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectStaticAdminRequest;
  static deserializeBinaryFromReader(message: UpdateProjectStaticAdminRequest, reader: jspb.BinaryReader): UpdateProjectStaticAdminRequest;
}

export namespace UpdateProjectStaticAdminRequest {
  export type AsObject = {
    username: string,
    password: string,
  }
}

export class UpdateProjectStaticAdminResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectStaticAdminResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectStaticAdminResponse): UpdateProjectStaticAdminResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectStaticAdminResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectStaticAdminResponse;
  static deserializeBinaryFromReader(message: UpdateProjectStaticAdminResponse, reader: jspb.BinaryReader): UpdateProjectStaticAdminResponse;
}

export namespace UpdateProjectStaticAdminResponse {
  export type AsObject = {
  }
}

export class UpdateProjectSSOConfigRequest extends jspb.Message {
  getSso(): pkg_model_project_pb.ProjectSSOConfig | undefined;
  setSso(value?: pkg_model_project_pb.ProjectSSOConfig): UpdateProjectSSOConfigRequest;
  hasSso(): boolean;
  clearSso(): UpdateProjectSSOConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectSSOConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectSSOConfigRequest): UpdateProjectSSOConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectSSOConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectSSOConfigRequest;
  static deserializeBinaryFromReader(message: UpdateProjectSSOConfigRequest, reader: jspb.BinaryReader): UpdateProjectSSOConfigRequest;
}

export namespace UpdateProjectSSOConfigRequest {
  export type AsObject = {
    sso?: pkg_model_project_pb.ProjectSSOConfig.AsObject,
  }
}

export class UpdateProjectSSOConfigResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectSSOConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectSSOConfigResponse): UpdateProjectSSOConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectSSOConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectSSOConfigResponse;
  static deserializeBinaryFromReader(message: UpdateProjectSSOConfigResponse, reader: jspb.BinaryReader): UpdateProjectSSOConfigResponse;
}

export namespace UpdateProjectSSOConfigResponse {
  export type AsObject = {
  }
}

export class UpdateProjectRBACConfigRequest extends jspb.Message {
  getRbac(): pkg_model_project_pb.ProjectRBACConfig | undefined;
  setRbac(value?: pkg_model_project_pb.ProjectRBACConfig): UpdateProjectRBACConfigRequest;
  hasRbac(): boolean;
  clearRbac(): UpdateProjectRBACConfigRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRBACConfigRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRBACConfigRequest): UpdateProjectRBACConfigRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRBACConfigRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRBACConfigRequest;
  static deserializeBinaryFromReader(message: UpdateProjectRBACConfigRequest, reader: jspb.BinaryReader): UpdateProjectRBACConfigRequest;
}

export namespace UpdateProjectRBACConfigRequest {
  export type AsObject = {
    rbac?: pkg_model_project_pb.ProjectRBACConfig.AsObject,
  }
}

export class UpdateProjectRBACConfigResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRBACConfigResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRBACConfigResponse): UpdateProjectRBACConfigResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRBACConfigResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRBACConfigResponse;
  static deserializeBinaryFromReader(message: UpdateProjectRBACConfigResponse, reader: jspb.BinaryReader): UpdateProjectRBACConfigResponse;
}

export namespace UpdateProjectRBACConfigResponse {
  export type AsObject = {
  }
}

export class EnableStaticAdminRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableStaticAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: EnableStaticAdminRequest): EnableStaticAdminRequest.AsObject;
  static serializeBinaryToWriter(message: EnableStaticAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableStaticAdminRequest;
  static deserializeBinaryFromReader(message: EnableStaticAdminRequest, reader: jspb.BinaryReader): EnableStaticAdminRequest;
}

export namespace EnableStaticAdminRequest {
  export type AsObject = {
  }
}

export class EnableStaticAdminResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): EnableStaticAdminResponse.AsObject;
  static toObject(includeInstance: boolean, msg: EnableStaticAdminResponse): EnableStaticAdminResponse.AsObject;
  static serializeBinaryToWriter(message: EnableStaticAdminResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): EnableStaticAdminResponse;
  static deserializeBinaryFromReader(message: EnableStaticAdminResponse, reader: jspb.BinaryReader): EnableStaticAdminResponse;
}

export namespace EnableStaticAdminResponse {
  export type AsObject = {
  }
}

export class DisableStaticAdminRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableStaticAdminRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableStaticAdminRequest): DisableStaticAdminRequest.AsObject;
  static serializeBinaryToWriter(message: DisableStaticAdminRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableStaticAdminRequest;
  static deserializeBinaryFromReader(message: DisableStaticAdminRequest, reader: jspb.BinaryReader): DisableStaticAdminRequest;
}

export namespace DisableStaticAdminRequest {
  export type AsObject = {
  }
}

export class DisableStaticAdminResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableStaticAdminResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableStaticAdminResponse): DisableStaticAdminResponse.AsObject;
  static serializeBinaryToWriter(message: DisableStaticAdminResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableStaticAdminResponse;
  static deserializeBinaryFromReader(message: DisableStaticAdminResponse, reader: jspb.BinaryReader): DisableStaticAdminResponse;
}

export namespace DisableStaticAdminResponse {
  export type AsObject = {
  }
}

export class GetMeRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeRequest): GetMeRequest.AsObject;
  static serializeBinaryToWriter(message: GetMeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeRequest;
  static deserializeBinaryFromReader(message: GetMeRequest, reader: jspb.BinaryReader): GetMeRequest;
}

export namespace GetMeRequest {
  export type AsObject = {
  }
}

export class GetMeResponse extends jspb.Message {
  getSubject(): string;
  setSubject(value: string): GetMeResponse;

  getAvatarUrl(): string;
  setAvatarUrl(value: string): GetMeResponse;

  getProjectId(): string;
  setProjectId(value: string): GetMeResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetMeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetMeResponse): GetMeResponse.AsObject;
  static serializeBinaryToWriter(message: GetMeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetMeResponse;
  static deserializeBinaryFromReader(message: GetMeResponse, reader: jspb.BinaryReader): GetMeResponse;
}

export namespace GetMeResponse {
  export type AsObject = {
    subject: string,
    avatarUrl: string,
    projectId: string,
  }
}

export class AddProjectRBACRoleRequest extends jspb.Message {
  getName(): string;
  setName(value: string): AddProjectRBACRoleRequest;

  getPoliciesList(): Array<pkg_model_project_pb.ProjectRBACPolicy>;
  setPoliciesList(value: Array<pkg_model_project_pb.ProjectRBACPolicy>): AddProjectRBACRoleRequest;
  clearPoliciesList(): AddProjectRBACRoleRequest;
  addPolicies(value?: pkg_model_project_pb.ProjectRBACPolicy, index?: number): pkg_model_project_pb.ProjectRBACPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectRBACRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectRBACRoleRequest): AddProjectRBACRoleRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectRBACRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectRBACRoleRequest;
  static deserializeBinaryFromReader(message: AddProjectRBACRoleRequest, reader: jspb.BinaryReader): AddProjectRBACRoleRequest;
}

export namespace AddProjectRBACRoleRequest {
  export type AsObject = {
    name: string,
    policiesList: Array<pkg_model_project_pb.ProjectRBACPolicy.AsObject>,
  }
}

export class AddProjectRBACRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectRBACRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectRBACRoleResponse): AddProjectRBACRoleResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectRBACRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectRBACRoleResponse;
  static deserializeBinaryFromReader(message: AddProjectRBACRoleResponse, reader: jspb.BinaryReader): AddProjectRBACRoleResponse;
}

export namespace AddProjectRBACRoleResponse {
  export type AsObject = {
  }
}

export class UpdateProjectRBACRoleRequest extends jspb.Message {
  getName(): string;
  setName(value: string): UpdateProjectRBACRoleRequest;

  getPoliciesList(): Array<pkg_model_project_pb.ProjectRBACPolicy>;
  setPoliciesList(value: Array<pkg_model_project_pb.ProjectRBACPolicy>): UpdateProjectRBACRoleRequest;
  clearPoliciesList(): UpdateProjectRBACRoleRequest;
  addPolicies(value?: pkg_model_project_pb.ProjectRBACPolicy, index?: number): pkg_model_project_pb.ProjectRBACPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRBACRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRBACRoleRequest): UpdateProjectRBACRoleRequest.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRBACRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRBACRoleRequest;
  static deserializeBinaryFromReader(message: UpdateProjectRBACRoleRequest, reader: jspb.BinaryReader): UpdateProjectRBACRoleRequest;
}

export namespace UpdateProjectRBACRoleRequest {
  export type AsObject = {
    name: string,
    policiesList: Array<pkg_model_project_pb.ProjectRBACPolicy.AsObject>,
  }
}

export class UpdateProjectRBACRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UpdateProjectRBACRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UpdateProjectRBACRoleResponse): UpdateProjectRBACRoleResponse.AsObject;
  static serializeBinaryToWriter(message: UpdateProjectRBACRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UpdateProjectRBACRoleResponse;
  static deserializeBinaryFromReader(message: UpdateProjectRBACRoleResponse, reader: jspb.BinaryReader): UpdateProjectRBACRoleResponse;
}

export namespace UpdateProjectRBACRoleResponse {
  export type AsObject = {
  }
}

export class DeleteProjectRBACRoleRequest extends jspb.Message {
  getName(): string;
  setName(value: string): DeleteProjectRBACRoleRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProjectRBACRoleRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProjectRBACRoleRequest): DeleteProjectRBACRoleRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteProjectRBACRoleRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProjectRBACRoleRequest;
  static deserializeBinaryFromReader(message: DeleteProjectRBACRoleRequest, reader: jspb.BinaryReader): DeleteProjectRBACRoleRequest;
}

export namespace DeleteProjectRBACRoleRequest {
  export type AsObject = {
    name: string,
  }
}

export class DeleteProjectRBACRoleResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProjectRBACRoleResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProjectRBACRoleResponse): DeleteProjectRBACRoleResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteProjectRBACRoleResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProjectRBACRoleResponse;
  static deserializeBinaryFromReader(message: DeleteProjectRBACRoleResponse, reader: jspb.BinaryReader): DeleteProjectRBACRoleResponse;
}

export namespace DeleteProjectRBACRoleResponse {
  export type AsObject = {
  }
}

export class AddProjectUserGroupRequest extends jspb.Message {
  getSsoGroup(): string;
  setSsoGroup(value: string): AddProjectUserGroupRequest;

  getRole(): string;
  setRole(value: string): AddProjectUserGroupRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectUserGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectUserGroupRequest): AddProjectUserGroupRequest.AsObject;
  static serializeBinaryToWriter(message: AddProjectUserGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectUserGroupRequest;
  static deserializeBinaryFromReader(message: AddProjectUserGroupRequest, reader: jspb.BinaryReader): AddProjectUserGroupRequest;
}

export namespace AddProjectUserGroupRequest {
  export type AsObject = {
    ssoGroup: string,
    role: string,
  }
}

export class AddProjectUserGroupResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddProjectUserGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddProjectUserGroupResponse): AddProjectUserGroupResponse.AsObject;
  static serializeBinaryToWriter(message: AddProjectUserGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddProjectUserGroupResponse;
  static deserializeBinaryFromReader(message: AddProjectUserGroupResponse, reader: jspb.BinaryReader): AddProjectUserGroupResponse;
}

export namespace AddProjectUserGroupResponse {
  export type AsObject = {
  }
}

export class DeleteProjectUserGroupRequest extends jspb.Message {
  getSsoGroup(): string;
  setSsoGroup(value: string): DeleteProjectUserGroupRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProjectUserGroupRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProjectUserGroupRequest): DeleteProjectUserGroupRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteProjectUserGroupRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProjectUserGroupRequest;
  static deserializeBinaryFromReader(message: DeleteProjectUserGroupRequest, reader: jspb.BinaryReader): DeleteProjectUserGroupRequest;
}

export namespace DeleteProjectUserGroupRequest {
  export type AsObject = {
    ssoGroup: string,
  }
}

export class DeleteProjectUserGroupResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteProjectUserGroupResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteProjectUserGroupResponse): DeleteProjectUserGroupResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteProjectUserGroupResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteProjectUserGroupResponse;
  static deserializeBinaryFromReader(message: DeleteProjectUserGroupResponse, reader: jspb.BinaryReader): DeleteProjectUserGroupResponse;
}

export namespace DeleteProjectUserGroupResponse {
  export type AsObject = {
  }
}

export class GetCommandRequest extends jspb.Message {
  getCommandId(): string;
  setCommandId(value: string): GetCommandRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCommandRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetCommandRequest): GetCommandRequest.AsObject;
  static serializeBinaryToWriter(message: GetCommandRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCommandRequest;
  static deserializeBinaryFromReader(message: GetCommandRequest, reader: jspb.BinaryReader): GetCommandRequest;
}

export namespace GetCommandRequest {
  export type AsObject = {
    commandId: string,
  }
}

export class GetCommandResponse extends jspb.Message {
  getCommand(): pkg_model_command_pb.Command | undefined;
  setCommand(value?: pkg_model_command_pb.Command): GetCommandResponse;
  hasCommand(): boolean;
  clearCommand(): GetCommandResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetCommandResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetCommandResponse): GetCommandResponse.AsObject;
  static serializeBinaryToWriter(message: GetCommandResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetCommandResponse;
  static deserializeBinaryFromReader(message: GetCommandResponse, reader: jspb.BinaryReader): GetCommandResponse;
}

export namespace GetCommandResponse {
  export type AsObject = {
    command?: pkg_model_command_pb.Command.AsObject,
  }
}

export class GenerateAPIKeyRequest extends jspb.Message {
  getName(): string;
  setName(value: string): GenerateAPIKeyRequest;

  getRole(): pkg_model_apikey_pb.APIKey.Role;
  setRole(value: pkg_model_apikey_pb.APIKey.Role): GenerateAPIKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateAPIKeyRequest): GenerateAPIKeyRequest.AsObject;
  static serializeBinaryToWriter(message: GenerateAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateAPIKeyRequest;
  static deserializeBinaryFromReader(message: GenerateAPIKeyRequest, reader: jspb.BinaryReader): GenerateAPIKeyRequest;
}

export namespace GenerateAPIKeyRequest {
  export type AsObject = {
    name: string,
    role: pkg_model_apikey_pb.APIKey.Role,
  }
}

export class GenerateAPIKeyResponse extends jspb.Message {
  getKey(): string;
  setKey(value: string): GenerateAPIKeyResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GenerateAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GenerateAPIKeyResponse): GenerateAPIKeyResponse.AsObject;
  static serializeBinaryToWriter(message: GenerateAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GenerateAPIKeyResponse;
  static deserializeBinaryFromReader(message: GenerateAPIKeyResponse, reader: jspb.BinaryReader): GenerateAPIKeyResponse;
}

export namespace GenerateAPIKeyResponse {
  export type AsObject = {
    key: string,
  }
}

export class DisableAPIKeyRequest extends jspb.Message {
  getId(): string;
  setId(value: string): DisableAPIKeyRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAPIKeyRequest): DisableAPIKeyRequest.AsObject;
  static serializeBinaryToWriter(message: DisableAPIKeyRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyRequest;
  static deserializeBinaryFromReader(message: DisableAPIKeyRequest, reader: jspb.BinaryReader): DisableAPIKeyRequest;
}

export namespace DisableAPIKeyRequest {
  export type AsObject = {
    id: string,
  }
}

export class DisableAPIKeyResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DisableAPIKeyResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DisableAPIKeyResponse): DisableAPIKeyResponse.AsObject;
  static serializeBinaryToWriter(message: DisableAPIKeyResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DisableAPIKeyResponse;
  static deserializeBinaryFromReader(message: DisableAPIKeyResponse, reader: jspb.BinaryReader): DisableAPIKeyResponse;
}

export namespace DisableAPIKeyResponse {
  export type AsObject = {
  }
}

export class ListAPIKeysRequest extends jspb.Message {
  getOptions(): ListAPIKeysRequest.Options | undefined;
  setOptions(value?: ListAPIKeysRequest.Options): ListAPIKeysRequest;
  hasOptions(): boolean;
  clearOptions(): ListAPIKeysRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListAPIKeysRequest): ListAPIKeysRequest.AsObject;
  static serializeBinaryToWriter(message: ListAPIKeysRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysRequest;
  static deserializeBinaryFromReader(message: ListAPIKeysRequest, reader: jspb.BinaryReader): ListAPIKeysRequest;
}

export namespace ListAPIKeysRequest {
  export type AsObject = {
    options?: ListAPIKeysRequest.Options.AsObject,
  }

  export class Options extends jspb.Message {
    getEnabled(): google_protobuf_wrappers_pb.BoolValue | undefined;
    setEnabled(value?: google_protobuf_wrappers_pb.BoolValue): Options;
    hasEnabled(): boolean;
    clearEnabled(): Options;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
      enabled?: google_protobuf_wrappers_pb.BoolValue.AsObject,
    }
  }

}

export class ListAPIKeysResponse extends jspb.Message {
  getKeysList(): Array<pkg_model_apikey_pb.APIKey>;
  setKeysList(value: Array<pkg_model_apikey_pb.APIKey>): ListAPIKeysResponse;
  clearKeysList(): ListAPIKeysResponse;
  addKeys(value?: pkg_model_apikey_pb.APIKey, index?: number): pkg_model_apikey_pb.APIKey;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListAPIKeysResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListAPIKeysResponse): ListAPIKeysResponse.AsObject;
  static serializeBinaryToWriter(message: ListAPIKeysResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListAPIKeysResponse;
  static deserializeBinaryFromReader(message: ListAPIKeysResponse, reader: jspb.BinaryReader): ListAPIKeysResponse;
}

export namespace ListAPIKeysResponse {
  export type AsObject = {
    keysList: Array<pkg_model_apikey_pb.APIKey.AsObject>,
  }
}

export class GetInsightDataRequest extends jspb.Message {
  getMetricsKind(): pkg_model_insight_pb.InsightMetricsKind;
  setMetricsKind(value: pkg_model_insight_pb.InsightMetricsKind): GetInsightDataRequest;

  getRangeFrom(): number;
  setRangeFrom(value: number): GetInsightDataRequest;

  getRangeTo(): number;
  setRangeTo(value: number): GetInsightDataRequest;

  getResolution(): pkg_model_insight_pb.InsightResolution;
  setResolution(value: pkg_model_insight_pb.InsightResolution): GetInsightDataRequest;

  getApplicationId(): string;
  setApplicationId(value: string): GetInsightDataRequest;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): GetInsightDataRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInsightDataRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetInsightDataRequest): GetInsightDataRequest.AsObject;
  static serializeBinaryToWriter(message: GetInsightDataRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInsightDataRequest;
  static deserializeBinaryFromReader(message: GetInsightDataRequest, reader: jspb.BinaryReader): GetInsightDataRequest;
}

export namespace GetInsightDataRequest {
  export type AsObject = {
    metricsKind: pkg_model_insight_pb.InsightMetricsKind,
    rangeFrom: number,
    rangeTo: number,
    resolution: pkg_model_insight_pb.InsightResolution,
    applicationId: string,
    labelsMap: Array<[string, string]>,
  }
}

export class GetInsightDataResponse extends jspb.Message {
  getUpdatedAt(): number;
  setUpdatedAt(value: number): GetInsightDataResponse;

  getType(): pkg_model_insight_pb.InsightResultType;
  setType(value: pkg_model_insight_pb.InsightResultType): GetInsightDataResponse;

  getVectorList(): Array<pkg_model_insight_pb.InsightSample>;
  setVectorList(value: Array<pkg_model_insight_pb.InsightSample>): GetInsightDataResponse;
  clearVectorList(): GetInsightDataResponse;
  addVector(value?: pkg_model_insight_pb.InsightSample, index?: number): pkg_model_insight_pb.InsightSample;

  getMatrixList(): Array<pkg_model_insight_pb.InsightSampleStream>;
  setMatrixList(value: Array<pkg_model_insight_pb.InsightSampleStream>): GetInsightDataResponse;
  clearMatrixList(): GetInsightDataResponse;
  addMatrix(value?: pkg_model_insight_pb.InsightSampleStream, index?: number): pkg_model_insight_pb.InsightSampleStream;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInsightDataResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetInsightDataResponse): GetInsightDataResponse.AsObject;
  static serializeBinaryToWriter(message: GetInsightDataResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInsightDataResponse;
  static deserializeBinaryFromReader(message: GetInsightDataResponse, reader: jspb.BinaryReader): GetInsightDataResponse;
}

export namespace GetInsightDataResponse {
  export type AsObject = {
    updatedAt: number,
    type: pkg_model_insight_pb.InsightResultType,
    vectorList: Array<pkg_model_insight_pb.InsightSample.AsObject>,
    matrixList: Array<pkg_model_insight_pb.InsightSampleStream.AsObject>,
  }
}

export class GetInsightApplicationCountRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInsightApplicationCountRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetInsightApplicationCountRequest): GetInsightApplicationCountRequest.AsObject;
  static serializeBinaryToWriter(message: GetInsightApplicationCountRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInsightApplicationCountRequest;
  static deserializeBinaryFromReader(message: GetInsightApplicationCountRequest, reader: jspb.BinaryReader): GetInsightApplicationCountRequest;
}

export namespace GetInsightApplicationCountRequest {
  export type AsObject = {
  }
}

export class GetInsightApplicationCountResponse extends jspb.Message {
  getUpdatedAt(): number;
  setUpdatedAt(value: number): GetInsightApplicationCountResponse;

  getCountsList(): Array<pkg_model_insight_pb.InsightApplicationCount>;
  setCountsList(value: Array<pkg_model_insight_pb.InsightApplicationCount>): GetInsightApplicationCountResponse;
  clearCountsList(): GetInsightApplicationCountResponse;
  addCounts(value?: pkg_model_insight_pb.InsightApplicationCount, index?: number): pkg_model_insight_pb.InsightApplicationCount;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetInsightApplicationCountResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetInsightApplicationCountResponse): GetInsightApplicationCountResponse.AsObject;
  static serializeBinaryToWriter(message: GetInsightApplicationCountResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetInsightApplicationCountResponse;
  static deserializeBinaryFromReader(message: GetInsightApplicationCountResponse, reader: jspb.BinaryReader): GetInsightApplicationCountResponse;
}

export namespace GetInsightApplicationCountResponse {
  export type AsObject = {
    updatedAt: number,
    countsList: Array<pkg_model_insight_pb.InsightApplicationCount.AsObject>,
  }
}

export class ListDeploymentChainsRequest extends jspb.Message {
  getOptions(): ListDeploymentChainsRequest.Options | undefined;
  setOptions(value?: ListDeploymentChainsRequest.Options): ListDeploymentChainsRequest;
  hasOptions(): boolean;
  clearOptions(): ListDeploymentChainsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListDeploymentChainsRequest;

  getCursor(): string;
  setCursor(value: string): ListDeploymentChainsRequest;

  getPageMinUpdatedAt(): number;
  setPageMinUpdatedAt(value: number): ListDeploymentChainsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentChainsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentChainsRequest): ListDeploymentChainsRequest.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentChainsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentChainsRequest;
  static deserializeBinaryFromReader(message: ListDeploymentChainsRequest, reader: jspb.BinaryReader): ListDeploymentChainsRequest;
}

export namespace ListDeploymentChainsRequest {
  export type AsObject = {
    options?: ListDeploymentChainsRequest.Options.AsObject,
    pageSize: number,
    cursor: string,
    pageMinUpdatedAt: number,
  }

  export class Options extends jspb.Message {
    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
    }
  }

}

export class ListDeploymentChainsResponse extends jspb.Message {
  getDeploymentChainsList(): Array<pkg_model_deployment_chain_pb.DeploymentChain>;
  setDeploymentChainsList(value: Array<pkg_model_deployment_chain_pb.DeploymentChain>): ListDeploymentChainsResponse;
  clearDeploymentChainsList(): ListDeploymentChainsResponse;
  addDeploymentChains(value?: pkg_model_deployment_chain_pb.DeploymentChain, index?: number): pkg_model_deployment_chain_pb.DeploymentChain;

  getCursor(): string;
  setCursor(value: string): ListDeploymentChainsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListDeploymentChainsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListDeploymentChainsResponse): ListDeploymentChainsResponse.AsObject;
  static serializeBinaryToWriter(message: ListDeploymentChainsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListDeploymentChainsResponse;
  static deserializeBinaryFromReader(message: ListDeploymentChainsResponse, reader: jspb.BinaryReader): ListDeploymentChainsResponse;
}

export namespace ListDeploymentChainsResponse {
  export type AsObject = {
    deploymentChainsList: Array<pkg_model_deployment_chain_pb.DeploymentChain.AsObject>,
    cursor: string,
  }
}

export class GetDeploymentChainRequest extends jspb.Message {
  getDeploymentChainId(): string;
  setDeploymentChainId(value: string): GetDeploymentChainRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDeploymentChainRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetDeploymentChainRequest): GetDeploymentChainRequest.AsObject;
  static serializeBinaryToWriter(message: GetDeploymentChainRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDeploymentChainRequest;
  static deserializeBinaryFromReader(message: GetDeploymentChainRequest, reader: jspb.BinaryReader): GetDeploymentChainRequest;
}

export namespace GetDeploymentChainRequest {
  export type AsObject = {
    deploymentChainId: string,
  }
}

export class GetDeploymentChainResponse extends jspb.Message {
  getDeploymentChain(): pkg_model_deployment_chain_pb.DeploymentChain | undefined;
  setDeploymentChain(value?: pkg_model_deployment_chain_pb.DeploymentChain): GetDeploymentChainResponse;
  hasDeploymentChain(): boolean;
  clearDeploymentChain(): GetDeploymentChainResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetDeploymentChainResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetDeploymentChainResponse): GetDeploymentChainResponse.AsObject;
  static serializeBinaryToWriter(message: GetDeploymentChainResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetDeploymentChainResponse;
  static deserializeBinaryFromReader(message: GetDeploymentChainResponse, reader: jspb.BinaryReader): GetDeploymentChainResponse;
}

export namespace GetDeploymentChainResponse {
  export type AsObject = {
    deploymentChain?: pkg_model_deployment_chain_pb.DeploymentChain.AsObject,
  }
}

export class ListEventsRequest extends jspb.Message {
  getOptions(): ListEventsRequest.Options | undefined;
  setOptions(value?: ListEventsRequest.Options): ListEventsRequest;
  hasOptions(): boolean;
  clearOptions(): ListEventsRequest;

  getPageSize(): number;
  setPageSize(value: number): ListEventsRequest;

  getCursor(): string;
  setCursor(value: string): ListEventsRequest;

  getPageMinUpdatedAt(): number;
  setPageMinUpdatedAt(value: number): ListEventsRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventsRequest.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventsRequest): ListEventsRequest.AsObject;
  static serializeBinaryToWriter(message: ListEventsRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventsRequest;
  static deserializeBinaryFromReader(message: ListEventsRequest, reader: jspb.BinaryReader): ListEventsRequest;
}

export namespace ListEventsRequest {
  export type AsObject = {
    options?: ListEventsRequest.Options.AsObject,
    pageSize: number,
    cursor: string,
    pageMinUpdatedAt: number,
  }

  export class Options extends jspb.Message {
    getName(): string;
    setName(value: string): Options;

    getStatusesList(): Array<pkg_model_event_pb.EventStatus>;
    setStatusesList(value: Array<pkg_model_event_pb.EventStatus>): Options;
    clearStatusesList(): Options;
    addStatuses(value: pkg_model_event_pb.EventStatus, index?: number): Options;

    getLabelsMap(): jspb.Map<string, string>;
    clearLabelsMap(): Options;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Options.AsObject;
    static toObject(includeInstance: boolean, msg: Options): Options.AsObject;
    static serializeBinaryToWriter(message: Options, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Options;
    static deserializeBinaryFromReader(message: Options, reader: jspb.BinaryReader): Options;
  }

  export namespace Options {
    export type AsObject = {
      name: string,
      statusesList: Array<pkg_model_event_pb.EventStatus>,
      labelsMap: Array<[string, string]>,
    }
  }

}

export class ListEventsResponse extends jspb.Message {
  getEventsList(): Array<pkg_model_event_pb.Event>;
  setEventsList(value: Array<pkg_model_event_pb.Event>): ListEventsResponse;
  clearEventsList(): ListEventsResponse;
  addEvents(value?: pkg_model_event_pb.Event, index?: number): pkg_model_event_pb.Event;

  getCursor(): string;
  setCursor(value: string): ListEventsResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ListEventsResponse.AsObject;
  static toObject(includeInstance: boolean, msg: ListEventsResponse): ListEventsResponse.AsObject;
  static serializeBinaryToWriter(message: ListEventsResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ListEventsResponse;
  static deserializeBinaryFromReader(message: ListEventsResponse, reader: jspb.BinaryReader): ListEventsResponse;
}

export namespace ListEventsResponse {
  export type AsObject = {
    eventsList: Array<pkg_model_event_pb.Event.AsObject>,
    cursor: string,
  }
}

