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

type Ability = typeof ALL_VALUE | "enabled" | "disabled";

interface FormState {
  status: ApplicationSyncStatus | typeof ALL_VALUE;
  kind: ApplicationKind | typeof ALL_VALUE;
  ability: Ability;
  env: string;
}

const initialState: FormState = {
  ability: ALL_VALUE,
  status: ALL_VALUE,
  env: ALL_VALUE,
  kind: ALL_VALUE,
};
type Actions =
  | { type: "update-status"; value: ApplicationSyncStatus | typeof ALL_VALUE }
  | { type: "update-kind"; value: ApplicationKind | typeof ALL_VALUE }
  | { type: "update-env"; value: string }
  | { type: "update-ability"; value: Ability }
  | { type: "clear-form" };
const reducer = (state: FormState, action: Actions): FormState => {
  switch (action.type) {
    case "update-ability":
      return { ...state, ability: action.value };
    case "update-status":
      return { ...state, status: action.value };
    case "update-kind":
      return { ...state, kind: action.value };
    case "update-env":
      return { ...state, env: action.value };
    case "clear-form":
      return initialState;
  }
};

interface Options {
  enabled?: { value: boolean };
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
    if (state.ability !== ALL_VALUE) {
      options.enabled = {
        value: state.ability === "enabled",
      };
    }
    if (state.kind !== ALL_VALUE) {
      options.kindsList = [state.kind];
    }
    if (state.env !== ALL_VALUE) {
      options.envIdsList = [state.env];
    }
    if (state.status !== ALL_VALUE) {
      options.syncStatusesList = [state.status];
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
        <InputLabel id="filter-ability">Ability</InputLabel>
        <Select
          labelId="filter-ability"
          id="filter-ability"
          value={state.ability}
          label="Ability"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-ability",
              value: e.target.value as Ability,
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
        <InputLabel id="filter-status">Status</InputLabel>
        <Select
          labelId="filter-status"
          id="filter-status"
          value={state.status}
          label="Status"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-status",
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
              key={`status-${key}`}
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
          <MenuItem value={ALL_VALUE} key="env-all">
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
          value={state.kind}
          label="Environment"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-kind",
              value:
                e.target.value === ""
                  ? ALL_VALUE
                  : (e.target.value as ApplicationKind),
            });
          }}
        >
          <MenuItem value={ALL_VALUE} key="kind-all">
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
