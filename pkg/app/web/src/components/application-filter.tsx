import {
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Select,
} from "@material-ui/core";
import React, { FC, memo } from "react";
import { useDispatch, useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  ApplicationFilterOptions,
  clearApplicationFilter,
  updateApplicationFilter,
} from "../modules/application-filter-options";
import {
  ApplicationKind,
  ApplicationKindKey,
  ApplicationSyncStatus,
  ApplicationSyncStatusKey,
} from "../modules/applications";
import { Environment, selectAll } from "../modules/environments";
import { FilterView } from "./filter-view";

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

interface Props {
  onChange: () => void;
}

const ALL_VALUE = "ALL";
const getActiveStatusText = (v: boolean): string =>
  v ? "enabled" : "disabled";

export const ApplicationFilter: FC<Props> = memo(function ApplicationFilter({
  onChange,
}) {
  const classes = useStyles();
  const dispatch = useDispatch();
  const envs = useSelector<AppState, Environment.AsObject[]>((state) =>
    selectAll(state.environments)
  );
  const options = useSelector<AppState, ApplicationFilterOptions>(
    (state) => state.applicationFilterOptions
  );

  const handleUpdateFilterValue = (
    options: Partial<ApplicationFilterOptions>
  ): void => {
    dispatch(updateApplicationFilter(options));
    onChange();
  };

  return (
    <FilterView
      onClear={() => {
        dispatch(clearApplicationFilter());
        onChange();
      }}
    >
      <FormControl className={classes.formItem} variant="outlined">
        <InputLabel id="filter-env">Environment</InputLabel>
        <Select
          labelId="filter-env"
          id="filter-env"
          value={options.envIdsList[0] || ALL_VALUE}
          label="Environment"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              envIdsList:
                e.target.value === ALL_VALUE ? [] : [e.target.value as string],
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
        <InputLabel id="filter-application-kind">Application Kind</InputLabel>
        <Select
          labelId="filter-application-kind"
          id="filter-application-kind"
          value={options.kindsList[0] ?? ALL_VALUE}
          label="Application Kind"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              kindsList:
                e.target.value === ALL_VALUE
                  ? []
                  : [e.target.value as ApplicationKind],
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
          value={options.syncStatusesList[0] ?? ALL_VALUE}
          label="Sync Status"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              syncStatusesList:
                e.target.value === ALL_VALUE
                  ? []
                  : [e.target.value as ApplicationSyncStatus],
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
            options.enabled === undefined
              ? ALL_VALUE
              : getActiveStatusText(options.enabled.value)
          }
          label="Active Status"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              enabled:
                e.target.value === ALL_VALUE
                  ? undefined
                  : { value: e.target.value === "enabled" },
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
    </FilterView>
  );
});
