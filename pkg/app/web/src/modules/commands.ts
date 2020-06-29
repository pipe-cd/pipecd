import {
  createSlice,
  createEntityAdapter,
  createAsyncThunk,
} from "@reduxjs/toolkit";
import { Command as CommandModel } from "pipe/pkg/app/web/model/command_pb";
import { getCommand } from "../api/commands";

export type Command = CommandModel.AsObject;

const commandsAdapter = createEntityAdapter<Command>();
const fetchCommand = createAsyncThunk(
  "commands/fetch",
  async (commandId: string) => {
    const { command } = await getCommand({ commandId });
    if (command === undefined) {
      throw Error("command not found");
    }
    return command;
  }
);

export const commandsSlice = createSlice({
  name: "commands",
  initialState: commandsAdapter.getInitialState(),
  reducers: {},
  extraReducers: (builder) => {
    builder.addCase(fetchCommand.fulfilled, (state, action) => {
      commandsAdapter.addOne(state, action.payload);
    });
  },
});

export { CommandStatus } from "pipe/pkg/app/web/model/command_pb";
