import {
  Button,
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Paper,
  Select,
  Typography,
} from "@material-ui/core";
import React, { FC, useReducer, memo, useEffect } from "react";
import { useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { APPLICATION_SYNC_STATUS_TEXT } from "../constants/application-sync-status-text";
import { AppState } from "../modules";
import {
  ApplicationKind,
  ApplicationKindKey,
  ApplicationSyncStatus,
  ApplicationSyncStatusKey,
} from "../modules/applications";
import { Environment, selectAll } from "../modules/environments";

const FILTER_PAPER_WIDTH = 360;

const useStyles = makeStyles((theme) => ({
  main: {
    display: "flex",
  },
  header: {
    display: "flex",
    justifyContent: "space-between",
  },
  toolbarSpacer: {
    flexGrow: 1,
  },
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
  filterPaper: {
    width: FILTER_PAPER_WIDTH,
    padding: theme.spacing(3),
  },
  select: {
    width: "100%",
  },
}));

const ALL_VALUE = "ALL";

type ActiveStatus = typeof ALL_VALUE | "enabled" | "disabled";

interface FormState {
  syncStatus: ApplicationSyncStatus | typeof ALL_VALUE;
  applicationKind: ApplicationKind | typeof ALL_VALUE;
  activeStatus: ActiveStatus;
  env: string;
}

const initialState: FormState = {
  activeStatus: "enabled",
  syncStatus: ALL_VALUE,
  env: ALL_VALUE,
  applicationKind: ALL_VALUE,
};

type Actions =
  | {
      type: "update-sync-status";
      value: ApplicationSyncStatus | typeof ALL_VALUE;
    }
  | {
      type: "update-application-kind";
      value: ApplicationKind | typeof ALL_VALUE;
    }
  | { type: "update-env"; value: string }
  | { type: "update-active-status"; value: ActiveStatus }
  | { type: "clear-form" };
const reducer = (state: FormState, action: Actions): FormState => {
  switch (action.type) {
    case "update-active-status":
      return { ...state, activeStatus: action.value };
    case "update-sync-status":
      return { ...state, syncStatus: action.value };
    case "update-application-kind":
      return { ...state, applicationKind: action.value };
    case "update-env":
      return { ...state, env: action.value };
    case "clear-form":
      return initialState;
  }
};

interface Options {
  enabled?: {
    value: boolean;
  };
  kindsList: ApplicationKind[];
  envIdsList: string[];
  syncStatusesList: ApplicationSyncStatus[];
}

interface Props {
  open: boolean;
  onChange: (props: Options) => void;
}

export const ApplicationFilter: FC<Props> = memo(function ApplicationFilter({
  open,
  onChange,
}) {
  const classes = useStyles();
  const envs = useSelector<AppState, Environment[]>((state) =>
    selectAll(state.environments)
  );

  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const options: Options = {
      kindsList: [],
      envIdsList: [],
      syncStatusesList: [],
    };
    if (state.activeStatus !== ALL_VALUE) {
      options.enabled = {
        value: state.activeStatus === "enabled",
      };
    }
    if (state.applicationKind !== ALL_VALUE) {
      options.kindsList = [state.applicationKind];
    }
    if (state.env !== ALL_VALUE) {
      options.envIdsList = [state.env];
    }
    if (state.syncStatus !== ALL_VALUE) {
      options.syncStatusesList = [state.syncStatus];
    }
    onChange(options);
  }, [state, onChange]);

  if (open === false) {
    return null;
  }

  return (
    <Paper className={classes.filterPaper} square>
      <div className={classes.header}>
        <Typography variant="h6">Filters</Typography>
        <Button
          color="primary"
          onClick={() => {
            dispatch({ type: "clear-form" });
          }}
        >
          Clear
        </Button>
      </div>

      <FormControl className={classes.formItem} variant="outlined">
        <InputLabel id="filter-active-status">Active Status</InputLabel>
        <Select
          labelId="filter-active-status"
          id="filter-active-status"
          value={state.activeStatus}
          label="Active Status"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-active-status",
              value: e.target.value as ActiveStatus,
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
        <InputLabel id="filter-sync-status">Sync Status</InputLabel>
        <Select
          labelId="filter-sync-status"
          id="filter-sync-status"
          value={state.syncStatus}
          label="Sync Status"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-sync-status",
              value:
                e.target.value === ""
                  ? ALL_VALUE
                  : (e.target.value as ApplicationSyncStatus),
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
        <InputLabel id="filter-env">Environment</InputLabel>
        <Select
          labelId="filter-env"
          id="filter-env"
          value={state.env}
          label="Environment"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-env",
              value: e.target.value as string,
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
          value={state.applicationKind}
          label="Environment"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-application-kind",
              value:
                e.target.value === ""
                  ? ALL_VALUE
                  : (e.target.value as ApplicationKind),
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
    </Paper>
  );
});
