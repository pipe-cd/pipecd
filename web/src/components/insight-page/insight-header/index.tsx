import DayJSUtils from "@date-io/dayjs";
import { Box, makeStyles, TextField } from "@material-ui/core";
import { Autocomplete } from "@material-ui/lab";
import { DatePicker, MuiPickersUtilsProvider } from "@material-ui/pickers";
import { MaterialUiPickersDate } from "@material-ui/pickers/typings/date";
import { FC, memo, useCallback, useEffect } from "react";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { Application, selectAll, selectById } from "~/modules/applications";
import { fetchDeploymentFrequency } from "~/modules/deployment-frequency";
import {
  changeApplication,
  changeRangeFrom,
  changeRangeTo,
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

  const [applicationId, rangeFrom, rangeTo] = useAppSelector<
    [string, number, number]
  >((state) => [
    state.insight.applicationId,
    state.insight.rangeFrom,
    state.insight.rangeTo,
  ]);

  const selectedApp = useAppSelector<Application.AsObject | null>(
    (state) => selectById(state.applications, applicationId) || null
  );
  const applications = useAppSelector<Application.AsObject[]>((state) =>
    selectAll(state.applications)
  );

  const Picker = DatePicker;

  const handleApplicationChange = useCallback(
    (_, newValue: Application.AsObject | null) => {
      if (newValue) {
        dispatch(changeApplication(newValue.id));
      } else {
        dispatch(changeApplication(""));
      }
    },
    [dispatch]
  );

  const handleRangeFromChange = useCallback(
    (time: MaterialUiPickersDate) => {
      if (time) {
        dispatch(changeRangeFrom(time.valueOf()));
      }
    },
    [dispatch]
  );

  const handleRangeToChange = useCallback(
    (time: MaterialUiPickersDate) => {
      if (time) {
        dispatch(changeRangeTo(time.valueOf()));
      }
    },
    [dispatch]
  );

  useEffect(() => {
    dispatch(fetchDeploymentFrequency());
  }, [dispatch, applicationId, rangeFrom, rangeTo]);

  return (
    <Box display="flex" alignItems="center" justifyContent="flex-end">
      <Autocomplete
        id="application"
        style={{ width: 300 }}
        value={selectedApp}
        options={applications}
        getOptionLabel={(option) => option.name}
        onChange={handleApplicationChange}
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

      <MuiPickersUtilsProvider utils={DayJSUtils}>
        <Picker
          views={["date"]}
          margin="dense"
          inputVariant="outlined"
          variant="dialog"
          label="From"
          value={rangeFrom}
          onChange={handleRangeFromChange}
          className={classes.headerItemMargin}
        />
        <Picker
          views={["date"]}
          margin="dense"
          inputVariant="outlined"
          variant="dialog"
          label="To"
          value={rangeTo}
          onChange={handleRangeToChange}
          className={classes.rangeMargin}
        />
      </MuiPickersUtilsProvider>
    </Box>
  );
});
