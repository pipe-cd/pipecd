import React from "react";
import { dummyApplicationSyncState } from "../../__fixtures__/dummy-application";

import { AppSyncStatus } from "./";

export default {
  title: "application/AppSyncStatus",
  component: AppSyncStatus,
};

export const overview: React.FC = () => (
  <AppSyncStatus deploying={false} syncState={dummyApplicationSyncState} />
);

export const large: React.FC = () => (
  <AppSyncStatus
    deploying={false}
    size="large"
    syncState={dummyApplicationSyncState}
  />
);

export const deploying: React.FC = () => (
  <AppSyncStatus deploying={true} syncState={dummyApplicationSyncState} />
);
