import * as jspb from 'google-protobuf'




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
    GITHUB_ENTERPRISE = 1,
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

export class ProjectRBACRole extends jspb.Message {
  getName(): string;
  setName(value: string): ProjectRBACRole;

  getType(): ProjectRBACRole.RoleType;
  setType(value: ProjectRBACRole.RoleType): ProjectRBACRole;

  getSubject(): string;
  setSubject(value: string): ProjectRBACRole;

  getPolicyList(): Array<ProjectRBACPolicy>;
  setPolicyList(value: Array<ProjectRBACPolicy>): ProjectRBACRole;
  clearPolicyList(): ProjectRBACRole;
  addPolicy(value?: ProjectRBACPolicy, index?: number): ProjectRBACPolicy;

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
    type: ProjectRBACRole.RoleType,
    subject: string,
    policyList: Array<ProjectRBACPolicy.AsObject>,
  }

  export enum RoleType { 
    GITHUB_TEAM = 0,
    GOOGLE_GROUP = 1,
  }
}

export class ProjectRBACPolicy extends jspb.Message {
  getResourceType(): ProjectRBACPolicy.ResourceType;
  setResourceType(value: ProjectRBACPolicy.ResourceType): ProjectRBACPolicy;

  getActionType(): ProjectRBACPolicy.ActionType;
  setActionType(value: ProjectRBACPolicy.ActionType): ProjectRBACPolicy;

  getLabelKey(): string;
  setLabelKey(value: string): ProjectRBACPolicy;

  getLabelValue(): string;
  setLabelValue(value: string): ProjectRBACPolicy;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ProjectRBACPolicy.AsObject;
  static toObject(includeInstance: boolean, msg: ProjectRBACPolicy): ProjectRBACPolicy.AsObject;
  static serializeBinaryToWriter(message: ProjectRBACPolicy, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ProjectRBACPolicy;
  static deserializeBinaryFromReader(message: ProjectRBACPolicy, reader: jspb.BinaryReader): ProjectRBACPolicy;
}

export namespace ProjectRBACPolicy {
  export type AsObject = {
    resourceType: ProjectRBACPolicy.ResourceType,
    actionType: ProjectRBACPolicy.ActionType,
    labelKey: string,
    labelValue: string,
  }

  export enum ResourceType { 
    APPLICATION = 0,
    DEPLOYMENT = 1,
    EVENT = 2,
    PIPED = 3,
    DEPLOYMENTCHAIN = 4,
    PROJECT = 5,
    ROLE = 6,
    APIKEY = 7,
  }

  export enum ActionType { 
    GET = 0,
    LIST = 1,
    CREATE = 2,
    UPDATE = 3,
    DELETE = 4,
  }
}

