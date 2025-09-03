import { useQuery, UseQueryResult } from "@tanstack/react-query";
import { getCommand } from "~/api/commands";
import { Command, CommandStatus } from "~~/model/command_pb";

const FETCH_COMMANDS_INTERVAL = 3000;

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

export const useFetchCommand = ({
  commandId,
}: {
  commandId: string;
}): UseQueryResult<Command.AsObject> => {
  return useQuery({
    queryKey: ["command", commandId],
    queryFn: async () => {
      const { command } = await getCommand({ commandId });
      return command;
    },
    refetchInterval(data) {
      // If the command is not handled yet, we will refetch it every 5 seconds.
      if (data?.status === CommandStatus.COMMAND_NOT_HANDLED_YET) {
        return FETCH_COMMANDS_INTERVAL;
      }
      // Otherwise, we don't need to refetch.
      return false;
    },
  });
};
