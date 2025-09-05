import {
  UseMutateFunction,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useMemo, useState } from "react";
import * as applicationsAPI from "~/api/applications";
import { useCommand } from "~/contexts/command-context";
import { SyncApplicationRequest } from "~~/api_client/service_pb";

const useIsSyncingApplication = ({
  commandId,
  isInitLoading,
}: {
  applicationId?: string;
  commandId?: string;
  isInitLoading: boolean;
}): boolean => {
  const { commandIds } = useCommand();

  const isSyncing = useMemo(() => {
    // fetchSync => store command => tracking command => remove command when finished
    const loadingStatus = {
      // from application syncing is start till result (commandId) is Store for tracking status
      isInitLoading,
      // when CommandId is stored for tracking status, after finished it will be removed
      isCommandRunning: commandIds?.has(commandId ?? ""),
    };

    return Object.values(loadingStatus).some((v) => Boolean(v));
  }, [commandId, commandIds, isInitLoading]);

  return isSyncing;
};

export const useSyncApplication = (): {
  mutate: UseMutateFunction<string, unknown, SyncApplicationRequest.AsObject>;
  isSyncing: boolean;
} => {
  const [commandId, setCommandId] = useState<string>();

  const queryClient = useQueryClient();
  const { addCommand } = useCommand();

  const { mutate, isLoading: isInitLoading } = useMutation({
    mutationFn: async (payload: SyncApplicationRequest.AsObject) => {
      const { commandId } = await applicationsAPI.syncApplication(payload);
      return commandId;
    },
    onSuccess: (commandId) => {
      addCommand(commandId);
      setCommandId(commandId);
      queryClient.invalidateQueries(["applications"]);
    },
  });

  const isSyncing = useIsSyncingApplication({
    isInitLoading,
    commandId,
  });

  return { mutate, isSyncing };
};
