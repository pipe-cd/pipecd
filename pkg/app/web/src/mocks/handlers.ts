import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
import { applicationHandlers } from "./services/application";
import { deploymentHandlers } from "./services/deployment";
import { projectHandlers } from "./services/project";
import { pipedHandlers } from "./services/piped";
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
  ...liveStateHandlers,
  ...stageLogHandlers,
  ...apiKeyHandlers,
  ...insightHandlers,
];
