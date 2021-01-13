import React from "react";
import { SyncStatusIcon } from "./sync-status-icon";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export default {
  title: "APPLICATION/SyncStatusIcon",
  component: SyncStatusIcon,
};

export const overview: React.FC = () => (
  <>
    <SyncStatusIcon status={ApplicationSyncStatus.UNKNOWN} deploying={false} />
    <SyncStatusIcon status={ApplicationSyncStatus.SYNCED} deploying={false} />
    <SyncStatusIcon status={ApplicationSyncStatus.DEPLOYING} deploying={true} />
    <SyncStatusIcon
      status={ApplicationSyncStatus.OUT_OF_SYNC}
      deploying={false}
    />
  </>
);
