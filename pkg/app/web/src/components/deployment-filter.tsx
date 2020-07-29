import {
  Button,
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Paper,
  Select,
  TextField,
  Typography,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import React, { FC, memo, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "../constants/deployment-status-text";
import { AppState } from "../modules";
import {
  Application,
  ApplicationKind,
  ApplicationKindKey,
  selectAll as selectAllApplications,
  selectById as selectApplicationById,
} from "../modules/applications";
import {
  clearDeploymentFilter,
  DeploymentFilterOptions,
  updateDeploymentFilter,
} from "../modules/deployment-filter-options";
import { DeploymentStatus, DeploymentStatusKey } from "../modules/deployments";
import { Environment, selectAll } from "../modules/environments";

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

interface Props {
  open: boolean;
  onChange: () => void;
}

export const DeploymentFilter: FC<Props> = memo(function DeploymentFilter({
  open,
  onChange,
}) {
  const classes = useStyles();

  const dispatch = useDispatch();
  const envs = useSelector<AppState, Environment[]>((state) =>
    selectAll(state.environments)
  );
  const applications = useSelector<AppState, Application[]>((state) =>
    selectAllApplications(state.applications)
  );
  const options = useSelector<AppState, DeploymentFilterOptions>(
    (state) => state.deploymentFilterOptions
  );
  const selectedApp = useSelector<AppState, Application | undefined>((state) =>
    selectApplicationById(state.applications, options.applicationIds[0])
  );
  const handleUpdateFilterValue = useCallback(
    (opts: Partial<DeploymentFilterOptions>): void => {
      dispatch(updateDeploymentFilter(opts));
      onChange();
    },
    [dispatch, onChange]
  );

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
            dispatch(clearDeploymentFilter());
            onChange();
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
          value={options.envIds[0] || ALL_VALUE}
          label="Environment"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              envIds:
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
          value={options.kinds[0] ?? ALL_VALUE}
          label="Application Kind"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              kinds:
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

      <div className={classes.formItem}>
        <Autocomplete
          id="application-select"
          options={applications}
          getOptionLabel={(option) => option.name}
          renderOption={(option) => <span>{option.name}</span>}
          value={selectedApp || null}
          onChange={(_, value) => {
            handleUpdateFilterValue({
              applicationIds: value ? [value.id] : [],
            });
          }}
          renderInput={(params) => (
            <TextField
              {...params}
              label="Application"
              variant="outlined"
              inputProps={{
                ...params.inputProps,
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
          value={options.statuses[0] ?? ALL_VALUE}
          label="Deployment Status"
          className={classes.select}
          onChange={(e) => {
            handleUpdateFilterValue({
              statuses:
                e.target.value === ALL_VALUE
                  ? []
                  : [e.target.value as DeploymentStatus],
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
