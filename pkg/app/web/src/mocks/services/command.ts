import { rest } from "msw";
import {
  GetCommandRequest,
  GetCommandResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { CommandModel } from "../../modules/commands";
import {
  dummyCommand,
  dummySyncSucceededCommand,
} from "../../__fixtures__/dummy-command";
import { serialize } from "../serializer";
import { createMask } from "../utils";
import { createHandler } from "../create-handler";

const createCommandModel = (
  commandObj: CommandModel.AsObject
): CommandModel => {
  const command = new CommandModel();
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
