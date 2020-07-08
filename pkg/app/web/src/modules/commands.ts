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
import { addToast } from "./toasts";

export type Command = CommandModel.AsObject;

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

    if (command.status === CommandStatus.COMMAND_SUCCEEDED) {
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
