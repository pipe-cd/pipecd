import { ApplicationActiveStatus } from "~/modules/applications";

export const APPLICATION_ACTIVE_STATUS_NAME: Record<
  ApplicationActiveStatus,
  string
> = {
  [ApplicationActiveStatus.ENABLED]: "ENABLED",
  [ApplicationActiveStatus.DISABLED]: "DISABLED",
  [ApplicationActiveStatus.DELETED]: "DELETED",
};
