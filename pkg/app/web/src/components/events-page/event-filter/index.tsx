import {
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Select,
  TextField,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FC, memo, useCallback, useState, useEffect } from "react";
import { FilterView } from "~/components/filter-view";
import { EVENT_STATE_TEXT } from "~/constants/event-status-text";
import { useAppSelector } from "~/hooks/redux";
import {
  Event,
  EventFilterOptions,
  EventStatus,
  EventStatusKey,
  selectAll as selectAllEvents,
} from "~/modules/events";

const useStyles = makeStyles((theme) => ({
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
  select: {
    width: "100%",
  },
}));

const ALL_VALUE = "ALL";

export interface EventFilterProps {
  options: EventFilterOptions;
  onClear: () => void;
  onChange: (options: EventFilterOptions) => void;
}

export const EventFilter: FC<EventFilterProps> = memo(function EventFilter({
  options,
  onChange,
  onClear,
}) {
  const classes = useStyles();
  const handleUpdateFilterValue = useCallback(
    (opts: Partial<EventFilterOptions>): void => {
      onChange({ ...options, ...opts });
    },
    [options, onChange]
  );

  const events = useAppSelector<Event.AsObject[]>((state) =>
    selectAllEvents(state.events)
  );
  const [allLabels, setAllLabels] = useState(new Array<string>());
  const [selectedLabels, setSelectedLabels] = useState(new Array<string>());

  useEffect(() => {
    const labels = new Set<string>();
    events
      .filter((app) => app.labelsMap.length > 0)
      .map((app) => {
        app.labelsMap.map((label) => {
          labels.add(`${label[0]}:${label[1]}`);
        });
      });
    setAllLabels(Array.from(labels));
  }, [events]);

  return (
    <FilterView
      onClear={() => {
        onClear();
      }}
    >
      {/* TODO: Suggest possible event names on the event filter */}
      <FormControl className={classes.formItem} variant="outlined">
        <TextField
          id="filter-event-name"
          value={options.name ?? ""}
          label="Name"
          className={classes.select}
          variant="outlined"
          onChange={(e) => {
            handleUpdateFilterValue({
              name:
                e.target.value === ALL_VALUE ? undefined : `${e.target.value}`,
            });
          }}
        ></TextField>
      </FormControl>

      <FormControl className={classes.formItem} variant="outlined">
        <InputLabel id="filter-event-status">Event Status</InputLabel>
        <Select
          labelId="filter-event-status"
          id="filter-event-status"
          value={options.status ?? ALL_VALUE}
          label="Event Status"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              status:
                e.target.value === ALL_VALUE ? undefined : `${e.target.value}`,
            });
          }}
        >
          <MenuItem value={ALL_VALUE}>
            <em>All</em>
          </MenuItem>

          {Object.keys(EventStatus).map((key) => (
            <MenuItem
              key={`event-status-${key}`}
              value={EventStatus[key as EventStatusKey]}
            >
              {EVENT_STATE_TEXT[EventStatus[key as EventStatusKey]]}
            </MenuItem>
          ))}
        </Select>
      </FormControl>

      <FormControl className={classes.formItem} variant="outlined">
        <Autocomplete
          multiple
          autoHighlight
          id="labels"
          noOptionsText="No selectable labels"
          options={allLabels}
          value={options.labels ?? selectedLabels}
          onInputChange={(_, value) => {
            const label = value.split(":");
            if (label.length !== 2) return;
            if (label[0].length === 0) return;
            if (label[1].length === 0) return;
            setAllLabels([value]);
          }}
          onChange={(_, newValue) => {
            setAllLabels([]);
            setSelectedLabels(newValue);
            handleUpdateFilterValue({
              labels: newValue,
            });
          }}
          renderInput={(params) => (
            <TextField
              {...params}
              variant="outlined"
              label="Labels"
              margin="dense"
              placeholder="key:value"
              fullWidth
            />
          )}
        />
      </FormControl>
    </FilterView>
  );
});
