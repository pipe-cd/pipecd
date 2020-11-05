import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
import { applicationHandlers } from "./services/application";
export const handlers = [
  ...meHandlers,
  ...commandHandlers,
  ...applicationHandlers,
];
