import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import {
  Command as CommandModel,
  CommandStatus,
} from "pipe/pkg/app/web/model/command_pb";
import { getCommand } from "../api/commands";
import { PAGE_PATH_DEPLOYMENTS } from "../constants/path";
import { findMetadataByKey } from "../utils/find-metadata-by-key";
import { addToast } from "./toasts";

export type Command = CommandModel.AsObject;

const METADATA_KEY = {
  TRIGGERED_DEPLOYMENT_ID: "TriggeredDeploymentID",
};

export const COMMAND_TYPE_TEXT: Record<CommandModel.Type, string> = {
  [CommandModel.Type.APPROVE_STAGE]: "Approve Stage",
  [CommandModel.Type.CANCEL_DEPLOYMENT]: "Cancel Deployment",
  [CommandModel.Type.SYNC_APPLICATION]: "Sync Application",
  [CommandModel.Type.UPDATE_APPLICATION_CONFIG]: "Update Application Config",
};

const commandsAdapter = createEntityAdapter<Command>();
export const fetchCommand = createAsyncThunk(
  "commands/fetch",
  async (commandId: string, thunkAPI) => {
    const { command } = await getCommand({ commandId });
    if (command === undefined) {
      throw Error("command not found");
    }

    if (command.status !== CommandStatus.COMMAND_SUCCEEDED) {
      return command;
    }

    switch (command.type) {
      case CommandModel.Type.SYNC_APPLICATION: {
        const deploymentId = findMetadataByKey(
          command.metadataMap,
          METADATA_KEY.TRIGGERED_DEPLOYMENT_ID
        );
        thunkAPI.dispatch(
          addToast({
            message: `Succeed "${COMMAND_TYPE_TEXT[command.type]}"`,
            severity: "success",
            to: deploymentId
              ? `${PAGE_PATH_DEPLOYMENTS}/${deploymentId}`
              : undefined,
          })
        );
        break;
      }
      default:
        thunkAPI.dispatch(
          addToast({
            message: `Succeed "${COMMAND_TYPE_TEXT[command.type]}"`,
            severity: "success",
          })
        );
    }

    return command;
  }
);

export const { selectIds } = commandsAdapter.getSelectors();

export const commandsSlice = createSlice({
  name: "commands",
  initialState: commandsAdapter.getInitialState(),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchCommand.fulfilled, (state, action) => {
      // If command process is finished, remove from processing command ids.
      if (action.payload.status !== CommandStatus.COMMAND_NOT_HANDLED_YET) {
        commandsAdapter.removeOne(state, action.payload.id);
      } else {
        commandsAdapter.upsertOne(state, action.payload);
      }
    });
  },
});

export {
  CommandStatus,
  Command as CommandModel,
} from "pipe/pkg/app/web/model/command_pb";
