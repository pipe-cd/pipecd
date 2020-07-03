import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export const APPLICATION_SYNC_STATUS_TEXT: Record<
  ApplicationSyncStatus,
  string
> = {
  [ApplicationSyncStatus.UNKNOWN]: "Unknown",
  [ApplicationSyncStatus.SYNCED]: "Synced",
  [ApplicationSyncStatus.DEPLOYING]: "Deploying",
  [ApplicationSyncStatus.OUT_OF_SYNC]: "Out of Sync",
};
