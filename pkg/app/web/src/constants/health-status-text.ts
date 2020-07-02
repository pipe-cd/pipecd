import { ApplicationLiveStateSnapshotModel } from "../modules/applications-live-state";

export const APPLICATION_HEALTH_STATUS_TEXT: Record<
  ApplicationLiveStateSnapshotModel.Status,
  string
> = {
  [ApplicationLiveStateSnapshotModel.Status.UNKNOWN]: "Unknown",
  [ApplicationLiveStateSnapshotModel.Status.HEALTHY]: "Healthy",
  [ApplicationLiveStateSnapshotModel.Status.OTHER]: "Other",
};
