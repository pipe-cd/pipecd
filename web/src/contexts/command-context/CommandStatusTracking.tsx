import { FC, memo, useEffect } from "react";
import {
  COMMAND_TYPE_TEXT,
  useFetchCommand,
} from "~/queries/commands/use-fetch-command";
import { useToast } from "../toast-context";
import { Command, CommandStatus } from "~~/model/command_pb";
import { findMetadataByKey } from "~/utils/find-metadata-by-key";
import { PAGE_PATH_DEPLOYMENTS } from "~/constants/path";

type Props = {
  commandId: string;
  onComplete: (commandId: string) => void;
  onFetched: (command: Command.AsObject) => void;
};

const METADATA_KEY = {
  TRIGGERED_DEPLOYMENT_ID: "TriggeredDeploymentID",
};

const CommandStatusTracking: FC<Props> = ({
  commandId,
  onComplete,
  onFetched,
}: Props) => {
  const { addToast } = useToast();
  const { data: command, isSuccess } = useFetchCommand({ commandId });

  // show toast when the command is successfully fetched and handled
  useEffect(() => {
    if (!isSuccess) return;
    if (command === undefined) return;
    if (command.status !== CommandStatus.COMMAND_SUCCEEDED) return;

    switch (command.type) {
      case Command.Type.SYNC_APPLICATION: {
        const deploymentId = findMetadataByKey(
          command.metadataMap,
          METADATA_KEY.TRIGGERED_DEPLOYMENT_ID
        );
        addToast({
          message: `Succeed "${COMMAND_TYPE_TEXT[command.type]}"`,
          severity: "success",
          to: deploymentId
            ? `${PAGE_PATH_DEPLOYMENTS}/${deploymentId}`
            : undefined,
        });
        break;
      }
      default:
        addToast({
          message: `Succeed "${COMMAND_TYPE_TEXT[command.type]}"`,
          severity: "success",
        });
    }
  }, [addToast, command, isSuccess]);

  // Stop tracking the command when it is handled.
  useEffect(() => {
    if (
      isSuccess &&
      command?.status !== CommandStatus.COMMAND_NOT_HANDLED_YET
    ) {
      onComplete(commandId);
    }
  }, [command?.status, commandId, isSuccess, onComplete]);

  useEffect(() => {
    if (command) {
      onFetched(command);
    }
  }, [command, onFetched]);

  return null;
};

export default memo(CommandStatusTracking, (prev, next) => {
  return prev.commandId === next.commandId;
});
