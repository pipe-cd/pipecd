import { meHandlers } from "./services/me";
import { commandHandlers } from "./services/command";
export const handlers = [...meHandlers, ...commandHandlers];
