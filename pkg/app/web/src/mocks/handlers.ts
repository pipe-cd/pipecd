import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
import { applicationHandlers } from "./services/application";
import { deploymentHandlers } from "./services/deployment";
import { projectHandlers } from "./services/project";
import { pipedHandlers } from "./services/piped";
import { environmentHandlers } from "./services/environment";
import { liveStateHandlers } from "./services/live-state";
import { stageLogHandlers } from "./services/stage-log";
import { apiKeyHandlers } from "./services/api-keys";
import { insightHandlers } from "./services/insight";

export const handlers = [
  ...meHandlers,
  ...commandHandlers,
  ...applicationHandlers,
  ...deploymentHandlers,
  ...projectHandlers,
  ...pipedHandlers,
  ...environmentHandlers,
  ...liveStateHandlers,
  ...stageLogHandlers,
  ...apiKeyHandlers,
  ...insightHandlers,
];
