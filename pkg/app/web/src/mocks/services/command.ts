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
  rest.post<Uint8Array>(createMask("/GetCommand"), (req, res, ctx) => {
    const response = new GetCommandResponse();

    const request = GetCommandRequest.deserializeBinary(req.body.slice(5));

    if (request.getCommandId() === dummySyncSucceededCommand.id) {
      response.setCommand(createCommandModel(dummySyncSucceededCommand));
    } else {
      response.setCommand(createCommandModel(dummyCommand));
    }

    const data = serialize(response.serializeBinary());
    return res(
      ctx.status(200),
      ctx.set("Content-Type", "application/grpc-web+proto"),
      ctx.body(data)
    );
  }),
];
