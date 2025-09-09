import React from "react";
import { Command } from "~~/model/command_pb";

export type CommandContextType = {
  fetchedCommands: Record<string, Command.AsObject>;
  commandIds?: Set<string>;
  addCommand: (commandId: string) => void;
  removeCommand: (commandId: string) => void;
};

export const CommandContext = React.createContext<CommandContextType>({
  fetchedCommands: {},
  commandIds: new Set<string>(),
  addCommand: () => Promise.resolve(),
  removeCommand: () => Promise.resolve(),
});
