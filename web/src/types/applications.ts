import { ApplicationSyncStatus } from "~~/model/application_pb";
import { ApplicationKind } from "~~/model/common_pb";
export type ApplicationSyncStatusKey = keyof typeof ApplicationSyncStatus;
export type ApplicationKindKey = keyof typeof ApplicationKind;

export {
  Application,
  ApplicationSyncState,
  ApplicationSyncStatus,
  ApplicationDeploymentReference,
} from "~~/model/application_pb";

export {
  ApplicationKind,
  ApplicationActiveStatus,
  ApplicationGitRepository,
  ApplicationInfo,
} from "pipecd/web/model/common_pb";
