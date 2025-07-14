import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from "@mui/material";
import Autocomplete from "@mui/material/Autocomplete";
import { FC, memo, useCallback, useMemo } from "react";
import { FilterView } from "~/components/filter-view";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { DEPLOYMENT_STATE_TEXT } from "~/constants/deployment-status-text";
import { ApplicationAutocomplete } from "../../applications-page/application-filter/application-autocomplete";
import { ApplicationKindKey } from "~/types/applications";
import { ApplicationKind } from "~~/model/common_pb";
import { DeploymentStatus, DeploymentStatusKey } from "~/types/deployment";
import { DeploymentFilterOptions } from "~/queries/deployment/use-get-deployments-infinite";
import { useGetApplications } from "~/queries/applications/use-get-applications";
import LabelAutoComplete from "~/components/label-auto-complete";

const ALL_VALUE = "ALL";

export interface DeploymentFilterProps {
  options: DeploymentFilterOptions;
  onClear: () => void;
  onChange: (options: DeploymentFilterOptions) => void;
}

export const DeploymentFilter: FC<DeploymentFilterProps> = memo(
  function DeploymentFilter({ options, onChange, onClear }) {
    const { data: applications = [] } = useGetApplications();

    const selectedApp = useMemo(() => {
      return applications.find((app) => app.id === options.applicationId);
    }, [applications, options.applicationId]);

    const handleUpdateFilterValue = useCallback(
      (opts: Partial<DeploymentFilterOptions>): void => {
        onChange({ ...options, ...opts });
      },
      [options, onChange]
    );

    const localApplications = useMemo(() => {
      if (!options.applicationName) return applications;

      return applications.filter((app) => app.name === options.applicationName);
    }, [applications, options.applicationName]);

    const allApplicationLabels = useMemo(() => {
      const labels = new Set<string>();
      applications.forEach((app) => {
        if (app.labelsMap.length > 0)
          app.labelsMap.forEach((label) => {
            labels.add(`${label[0]}:${label[1]}`);
          });
      });
      return Array.from(labels);
    }, [applications]);

    return (
      <FilterView
        onClear={() => {
          onClear();
        }}
      >
        <Box
          sx={{
            width: "100%",
            mt: 4,
          }}
        >
          <ApplicationAutocomplete
            value={options.applicationName ?? null}
            onChange={(value) =>
              handleUpdateFilterValue({
                applicationName: value,
                applicationId:
                  options.applicationName === value
                    ? options.applicationId
                    : undefined,
              })
            }
          />
        </Box>
        <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
          <InputLabel id="filter-application-kind">Application Kind</InputLabel>
          <Select
            labelId="filter-application-kind"
            id="filter-application-kind"
            value={options.kind ?? ALL_VALUE}
            label="Application Kind"
            fullWidth
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
        <Box sx={{ width: "100%", mt: 4 }}>
          <Autocomplete
            id="application-select"
            options={localApplications}
            getOptionLabel={(option) => option.id}
            renderOption={(props, option) => (
              <span {...props}>{`${option.name} (${option.id})`}</span>
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
                slotProps={{
                  htmlInput: {
                    ...params.inputProps,
                  },
                }}
              />
            )}
          />
        </Box>
        <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
          <InputLabel id="filter-deployment-status">
            Deployment Status
          </InputLabel>
          <Select
            labelId="filter-deployment-status"
            id="filter-deployment-status"
            value={options.status ?? ALL_VALUE}
            label="Deployment Status"
            fullWidth
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
        <FormControl sx={{ width: "100%", mt: 4 }} variant="outlined">
          <LabelAutoComplete
            value={options.labels ?? []}
            options={allApplicationLabels}
            onChange={(newLabels) => {
              handleUpdateFilterValue({
                labels: newLabels,
              });
            }}
          />
        </FormControl>
      </FilterView>
    );
  },
  (prevProps, nextProps) => {
    return (
      prevProps.options.applicationId === nextProps.options.applicationId &&
      prevProps.options.applicationName === nextProps.options.applicationName &&
      prevProps.options.kind === nextProps.options.kind &&
      prevProps.options.status === nextProps.options.status &&
      JSON.stringify(prevProps.options.labels) ===
        JSON.stringify(nextProps.options.labels)
    );
  }
);
