import {
  Box,
  Grid,
  TextField,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  makeStyles,
} from "@material-ui/core";
import { Autocomplete } from "@material-ui/lab";
import { FC, memo, useState, useEffect } from "react";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { Application, selectAll, selectById } from "~/modules/applications";

import {
  changeRange,
  changeResolution,
  changeApplication,
  changeLabels,
  InsightResolutions,
  InsightRanges,
  INSIGHT_RESOLUTION_TEXT,
  INSIGHT_RANGE_TEXT,
  InsightResolution,
  InsightRange,
} from "~/modules/insight";

const useStyles = makeStyles((theme) => ({
  headerItemMargin: {
    marginLeft: theme.spacing(2),
  },
  rangeMargin: {
    marginLeft: theme.spacing(1),
  },
}));

export const InsightHeader: FC = memo(function InsightHeader() {
  const classes = useStyles();
  const dispatch = useAppDispatch();

  const selectedApp = useAppSelector<Application.AsObject | null>(
    (state) =>
      selectById(state.applications, state.insight.applicationId) || null
  );

  const [allLabels, setAllLabels] = useState(new Array<string>());
  const [selectedLabels, setSelectedLabels] = useState(new Array<string>());

  const [selectedRange, setSelectedRange] = useState(InsightRange.LAST_1_MONTH);
  const [selectedResolution, setSelectedResolution] = useState(
    InsightResolution.DAILY
  );

  const applications = useAppSelector<Application.AsObject[]>((state) =>
    selectAll(state.applications)
  );

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
    <Grid container spacing={2} style={{ marginTop: 26, marginBottom: 26 }}>
      <Grid item xs={8}>
        <Box display="flex" alignItems="left" justifyContent="flex-start">
          <Autocomplete
            id="application"
            style={{ minWidth: 300 }}
            value={selectedApp}
            options={applications}
            getOptionLabel={(option) => option.name}
            onChange={(_, value) => {
              if (value) {
                dispatch(changeApplication(value.id));
              } else {
                dispatch(changeApplication(""));
              }
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
            className={classes.headerItemMargin}
            options={allLabels}
            value={selectedLabels}
            onInputChange={(_, value) => {
              const label = value.split(":");
              if (label.length !== 2) return;
              if (label[0].length === 0) return;
              if (label[1].length === 0) return;
              setAllLabels([value]);
            }}
            onChange={(_, value) => {
              setSelectedLabels(value);
              dispatch(changeLabels(value));
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
      <Grid item xs={4}>
        <Box display="flex" alignItems="right" justifyContent="flex-end">
          <FormControl className={classes.headerItemMargin} variant="outlined">
            <InputLabel id="range-input">Range</InputLabel>
            <Select
              id="range"
              label="Range"
              value={selectedRange}
              onChange={(e) => {
                const value = e.target.value as InsightRange;
                setSelectedRange(value);
                dispatch(changeRange(value));
              }}
            >
              {InsightRanges.map((e) => (
                <MenuItem key={e} value={e}>
                  {INSIGHT_RANGE_TEXT[e]}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          <FormControl className={classes.headerItemMargin} variant="outlined">
            <InputLabel id="resolution-input">Resolution</InputLabel>
            <Select
              id="resolution"
              label="Resolution"
              value={selectedResolution}
              onChange={(e) => {
                const value = e.target.value as InsightResolution;
                setSelectedResolution(value);
                dispatch(changeResolution(value));
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
