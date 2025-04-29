import {
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import Autocomplete from "@mui/material/Autocomplete";
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
    const [localApplications, setLocalApplications] = useState<
      Application.AsObject[]
    >([]);
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
        setLocalApplications(
          applications.filter((app) => app.name === options.applicationName)
        );
      } else {
        setLocalApplications(applications);
      }
    }, [applications, options]);

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
            value={options.applicationName ?? null}
            onChange={(value) =>
              handleUpdateFilterValue({ applicationName: value })
            }
          />
        </div>

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
            renderOption={(props, option) => (
              // TODO check this changes, add prop to span
              <span {...props}>
                {option.name} ({option.id})
              </span>
            )}
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
