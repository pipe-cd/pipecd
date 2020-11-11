import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
import { applicationHandlers } from "./services/application";
import { deploymentHandlers } from "./services/deployment";
export const handlers = [
  ...meHandlers,
  ...commandHandlers,
  ...applicationHandlers,
  ...deploymentHandlers,
];
