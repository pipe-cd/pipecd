import {
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Select,
  TextField,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FC, memo, useState, useEffect } from "react";
import { FilterView } from "~/components/filter-view";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { APPLICATION_SYNC_STATUS_TEXT } from "~/constants/application-sync-status-text";
import { useAppSelector } from "~/hooks/redux";
import {
  Application,
  ApplicationKind,
  ApplicationKindKey,
  ApplicationsFilterOptions,
  ApplicationSyncStatus,
  ApplicationSyncStatusKey,
  selectAll as selectAllApplications,
} from "~/modules/applications";
import { ApplicationAutocomplete } from "./application-autocomplete";
import { PipedSelect } from "./piped-select";

const useStyles = makeStyles((theme) => ({
  toolbarSpacer: {
    flexGrow: 1,
  },
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
  select: {
    width: "100%",
  },
}));

export interface ApplicationFilterProps {
  options: ApplicationsFilterOptions;
  onChange: (options: ApplicationsFilterOptions) => void;
  onClear: () => void;
}

const ALL_VALUE = "ALL";

export const ApplicationFilter: FC<ApplicationFilterProps> = memo(
  function ApplicationFilter({ options, onChange, onClear }) {
    const classes = useStyles();
    const applications = useAppSelector<Application.AsObject[]>((state) =>
      selectAllApplications(state.applications)
    );

    const handleUpdateFilterValue = (
      optionPart: Partial<ApplicationsFilterOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const [allLabels, setAllLabels] = useState(new Array<string>());
    const [selectedLabels, setSelectedLabels] = useState(new Array<string>());

    useEffect(() => {
      const labels = new Set<string>();
      applications
        .filter((app) => app.labelsMap.length > 0)
        .map((app) => {
          app.labelsMap.map((label) => {
            labels.add(`${label[0]}:${label[1]}`);
          });
        });
      setAllLabels(Array.from(labels));
    }, [applications]);


    return (
      <FilterView
        onClear={() => {
          onClear();
          setSelectedLabels([]);
        }}
      >
        <div className={classes.formItem}>
          <ApplicationAutocomplete
            value={options.name ?? null}
            onChange={(value) => handleUpdateFilterValue({ name: value })}
          />
        </div>

        <div className={classes.formItem}>
          <PipedSelect
            value={options.pipedId ?? null}
            onChange={(value) => handleUpdateFilterValue({ pipedId: value })}
            className={classes.select}
          />
        </div>

        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-kind">Kind</InputLabel>
          <Select
            labelId="filter-kind"
            id="filter-kind"
            value={options.kind ?? ALL_VALUE}
            label="Kind"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                kind:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : (e.target.value as string),
              });
            }}
          >
            <MenuItem value={ALL_VALUE}>
              <em>All</em>
            </MenuItem>

            {Object.keys(ApplicationKind).map((key) => (
              <MenuItem
                value={ApplicationKind[key as ApplicationKindKey]}
                key={`status-${key}`}
              >
                {
                  APPLICATION_KIND_TEXT[
                    ApplicationKind[key as ApplicationKindKey]
                  ]
                }
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-sync-status">Sync Status</InputLabel>
          <Select
            labelId="filter-sync-status"
            id="filter-sync-status"
            value={options.syncStatus ?? ALL_VALUE}
            label="Sync Status"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                syncStatus:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : (e.target.value as string),
              });
            }}
          >
            <MenuItem value={ALL_VALUE}>
              <em>All</em>
            </MenuItem>

            {Object.keys(ApplicationSyncStatus).map((key) => (
              <MenuItem
                value={ApplicationSyncStatus[key as ApplicationSyncStatusKey]}
                key={`sync-status-${key}`}
              >
                {
                  APPLICATION_SYNC_STATUS_TEXT[
                    ApplicationSyncStatus[key as ApplicationSyncStatusKey]
                  ]
                }
              </MenuItem>
            ))}
          </Select>
        </FormControl>

        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-active-status">Active Status</InputLabel>
          <Select
            labelId="filter-active-status"
            id="filter-active-status"
            value={
              options.activeStatus === undefined
                ? ALL_VALUE
                : options.activeStatus
            }
            label="Active Status"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                activeStatus:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : (e.target.value as string),
              });
            }}
          >
            <MenuItem value={ALL_VALUE}>
              <em>All</em>
            </MenuItem>
            <MenuItem value="enabled">Enabled</MenuItem>
            <MenuItem value="disabled">Disabled</MenuItem>
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
  }
);
