import * as jspb from 'google-protobuf'


import * as pkg_model_deployment_pb from 'pipecd/pkg/app/web/model/deployment_pb';


export class DeploymentChain extends jspb.Message {
  getId(): string;
  setId(value: string): DeploymentChain;

  getProjectId(): string;
  setProjectId(value: string): DeploymentChain;

  getStatus(): ChainStatus;
  setStatus(value: ChainStatus): DeploymentChain;

  getBlocksList(): Array<ChainBlock>;
  setBlocksList(value: Array<ChainBlock>): DeploymentChain;
  clearBlocksList(): DeploymentChain;
  addBlocks(value?: ChainBlock, index?: number): ChainBlock;

  getCompletedAt(): number;
  setCompletedAt(value: number): DeploymentChain;

  getCreatedAt(): number;
  setCreatedAt(value: number): DeploymentChain;

  getUpdatedAt(): number;
  setUpdatedAt(value: number): DeploymentChain;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeploymentChain.AsObject;
  static toObject(includeInstance: boolean, msg: DeploymentChain): DeploymentChain.AsObject;
  static serializeBinaryToWriter(message: DeploymentChain, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeploymentChain;
  static deserializeBinaryFromReader(message: DeploymentChain, reader: jspb.BinaryReader): DeploymentChain;
}

export namespace DeploymentChain {
  export type AsObject = {
    id: string,
    projectId: string,
    status: ChainStatus,
    blocksList: Array<ChainBlock.AsObject>,
    completedAt: number,
    createdAt: number,
    updatedAt: number,
  }
}

export class ChainApplicationRef extends jspb.Message {
  getApplicationId(): string;
  setApplicationId(value: string): ChainApplicationRef;

  getApplicationName(): string;
  setApplicationName(value: string): ChainApplicationRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChainApplicationRef.AsObject;
  static toObject(includeInstance: boolean, msg: ChainApplicationRef): ChainApplicationRef.AsObject;
  static serializeBinaryToWriter(message: ChainApplicationRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChainApplicationRef;
  static deserializeBinaryFromReader(message: ChainApplicationRef, reader: jspb.BinaryReader): ChainApplicationRef;
}

export namespace ChainApplicationRef {
  export type AsObject = {
    applicationId: string,
    applicationName: string,
  }
}

export class ChainDeploymentRef extends jspb.Message {
  getDeploymentId(): string;
  setDeploymentId(value: string): ChainDeploymentRef;

  getStatus(): pkg_model_deployment_pb.DeploymentStatus;
  setStatus(value: pkg_model_deployment_pb.DeploymentStatus): ChainDeploymentRef;

  getStatusReason(): string;
  setStatusReason(value: string): ChainDeploymentRef;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChainDeploymentRef.AsObject;
  static toObject(includeInstance: boolean, msg: ChainDeploymentRef): ChainDeploymentRef.AsObject;
  static serializeBinaryToWriter(message: ChainDeploymentRef, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChainDeploymentRef;
  static deserializeBinaryFromReader(message: ChainDeploymentRef, reader: jspb.BinaryReader): ChainDeploymentRef;
}

export namespace ChainDeploymentRef {
  export type AsObject = {
    deploymentId: string,
    status: pkg_model_deployment_pb.DeploymentStatus,
    statusReason: string,
  }
}

export class ChainNode extends jspb.Message {
  getApplicationRef(): ChainApplicationRef | undefined;
  setApplicationRef(value?: ChainApplicationRef): ChainNode;
  hasApplicationRef(): boolean;
  clearApplicationRef(): ChainNode;

  getDeploymentRef(): ChainDeploymentRef | undefined;
  setDeploymentRef(value?: ChainDeploymentRef): ChainNode;
  hasDeploymentRef(): boolean;
  clearDeploymentRef(): ChainNode;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChainNode.AsObject;
  static toObject(includeInstance: boolean, msg: ChainNode): ChainNode.AsObject;
  static serializeBinaryToWriter(message: ChainNode, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChainNode;
  static deserializeBinaryFromReader(message: ChainNode, reader: jspb.BinaryReader): ChainNode;
}

export namespace ChainNode {
  export type AsObject = {
    applicationRef?: ChainApplicationRef.AsObject,
    deploymentRef?: ChainDeploymentRef.AsObject,
  }
}

export class ChainBlock extends jspb.Message {
  getNodesList(): Array<ChainNode>;
  setNodesList(value: Array<ChainNode>): ChainBlock;
  clearNodesList(): ChainBlock;
  addNodes(value?: ChainNode, index?: number): ChainNode;

  getStatus(): ChainBlockStatus;
  setStatus(value: ChainBlockStatus): ChainBlock;

  getStartedAt(): number;
  setStartedAt(value: number): ChainBlock;

  getCompletedAt(): number;
  setCompletedAt(value: number): ChainBlock;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): ChainBlock.AsObject;
  static toObject(includeInstance: boolean, msg: ChainBlock): ChainBlock.AsObject;
  static serializeBinaryToWriter(message: ChainBlock, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): ChainBlock;
  static deserializeBinaryFromReader(message: ChainBlock, reader: jspb.BinaryReader): ChainBlock;
}

export namespace ChainBlock {
  export type AsObject = {
    nodesList: Array<ChainNode.AsObject>,
    status: ChainBlockStatus,
    startedAt: number,
    completedAt: number,
  }
}

export enum ChainStatus { 
  DEPLOYMENT_CHAIN_PENDING = 0,
  DEPLOYMENT_CHAIN_RUNNING = 1,
  DEPLOYMENT_CHAIN_SUCCESS = 2,
  DEPLOYMENT_CHAIN_FAILURE = 3,
  DEPLOYMENT_CHAIN_CANCELLED = 4,
}
export enum ChainBlockStatus { 
  DEPLOYMENT_BLOCK_PENDING = 0,
  DEPLOYMENT_BLOCK_RUNNING = 1,
  DEPLOYMENT_BLOCK_SUCCESS = 2,
  DEPLOYMENT_BLOCK_FAILURE = 3,
  DEPLOYMENT_BLOCK_CANCELLED = 4,
}
