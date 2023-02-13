import { makeStyles } from "@material-ui/core";
import {
  CheckCircle,
  Error,
  IndeterminateCheckBox,
  Block,
} from "@material-ui/icons";
import { EventStatus } from "~/modules/events";
import { FC } from "react";
import clsx from "clsx";

const useStyles = makeStyles((theme) => ({
  [EventStatus.EVENT_NOT_HANDLED]: {
    color: theme.palette.grey[500],
  },
  [EventStatus.EVENT_SUCCESS]: {
    color: theme.palette.success.main,
  },
  [EventStatus.EVENT_FAILURE]: {
    color: theme.palette.error.main,
  },
  [EventStatus.EVENT_OUTDATED]: {
    color: theme.palette.warning.main,
  },
}));

export interface EventStatusIconProps {
  status: EventStatus;
  className?: string;
}

export const EventStatusIcon: FC<EventStatusIconProps> = ({
  status,
  className,
}) => {
  const classes = useStyles();

  switch (status) {
    case EventStatus.EVENT_NOT_HANDLED:
      return (
        <IndeterminateCheckBox
          className={clsx(classes[status], className)}
          data-testid="event-not-handled-icon"
        />
      );
    case EventStatus.EVENT_SUCCESS:
      return (
        <CheckCircle
          className={clsx(classes[status], className)}
          data-testid="event-success-icon"
        />
      );
    case EventStatus.EVENT_FAILURE:
      return (
        <Error
          className={clsx(classes[status], className)}
          data-testid="event-failure-icon"
        />
      );
    case EventStatus.EVENT_OUTDATED:
      return (
        <Block
          className={clsx(classes[status], className)}
          data-testid="event-outdated-icon"
        />
      );
  }
};
