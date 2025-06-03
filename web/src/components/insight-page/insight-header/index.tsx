import {
  Box,
  Grid,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
} from "@mui/material";
import { Autocomplete } from "@mui/material";
import { FC, memo, useMemo } from "react";
import { useGetApplications } from "~/queries/applications/use-get-applications";
import { InsightFilterValues } from "..";
import {
  INSIGHT_RESOLUTION_TEXT,
  INSIGHT_RANGE_TEXT,
  InsightResolution,
  InsightRange,
  InsightRanges,
  InsightResolutions,
} from "~/queries/insight/insight.config";

type Props = {
  filterValues: InsightFilterValues;
  onChangeFilter: (filterValues: Partial<InsightFilterValues>) => void;
};

export const InsightHeader: FC<Props> = memo(function InsightHeader({
  onChangeFilter,
  filterValues,
}) {
  const { data: applications = [] } = useGetApplications();

  const allLabels = useMemo(() => {
    const labels = new Set<string>();
    applications
      .filter((app) => app.labelsMap.length > 0)
      .map((app) => {
        app.labelsMap.map((label) => {
          labels.add(`${label[0]}:${label[1]}`);
        });
      });
    return Array.from(labels);
  }, [applications]);

  return (
    <Grid container spacing={2} style={{ marginTop: 26, marginBottom: 26 }}>
      <Grid size={8}>
        <Box
          sx={{
            display: "flex",
            alignItems: "left",
            justifyContent: "flex-start",
          }}
        >
          <Autocomplete
            id="application"
            style={{ minWidth: 300 }}
            value={
              applications.find(
                (item) => item.id === filterValues?.applicationId
              ) ?? null
            }
            options={applications}
            getOptionKey={(option) => option.id}
            getOptionLabel={(option) => option.name}
            onChange={(_, value) => {
              onChangeFilter({ applicationId: value ? value.id : "" });
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                label="Application"
                margin="dense"
                variant="outlined"
                required
              />
            )}
          />

          <Autocomplete
            multiple
            autoHighlight
            id="labels"
            noOptionsText="No selectable labels"
            style={{ minWidth: 300 }}
            sx={{ ml: 2 }}
            options={allLabels}
            value={filterValues?.labels ?? []}
            onChange={(_, value) => {
              onChangeFilter({ labels: value });
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
        </Box>
      </Grid>
      <Grid size={4}>
        <Box
          sx={{
            display: "flex",
            alignItems: "right",
            justifyContent: "flex-end",
          }}
        >
          <FormControl sx={{ ml: 2 }} variant="outlined">
            <InputLabel id="range-input">Range</InputLabel>
            <Select
              id="range"
              label="Range"
              value={filterValues?.range}
              onChange={(e) => {
                const value = e.target.value as InsightRange;
                onChangeFilter({ range: value });
              }}
            >
              {InsightRanges.map((e) => (
                <MenuItem key={e} value={e}>
                  {INSIGHT_RANGE_TEXT[e]}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl sx={{ ml: 2 }} variant="outlined">
            <InputLabel id="resolution-input">Resolution</InputLabel>
            <Select
              id="resolution"
              label="Resolution"
              value={filterValues?.resolution}
              onChange={(e) => {
                const value = e.target.value as InsightResolution;
                onChangeFilter({ resolution: value });
              }}
            >
              {InsightResolutions.map((e) => (
                <MenuItem key={e} value={e}>
                  {INSIGHT_RESOLUTION_TEXT[e]}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Box>
      </Grid>
    </Grid>
  );
});
