import * as jspb from 'google-protobuf'


import * as google_protobuf_descriptor_pb from 'google-protobuf/google/protobuf/descriptor_pb';


export class Project extends jspb.Message {
  getId(): string;
  setId(value: string): Project;

  getDesc(): string;
  setDesc(value: string): Project;

  getStaticAdmin(): ProjectStaticUser | undefined;
  setStaticAdmin(value?: ProjectStaticUser): Project;
  hasStaticAdmin(): boolean;
  clearStaticAdmin(): Project;

  getStaticAdminDisabled(): boolean;
  setStaticAdminDisabled(value: boolean): Project;

  getSso(): ProjectSSOConfig | undefined;
  setSso(value?: ProjectSSOConfig): Project;
  hasSso(): boolean;
  clearSso(): Project;

  getRbac(): ProjectRBACConfig | undefined;
  setRbac(value?: ProjectRBACConfig): Project;
  hasRbac(): boolean;
  clearRbac(): Project;

  getSharedSsoName(): string;
  setSharedSsoName(value: string): Project;

  getAllowStrayAsViewer(): boolean;
  setAllowStrayAsViewer(value: boolean): Project;

  getRbacRolesList(): Array<ProjectRBACRole>;
  setRbacRolesList(value: Array<ProjectRBACRole>): Project;
  clearRbacRolesList(): Project;
  addRbacRoles(value?: ProjectRBACRole, index?: number): ProjectRBACRole;

  getUserGroupsList(): Array<ProjectUserGroup>;
  setUserGroupsList(value: Array<ProjectUserGroup>): Project;
  clearUserGroupsList(): Project;
  addUserGroups(value?: ProjectUserGroup, index?: number): ProjectUserGroup;

  getCreatedAt(): number;
  setCreatedAt(value: number): Project;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): Project;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Project.AsObject;
  static toObject(includeInstance: boolean, msg: Project): Project.AsObject;
  static serializeBinaryToWriter(message: Project, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Project;
  static deserializeBinaryFromReader(message: Project, reader: jspb.BinaryReader): Project;
}

export namespace Project {
  export type AsObject = {
    id: string,
    desc: string,
    staticAdmin?: ProjectStaticUser.AsObject,
    staticAdminDisabled: boolean,
    sso?: ProjectSSOConfig.AsObject,
    rbac?: ProjectRBACConfig.AsObject,
    sharedSsoName: string,
    allowStrayAsViewer: boolean,
    rbacRolesList: Array<ProjectRBACRole.AsObject>,
    userGroupsList: Array<ProjectUserGroup.AsObject>,
    createdAt: number,
    updatedAt: number,
  }
}

export class ProjectStaticUser extends jspb.Message {
  getUsername(): string;
  setUsername(value: string): ProjectStaticUser;

  getPasswordHash(): string;
  setPasswordHash(value: string): ProjectStaticUser;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectStaticUser.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectStaticUser): ProjectStaticUser.AsObject;
  static serializeBinaryToWriter(message: ProjectStaticUser, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectStaticUser;
  static deserializeBinaryFromReader(message: ProjectStaticUser, reader: jspb.BinaryReader): ProjectStaticUser;
}

export namespace ProjectStaticUser {
  export type AsObject = {
    username: string,
    passwordHash: string,
  }
}

export class ProjectSSOConfig extends jspb.Message {
  getProvider(): ProjectSSOConfig.Provider;
  setProvider(value: ProjectSSOConfig.Provider): ProjectSSOConfig;

  getSessionTtl(): number;
  setSessionTtl(value: number): ProjectSSOConfig;

  getGithub(): ProjectSSOConfig.GitHub | undefined;
  setGithub(value?: ProjectSSOConfig.GitHub): ProjectSSOConfig;
  hasGithub(): boolean;
  clearGithub(): ProjectSSOConfig;

  getGoogle(): ProjectSSOConfig.Google | undefined;
  setGoogle(value?: ProjectSSOConfig.Google): ProjectSSOConfig;
  hasGoogle(): boolean;
  clearGoogle(): ProjectSSOConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectSSOConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectSSOConfig): ProjectSSOConfig.AsObject;
  static serializeBinaryToWriter(message: ProjectSSOConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectSSOConfig;
  static deserializeBinaryFromReader(message: ProjectSSOConfig, reader: jspb.BinaryReader): ProjectSSOConfig;
}

export namespace ProjectSSOConfig {
  export type AsObject = {
    provider: ProjectSSOConfig.Provider,
    sessionTtl: number,
    github?: ProjectSSOConfig.GitHub.AsObject,
    google?: ProjectSSOConfig.Google.AsObject,
  }

  export class GitHub extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): GitHub;

    getClientSecret(): string;
    setClientSecret(value: string): GitHub;

    getBaseUrl(): string;
    setBaseUrl(value: string): GitHub;

    getUploadUrl(): string;
    setUploadUrl(value: string): GitHub;

    getProxyUrl(): string;
    setProxyUrl(value: string): GitHub;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GitHub.AsObject;
    static toObject(includeInstance: boolean, msg: GitHub): GitHub.AsObject;
    static serializeBinaryToWriter(message: GitHub, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GitHub;
    static deserializeBinaryFromReader(message: GitHub, reader: jspb.BinaryReader): GitHub;
  }

  export namespace GitHub {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
      baseUrl: string,
      uploadUrl: string,
      proxyUrl: string,
    }
  }


  export class Google extends jspb.Message {
    getClientId(): string;
    setClientId(value: string): Google;

    getClientSecret(): string;
    setClientSecret(value: string): Google;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Google.AsObject;
    static toObject(includeInstance: boolean, msg: Google): Google.AsObject;
    static serializeBinaryToWriter(message: Google, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Google;
    static deserializeBinaryFromReader(message: Google, reader: jspb.BinaryReader): Google;
  }

  export namespace Google {
    export type AsObject = {
      clientId: string,
      clientSecret: string,
    }
  }


  export enum Provider { 
    GITHUB = 0,
    GOOGLE = 2,
  }
}

