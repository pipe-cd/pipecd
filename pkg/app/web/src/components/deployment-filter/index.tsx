import {
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Select,
  TextField,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FC, memo, useCallback } from "react";
import { useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../../constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "../../constants/deployment-status-text";
import { AppState } from "../../modules";
import {
  Application,
  ApplicationKind,
  ApplicationKindKey,
  selectAll as selectAllApplications,
  selectById as selectApplicationById,
} from "../../modules/applications";
import {
  DeploymentStatus,
  DeploymentStatusKey,
  DeploymentFilterOptions,
} from "../../modules/deployments";
import { selectAllEnvs } from "../../modules/environments";
import { FilterView } from "../filter-view";

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

export interface DeploymentFilterProps {
  options: DeploymentFilterOptions;
  onClear: () => void;
  onChange: (options: DeploymentFilterOptions) => void;
}

export const DeploymentFilter: FC<DeploymentFilterProps> = memo(
  function DeploymentFilter({ options, onChange, onClear }) {
    const classes = useStyles();
    const envs = useSelector(selectAllEnvs);
    const applications = useSelector<AppState, Application.AsObject[]>(
      (state) => selectAllApplications(state.applications)
    );
    const selectedApp = useSelector<AppState, Application.AsObject | undefined>(
      (state) =>
        options.applicationId
          ? selectApplicationById(state.applications, options.applicationId)
          : undefined
    );
    const handleUpdateFilterValue = useCallback(
      (opts: Partial<DeploymentFilterOptions>): void => {
        onChange({ ...options, ...opts });
      },
      [options, onChange]
    );

    return (
      <FilterView
        onClear={() => {
          onClear();
        }}
      >
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
          <InputLabel id="filter-application-kind">Application Kind</InputLabel>
          <Select
            labelId="filter-application-kind"
            id="filter-application-kind"
            value={options.kind ?? ALL_VALUE}
            label="Application Kind"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                kind:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : `${e.target.value}`,
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
                applicationId: value ? value.id : undefined,
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
          <InputLabel id="filter-deployment-status">
            Deployment Status
          </InputLabel>
          <Select
            labelId="filter-deployment-status"
            id="filter-deployment-status"
            value={options.status ?? ALL_VALUE}
            label="Deployment Status"
            className={classes.select}
            onChange={(e) => {
              handleUpdateFilterValue({
                status:
                  e.target.value === ALL_VALUE
                    ? undefined
                    : `${e.target.value}`,
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
      </FilterView>
    );
  }
);
