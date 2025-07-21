import {
  GetCommandRequest,
  GetCommandResponse,
} from "pipecd/web/api_client/service_pb";
import { Command } from "~/types/commands";
import {
  dummyCommand,
  dummySyncSucceededCommand,
} from "~/__fixtures__/dummy-command";
import { createHandler } from "../create-handler";

const createCommandModel = (commandObj: Command.AsObject): Command => {
  const command = new Command();
  command.setId(commandObj.id);
  command.setApplicationId(commandObj.applicationId);
  command.setPipedId(commandObj.pipedId);
  command.setDeploymentId(commandObj.deploymentId);
  command.setStageId(commandObj.stageId);
  command.setCommander(commandObj.commander);
  command.setStatus(commandObj.status);
  command.setHandledAt(commandObj.handledAt);
  command.setType(commandObj.type);
  command.setCreatedAt(commandObj.createdAt);
  command.setUpdatedAt(commandObj.updatedAt);
  commandObj.metadataMap.forEach(([key, value]) => {
    command.getMetadataMap().set(key, value);
  });
  return command;
};

export const commandHandlers = [
  createHandler<GetCommandResponse>("/GetCommand", (requestBody) => {
    const response = new GetCommandResponse();
    const request = GetCommandRequest.deserializeBinary(requestBody);

    if (request.getCommandId() === dummySyncSucceededCommand.id) {
      response.setCommand(createCommandModel(dummySyncSucceededCommand));
    } else {
      response.setCommand(createCommandModel(dummyCommand));
    }

    return response;
  }),
];
