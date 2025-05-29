import {
  CheckCircle,
  Error,
  IndeterminateCheckBox,
  Block,
} from "@mui/icons-material";
import { EventStatus } from "~/modules/events";
import { FC } from "react";

export interface EventStatusIconProps {
  status: EventStatus;
}

export const EventStatusIcon: FC<EventStatusIconProps> = ({ status }) => {
  switch (status) {
    case EventStatus.EVENT_NOT_HANDLED:
      return (
        <IndeterminateCheckBox
          data-testid="event-not-handled-icon"
          sx={{
            color: "grey.500",
          }}
        />
      );
    case EventStatus.EVENT_SUCCESS:
      return (
        <CheckCircle
          data-testid="event-success-icon"
          sx={{
            color: "success.main",
          }}
        />
      );
    case EventStatus.EVENT_FAILURE:
      return (
        <Error
          data-testid="event-failure-icon"
          sx={{
            color: "error.main",
          }}
        />
      );
    case EventStatus.EVENT_OUTDATED:
      return (
        <Block
          data-testid="event-outdated-icon"
          sx={{
            color: "warning.main",
          }}
        />
      );
  }
};