export class ProjectRBACConfig extends jspb.Message {
  getAdmin(): string;
  setAdmin(value: string): ProjectRBACConfig;

  getEditor(): string;
  setEditor(value: string): ProjectRBACConfig;

  getViewer(): string;
  setViewer(value: string): ProjectRBACConfig;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRBACConfig.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRBACConfig): ProjectRBACConfig.AsObject;
  static serializeBinaryToWriter(message: ProjectRBACConfig, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRBACConfig;
  static deserializeBinaryFromReader(message: ProjectRBACConfig, reader: jspb.BinaryReader): ProjectRBACConfig;
}

export namespace ProjectRBACConfig {
  export type AsObject = {
    admin: string,
    editor: string,
    viewer: string,
  }
}

export class ProjectUserGroup extends jspb.Message {
  getSsoGroup(): string;
  setSsoGroup(value: string): ProjectUserGroup;

  getRole(): string;
  setRole(value: string): ProjectUserGroup;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectUserGroup.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectUserGroup): ProjectUserGroup.AsObject;
  static serializeBinaryToWriter(message: ProjectUserGroup, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectUserGroup;
  static deserializeBinaryFromReader(message: ProjectUserGroup, reader: jspb.BinaryReader): ProjectUserGroup;
}

export namespace ProjectUserGroup {
  export type AsObject = {
    ssoGroup: string,
    role: string,
  }
}

export class ProjectRBACRole extends jspb.Message {
  getName(): string;
  setName(value: string): ProjectRBACRole;

  getPoliciesList(): Array<ProjectRBACPolicy>;
  setPoliciesList(value: Array<ProjectRBACPolicy>): ProjectRBACRole;
  clearPoliciesList(): ProjectRBACRole;
  addPolicies(value?: ProjectRBACPolicy, index?: number): ProjectRBACPolicy;

  getIsBuiltin(): boolean;
  setIsBuiltin(value: boolean): ProjectRBACRole;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRBACRole.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRBACRole): ProjectRBACRole.AsObject;
  static serializeBinaryToWriter(message: ProjectRBACRole, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRBACRole;
  static deserializeBinaryFromReader(message: ProjectRBACRole, reader: jspb.BinaryReader): ProjectRBACRole;
}

export namespace ProjectRBACRole {
  export type AsObject = {
    name: string,
    policiesList: Array<ProjectRBACPolicy.AsObject>,
    isBuiltin: boolean,
  }
}

export class ProjectRBACResource extends jspb.Message {
  getType(): ProjectRBACResource.ResourceType;
  setType(value: ProjectRBACResource.ResourceType): ProjectRBACResource;

  getLabelsMap(): jspb.Map<string, string>;
  clearLabelsMap(): ProjectRBACResource;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRBACResource.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRBACResource): ProjectRBACResource.AsObject;
  static serializeBinaryToWriter(message: ProjectRBACResource, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRBACResource;
  static deserializeBinaryFromReader(message: ProjectRBACResource, reader: jspb.BinaryReader): ProjectRBACResource;
}

export namespace ProjectRBACResource {
  export type AsObject = {
    type: ProjectRBACResource.ResourceType,
    labelsMap: Array<[string, string]>,
  }

  export enum ResourceType { 
    ALL = 0,
    APPLICATION = 1,
    DEPLOYMENT = 2,
    EVENT = 3,
    PIPED = 4,
    DEPLOYMENT_CHAIN = 5,
    PROJECT = 6,
    API_KEY = 7,
    INSIGHT = 8,
  }
}

export class ProjectRBACPolicy extends jspb.Message {
  getResourcesList(): Array<ProjectRBACResource>;
  setResourcesList(value: Array<ProjectRBACResource>): ProjectRBACPolicy;
  clearResourcesList(): ProjectRBACPolicy;
  addResources(value?: ProjectRBACResource, index?: number): ProjectRBACResource;

  getActionsList(): Array<ProjectRBACPolicy.Action>;
  setActionsList(value: Array<ProjectRBACPolicy.Action>): ProjectRBACPolicy;
  clearActionsList(): ProjectRBACPolicy;
  addActions(value: ProjectRBACPolicy.Action, index?: number): ProjectRBACPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRBACPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRBACPolicy): ProjectRBACPolicy.AsObject;
  static serializeBinaryToWriter(message: ProjectRBACPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRBACPolicy;
  static deserializeBinaryFromReader(message: ProjectRBACPolicy, reader: jspb.BinaryReader): ProjectRBACPolicy;
}

export namespace ProjectRBACPolicy {
  export type AsObject = {
    resourcesList: Array<ProjectRBACResource.AsObject>,
    actionsList: Array<ProjectRBACPolicy.Action>,
  }

  export enum Action { 
    ALL = 0,
    GET = 1,
    LIST = 2,
    CREATE = 3,
    UPDATE = 4,
    DELETE = 5,
  }
}

