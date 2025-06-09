import {
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from "@mui/material";
import Autocomplete from "@mui/material/Autocomplete";
import { FC, memo, useCallback, useState, useEffect } from "react";
import { FilterView } from "~/components/filter-view";
import { EVENT_STATE_TEXT } from "~/constants/event-status-text";
import { Event, EventStatus } from "pipecd/web/model/event_pb";
import {
  EventFilterOptions,
  EventStatusKey,
} from "~/queries/events/use-get-events-infinite";

const ALL_VALUE = "ALL";

export interface EventFilterProps {
  events: Event.AsObject[];
  options: EventFilterOptions;
  onClear: () => void;
  onChange: (options: EventFilterOptions) => void;
}

export const EventFilter: FC<EventFilterProps> = memo(function EventFilter({
  events,
  options,
  onChange,
  onClear,
}) {
  const handleUpdateFilterValue = useCallback(
    (opts: Partial<EventFilterOptions>): void => {
      onChange({ ...options, ...opts });
    },
    [options, onChange]
  );

  const [allNames, setAllNames] = useState(new Array<string>());

  useEffect(() => {
    const names = new Set<string>();
    events.map((event) => {
      names.add(event.name);
    });
    setAllNames(Array.from(names));
  }, [events]);

  const [allLabels, setAllLabels] = useState(new Array<string>());
  const [selectedLabels, setSelectedLabels] = useState(new Array<string>());

  useEffect(() => {
    const labels = new Set<string>();
    events
      .filter((app) => app.labelsMap.length > 0)
      .forEach((app) => {
        app.labelsMap.forEach((label) => {
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
      <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
        <Autocomplete
          autoHighlight
          id="filter-event-name"
          noOptionsText="No selectable name"
          options={allNames}
          value={options.name ?? ""}
          onInputChange={(_, value) => {
            setAllNames([value]);
          }}
          onChange={(_, newValue) => {
            setAllNames([]);
            handleUpdateFilterValue({
              name: newValue !== null ? newValue : "",
            });
          }}
          renderInput={(params) => (
            <TextField
              {...params}
              variant="outlined"
              label="Name"
              margin="dense"
              fullWidth
            />
          )}
        />
      </FormControl>

      <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
        <InputLabel id="filter-event-status">Event Status</InputLabel>
        <Select
          labelId="filter-event-status"
          id="filter-event-status"
          value={options.status ?? ALL_VALUE}
          label="Event Status"
          fullWidth
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

      <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
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
