import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
import { applicationHandlers } from "./services/application";
import { deploymentHandlers } from "./services/deployment";
import { projectHandlers } from "./services/project";
import { pipedHandlers } from "./services/piped";
export const handlers = [
  ...meHandlers,
  ...commandHandlers,
  ...applicationHandlers,
  ...deploymentHandlers,
  ...projectHandlers,
  ...pipedHandlers,
];
