import React, { FC, useReducer, memo, useEffect } from "react";
import {
  makeStyles,
  Paper,
  Button,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import {
  ApplicationKind,
  ApplicationKindKey,
  Application,
  selectAll as selectAllApplications,
} from "../modules/applications";
import { DeploymentStatus, DeploymentStatusKey } from "../modules/deployments";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { Environment, selectAll } from "../modules/environments";
import { AppState } from "../modules";
import { useSelector } from "react-redux";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";

const FILTER_PAPER_WIDTH = 360;

const useStyles = makeStyles((theme) => ({
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
    height: "100%",
  },
  select: {
    width: "100%",
  },
}));

const ALL_VALUE = "ALL";

interface FormState {
  deploymentStatus: DeploymentStatus | typeof ALL_VALUE;
  applicationKind: ApplicationKind | typeof ALL_VALUE;
  application: Application | null;
  env: string;
}

const initialState: FormState = {
  deploymentStatus: ALL_VALUE,
  applicationKind: ALL_VALUE,
  application: null,
  env: ALL_VALUE,
};

type Actions =
  | {
      type: "update-deployment-status";
      value: DeploymentStatus | typeof ALL_VALUE;
    }
  | {
      type: "update-application-kind";
      value: ApplicationKind | typeof ALL_VALUE;
    }
  | { type: "update-application"; value: Application | null }
  | { type: "update-env"; value: string }
  | {
      type: "clear-form";
    };

const reducer = (state: FormState, action: Actions): FormState => {
  switch (action.type) {
    case "clear-form":
      return initialState;
    case "update-deployment-status":
      return { ...state, deploymentStatus: action.value };
    case "update-application-kind":
      return { ...state, applicationKind: action.value };
    case "update-application":
      return { ...state, application: action.value };
    case "update-env":
      return { ...state, env: action.value };
  }
};

interface Options {
  statusesList: DeploymentStatus[];
  kindsList: ApplicationKind[];
  applicationIdsList: string[];
  envIdsList: string[];
}

interface Props {
  open: boolean;
  onChange: (options: Options) => void;
}

export const DeploymentFilter: FC<Props> = memo(function DeploymentFilter({
  open,
  onChange,
}) {
  const classes = useStyles();

  const envs = useSelector<AppState, Environment[]>((state) =>
    selectAll(state.environments)
  );
  const applications = useSelector<AppState, Application[]>((state) =>
    selectAllApplications(state.applications)
  );

  const [state, dispatch] = useReducer(reducer, initialState);

  useEffect(() => {
    const options: Options = {
      statusesList: [],
      kindsList: [],
      envIdsList: [],
      applicationIdsList: [],
    };
    if (state.deploymentStatus !== ALL_VALUE) {
      options.statusesList = [state.deploymentStatus];
    }
    if (state.applicationKind !== ALL_VALUE) {
      options.kindsList = [state.applicationKind];
    }
    if (state.env !== ALL_VALUE) {
      options.envIdsList = [state.env];
    }
    if (state.application) {
      options.applicationIdsList = [state.application.id];
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
          label="Application Kind"
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

      <div className={classes.formItem}>
        <Autocomplete
          id="application-select"
          options={applications}
          getOptionLabel={(option) => option.name}
          renderOption={(option) => <span>{option.name}</span>}
          value={state.application}
          onChange={(_, value) => {
            dispatch({
              type: "update-application",
              value: value ? value : null,
            });
          }}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Application"
              variant="outlined"
              inputProps={{
                ...params.inputProps,
                autoComplete: "new-password", // disable autocomplete and autofill
              }}
            />
          )}
        />
      </div>

      <FormControl className={classes.formItem} variant="outlined">
        <InputLabel id="filter-deployment-status">Deployment Status</InputLabel>
        <Select
          labelId="filter-deployment-status"
          id="filter-deployment-status"
          value={state.deploymentStatus}
          label="Deployment Status"
          className={classes.select}
          onChange={(e) => {
            dispatch({
              type: "update-deployment-status",
              value:
                e.target.value === ALL_VALUE
                  ? ALL_VALUE
                  : (e.target.value as DeploymentStatus),
            });
          }}
        >
          <MenuItem value={ALL_VALUE}>
            <em>All</em>
          </MenuItem>

          {Object.keys(DeploymentStatus).map((key) => (
            <MenuItem
              key={`deployment-status-${key}`}
              value={DeploymentStatus[key as DeploymentStatusKey]}
            >
              {
                DEPLOYMENT_STATE_TEXT[
                  DeploymentStatus[key as DeploymentStatusKey]
                ]
              }
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </Paper>
  );
});
