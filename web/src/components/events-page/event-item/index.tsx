import { Box, ListItem, Typography, Chip } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import dayjs from "dayjs";
import { FC, memo } from "react";
import { EVENT_STATE_TEXT } from "~/constants/event-status-text";
import { useAppSelector } from "~/hooks/redux";
import { Event, selectById as selectEventById } from "~/modules/events";
import { EventStatusIcon } from "~/components/event-status-icon";

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    backgroundColor: theme.palette.background.paper,
  },
  info: {
    marginLeft: theme.spacing(1),
  },
  statusText: {
    marginLeft: theme.spacing(1),
    lineHeight: "1.5rem",
    // Fixed width to prevent misalignment of application name.
    width: "100px",
  },
  description: {
    color: theme.palette.text.secondary, // TODO check this color from hint #aaa to secondary #666
  },
  labelChip: {
    marginLeft: theme.spacing(1),
  },
}));

export interface EventItemProps {
  id: string;
}

const NO_DESCRIPTION = "No description.";

export const EventItem: FC<EventItemProps> = memo(function EventItem({ id }) {
  const classes = useStyles();
  const event = useAppSelector<Event.AsObject | undefined>((state) =>
    selectEventById(state.events, id)
  );

  if (!event) {
    return null;
  }

  return (
    <ListItem className={classes.root} dense divider>
      <Box display="flex" alignItems="center">
        <EventStatusIcon status={event.status} />
        <Typography
          variant="subtitle2"
          className={classes.statusText}
          component="span"
        >
          {EVENT_STATE_TEXT[event.status]}
        </Typography>
      </Box>
      <Box display="flex" flexDirection="column" flex={1} pl={2}>
        <Box display="flex" alignItems="baseline">
          <Typography variant="h6" component="span">
            {event.name}
          </Typography>
          <Typography
            variant="body2"
            color="textSecondary"
            className={classes.info}
          >
            {event.id}
            {event.labelsMap.map(([key, value], i) => (
              <Chip
                label={key + ": " + value}
                className={classes.labelChip}
                key={i}
              />
            ))}
          </Typography>
        </Box>
        <Typography variant="body1" className={classes.description}>
          {event.statusDescription || NO_DESCRIPTION}
        </Typography>
      </Box>
      <div>{dayjs(event.createdAt * 1000).fromNow()}</div>
    </ListItem>
  );
});
