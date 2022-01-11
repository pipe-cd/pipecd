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
  ApplicationKind,
  ApplicationKindKey,
  ApplicationsFilterOptions,
  ApplicationSyncStatus,
  ApplicationSyncStatusKey,
} from "~/modules/applications";
import { selectAllEnvs } from "~/modules/environments";
import { ApplicationAutocomplete } from "./application-autocomplete";

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
    const envs = useAppSelector(selectAllEnvs);

    const handleUpdateFilterValue = (
      optionPart: Partial<ApplicationsFilterOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    const [labelOptions, setLabelOptions] = useState(new Array<string>());
    const [selectedLabels, setSelectedLabels] = useState(new Array<string>());

    const allLabels = useAppSelector((state) => state.applications.allLabels);
    useEffect(() => {
      const opts = new Array<string>();
      for (const label in allLabels) {
        opts.push(label);
      }
      setLabelOptions(opts);
    }, [allLabels]);

    return (
      <FilterView
        onClear={() => {
          onClear();
        }}
      >
        <div className={classes.formItem}>
          <ApplicationAutocomplete
            value={options.name ?? null}
            onChange={(value) => handleUpdateFilterValue({ name: value })}
          />
        </div>

        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-env">Environment</InputLabel>
          <Select
            labelId="filter-env"
            id="filter-env"
            value={options.envId ?? ALL_VALUE}
            label="Environment"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                envId:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : (e.target.value as string),
              });
            }}
          >
            <MenuItem value={ALL_VALUE}>
              <em>All</em>
            </MenuItem>
            {envs.map((e) => (
              <MenuItem value={e.id} key={`env-${e.id}`}>
                {e.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>

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
            options={labelOptions}
            value={options.labels ?? selectedLabels}
            onInputChange={(_, value) => {
              const label = value.split(":");
              if (label.length !== 2) return;
              if (label[0].length === 0) return;
              if (label[1].length === 0) return;
              setLabelOptions([value]);
            }}
            onChange={(_, newValue) => {
              setLabelOptions([]);
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
