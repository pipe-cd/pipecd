import { Box, ListItem, Typography, Chip } from "@mui/material";
import dayjs from "dayjs";
import { FC, memo } from "react";
import { EVENT_STATE_TEXT } from "~/constants/event-status-text";
import { useAppSelector } from "~/hooks/redux";
import { Event, selectById as selectEventById } from "~/modules/events";
import { EventStatusIcon } from "~/components/event-status-icon";

export interface EventItemProps {
  id: string;
}

const NO_DESCRIPTION = "No description.";

export const EventItem: FC<EventItemProps> = memo(function EventItem({ id }) {
  const event = useAppSelector<Event.AsObject | undefined>((state) =>
    selectEventById(state.events, id)
  );

  if (!event) {
    return null;
  }

  return (
    <ListItem
      sx={{
        flex: 1,
        padding: 2,
        display: "flex",
        alignItems: "center",
        backgroundColor: "background.paper",
      }}
      dense
      divider
    >
      <Box display="flex" alignItems="center">
        <EventStatusIcon status={event.status} />
        <Typography
          variant="subtitle2"
          sx={{
            marginLeft: 1,
            lineHeight: "1.5rem",
            // Fixed width to prevent misalignment of application name.
            width: "100px",
          }}
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
          <Typography variant="body2" color="textSecondary" ml={1}>
            {event.id}
            {event.labelsMap.map(([key, value], i) => (
              <Chip label={key + ": " + value} sx={{ ml: 1 }} key={i} />
            ))}
          </Typography>
        </Box>
        <Typography variant="body1" sx={{ color: "text.secondary" }}>
          {event.statusDescription || NO_DESCRIPTION}
        </Typography>
      </Box>
      <div>{dayjs(event.createdAt * 1000).fromNow()}</div>
    </ListItem>
  );
});
