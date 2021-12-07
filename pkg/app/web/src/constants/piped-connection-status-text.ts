import { Piped } from "pipe/pkg/app/web/model/piped_pb";

export const PIPED_CONNECTION_STATUS_TEXT: Record<Piped.ConnectionStatus, string> = {
  [Piped.ConnectionStatus.UNKNOWN]: "Unknown",
  [Piped.ConnectionStatus.ONLINE]: "Online",
  [Piped.ConnectionStatus.OFFLINE]: "Offline",
};
