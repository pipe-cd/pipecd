import { Command, CommandModel, CommandStatus } from "../modules/commands";
import { dummyDeployment } from "./dummy-deployment";

export const dummyCommand: Command = {
  id: "command-1",
  pipedId: "piped-1",
  applicationId: "app-1",
  deploymentId: "",
  stageId: "",
  commander: "user",
  status: CommandStatus.COMMAND_NOT_HANDLED_YET,
  metadataMap: [],
  handledAt: 0,
  type: CommandModel.Type.SYNC_APPLICATION,
  syncApplication: {
    applicationId: "app-1",
  },
  createdAt: 0,
  updatedAt: 0,
};

export const dummySyncSucceededCommand: Command = {
  ...dummyCommand,
  id: "sync-succeeded",
  status: CommandStatus.COMMAND_SUCCEEDED,
  metadataMap: [["TriggeredDeploymentID", dummyDeployment.id]],
};
