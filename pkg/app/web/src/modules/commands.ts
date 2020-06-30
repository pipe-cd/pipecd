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

export type Command = CommandModel.AsObject;

const commandsAdapter = createEntityAdapter<Command>();
export const fetchCommand = createAsyncThunk(
  "commands/fetch",
  async (commandId: string) => {
    const { command } = await getCommand({ commandId });
    if (command === undefined) {
      throw Error("command not found");
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
        commandsAdapter.addOne(state, action.payload);
      }
    });
  },
});

export {
  CommandStatus,
  Command as CommandModel,
} from "pipe/pkg/app/web/model/command_pb";
