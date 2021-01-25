import { SyncStrategy } from "pipe/pkg/app/web/model/deployment_pb";
import { Command, CommandStatus } from "../modules/commands";
import { dummyDeployment } from "./dummy-deployment";

export const dummyCommand: Command.AsObject = {
  id: "command-1",
  pipedId: "piped-1",
  applicationId: "app-1",
  deploymentId: "",
  stageId: "",
  commander: "user",
  status: CommandStatus.COMMAND_NOT_HANDLED_YET,
  metadataMap: [],
  handledAt: 0,
  type: Command.Type.SYNC_APPLICATION,
  syncApplication: {
    applicationId: "app-1",
    syncStrategy: SyncStrategy.AUTO,
  },
  createdAt: 0,
  updatedAt: 0,
};

export const dummySyncSucceededCommand: Command.AsObject = {
  ...dummyCommand,
  id: "sync-succeeded",
  status: CommandStatus.COMMAND_SUCCEEDED,
  metadataMap: [["TriggeredDeploymentID", dummyDeployment.id]],
};
