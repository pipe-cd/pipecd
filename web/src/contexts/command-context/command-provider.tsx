import { FC, PropsWithChildren, useCallback, useState } from "react";
import { CommandContext } from "./command-context";
import CommandStatusTracking from "./CommandStatusTracking";
import { Command } from "~~/model/command_pb";

export const CommandProvider: FC<PropsWithChildren<unknown>> = ({
  children,
}) => {
  const [commandIds, setCommandIds] = useState(new Set<string>());

  const [fetchedCommands, setFetchedCommands] = useState<
    Record<string, Command.AsObject>
  >({});

  const addCommand = useCallback((commandId: string): void => {
    setCommandIds((prev) => new Set(prev).add(commandId));
  }, []);

  const removeCommand = useCallback((commandId: string): void => {
    setCommandIds((prev) => {
      const newIds = new Set(prev);
      newIds.delete(commandId);
      return newIds;
    });
    setFetchedCommands((prev) => {
      const newFetchedCommands = { ...prev };
      delete newFetchedCommands[commandId];
      return newFetchedCommands;
    });
  }, []);

  const saveFetchedCommand = useCallback((command: Command.AsObject): void => {
    setFetchedCommands((prev) => ({
      ...prev,
      [command.id]: command,
    }));
  }, []);

  return (
    <CommandContext.Provider
      value={{ fetchedCommands, addCommand, removeCommand, commandIds }}
    >
      {children}
      {[...commandIds].map((commandId) => (
        <CommandStatusTracking
          commandId={commandId}
          key={commandId}
          onComplete={() => {
            removeCommand(commandId);
          }}
          onFetched={(command) => {
            saveFetchedCommand(command);
          }}
        />
      ))}
    </CommandContext.Provider>
  );
};
