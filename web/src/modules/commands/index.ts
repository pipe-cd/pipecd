import {
  createAsyncThunk,
  createEntityAdapter,
  createSlice,
} from "@reduxjs/toolkit";
import { Command, CommandStatus } from "pipecd/web/model/command_pb";
import { getCommand } from "~/api/commands";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";
import { addToast } from "../toasts";

const METADATA_KEY = {
  TRIGGERED_DEPLOYMENT_ID: "TriggeredDeploymentID",
};

export const COMMAND_TYPE_TEXT: Record<Command.Type, string> = {
  [Command.Type.APPROVE_STAGE]: "Approve Stage",
  [Command.Type.CANCEL_DEPLOYMENT]: "Cancel Deployment",
  [Command.Type.SYNC_APPLICATION]: "Sync Application",
  [Command.Type.UPDATE_APPLICATION_CONFIG]: "Update Application Config",
  [Command.Type.BUILD_PLAN_PREVIEW]: "Build Plan Preview",
  [Command.Type.CHAIN_SYNC_APPLICATION]: "Chain Sync Application",
  [Command.Type.SKIP_STAGE]: "Skip Stage",
  [Command.Type.RESTART_PIPED]: "Restart Piped",
};

const commandsAdapter = createEntityAdapter<Command.AsObject>();
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
      case Command.Type.SYNC_APPLICATION: {
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

export { CommandStatus, Command } from "pipecd/web/model/command_pb";
