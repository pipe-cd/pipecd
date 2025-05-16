import { SyncStrategy } from "pipecd/web/model/common_pb";
import { Command, CommandStatus } from "~/modules/commands";
import { dummyDeployment } from "./dummy-deployment";
import { createRandTimes } from "./utils";

const [createdAt, handledAt] = createRandTimes(3);

export const dummyCommand: Command.AsObject = {
  id: "command-1",
  pipedId: "piped-1",
  applicationId: "app-1",
  projectId: "project-1",
  deploymentId: "",
  stageId: "",
  commander: "user",
  status: CommandStatus.COMMAND_NOT_HANDLED_YET,
  metadataMap: [],
  type: Command.Type.SYNC_APPLICATION,
  syncApplication: {
    applicationId: "app-1",
    syncStrategy: SyncStrategy.AUTO,
  },
  createdAt: createdAt.unix(),
  updatedAt: handledAt.unix(),
  handledAt: handledAt.unix(),
  errorReason: "",
};

export const dummySyncSucceededCommand: Command.AsObject = {
  ...dummyCommand,
  id: "sync-succeeded",
  status: CommandStatus.COMMAND_SUCCEEDED,
  metadataMap: [["TriggeredDeploymentID", dummyDeployment.id]],
};
