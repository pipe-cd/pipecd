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
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { useAppSelector } from "~/hooks/redux";
import {
  Application,
  ApplicationKind,
  ApplicationKindKey,
  selectAll as selectAllApplications,
  selectById as selectApplicationById,
} from "~/modules/applications";
import {
  DeploymentFilterOptions,
  DeploymentStatus,
  DeploymentStatusKey,
} from "~/modules/deployments";
import { selectAllEnvs } from "~/modules/environments";
import { ApplicationAutocomplete } from "../../applications-page/application-filter/application-autocomplete";

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
    const envs = useAppSelector(selectAllEnvs);
    const [localApplications, setLocalApplications] = useState<Application.AsObject[]>([]);
    const applications = useAppSelector<Application.AsObject[]>((state) =>
      selectAllApplications(state.applications)
    );
    const selectedApp = useAppSelector<Application.AsObject | undefined>(
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

    useEffect(() => {
      if (options.applicationName) {
        setLocalApplications(applications.filter(app => app.name === options.applicationName));
      } else {
        setLocalApplications(applications);
      }
    }, [applications, options]);

    return (
      <FilterView
        onClear={() => {
          onClear();
        }}
      >
        <div className={classes.formItem}>
          <ApplicationAutocomplete
            value={options.applicationName ?? null}
            onChange={(value) => handleUpdateFilterValue({ applicationName: value })}
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
            options={localApplications}
            getOptionLabel={(option) => option.id}
            renderOption={(option) => <span>{option.name} ({option.id})</span>}
            value={selectedApp || null}
            onChange={(_, value) => {
              handleUpdateFilterValue({
                applicationId: value ? value.id : undefined,
              });
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                label="Application Id"
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
