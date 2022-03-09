import { EventStatus } from "pipecd/web/model/event_pb";

export const EVENT_STATE_TEXT: Record<EventStatus, string> = {
  [EventStatus.EVENT_NOT_HANDLED]: "NOT HANDLED",
  [EventStatus.EVENT_SUCCESS]: "SUCCESS",
  [EventStatus.EVENT_FAILURE]: "FAILURE",
  [EventStatus.EVENT_OUTDATED]: "OUTDATED",
};
