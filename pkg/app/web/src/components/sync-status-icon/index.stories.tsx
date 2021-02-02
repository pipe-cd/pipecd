import React from "react";
import { SyncStatusIcon } from "./";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export default {
  title: "APPLICATION/SyncStatusIcon",
  component: SyncStatusIcon,
};

export const overview: React.FC = () => (
  <>
    <SyncStatusIcon status={ApplicationSyncStatus.UNKNOWN} />
    <SyncStatusIcon status={ApplicationSyncStatus.SYNCED} />
    <SyncStatusIcon status={ApplicationSyncStatus.DEPLOYING} />
    <SyncStatusIcon status={ApplicationSyncStatus.OUT_OF_SYNC} />
  </>
);
