import { ApplicationLiveStateSnapshot } from "~/modules/applications-live-state";

export const APPLICATION_HEALTH_STATUS_TEXT: Record<
  ApplicationLiveStateSnapshot.Status,
  string
> = {
  [ApplicationLiveStateSnapshot.Status.UNKNOWN]: "Unknown",
  [ApplicationLiveStateSnapshot.Status.HEALTHY]: "Healthy",
  [ApplicationLiveStateSnapshot.Status.OTHER]: "Other",
  [ApplicationLiveStateSnapshot.Status.UNHEALTHY]: "Unhealthy",
};
