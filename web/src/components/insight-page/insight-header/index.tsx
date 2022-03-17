// import DayJSUtils from "@date-io/dayjs";
// import { Box, makeStyles, TextField } from "@material-ui/core";
// import {
//   Autocomplete,
//   ToggleButton,
//   ToggleButtonGroup,
// } from "@material-ui/lab";
// import {
//   DatePicker,
//   DatePickerView,
//   MuiPickersUtilsProvider,
// } from "@material-ui/pickers";
// import { MaterialUiPickersDate } from "@material-ui/pickers/typings/date";
// import { FC, memo, useCallback } from "react";
// import { INSIGHT_STEP_TEXT } from "~/constants/insight-step-text";
// import { useAppDispatch, useAppSelector } from "~/hooks/redux";
// import { Application, selectAll, selectById } from "~/modules/applications";
// import {
//   changeApplication,
//   changeRangeFrom,
//   changeRangeTo,
//   changeStep,
//   InsightStep,
// } from "~/modules/insight";
// import { WeekPicker } from "./week-picker";

// const useStyles = makeStyles((theme) => ({
//   headerItemMargin: {
//     marginLeft: theme.spacing(2),
//   },
//   rangeMargin: {
//     marginLeft: theme.spacing(1),
//   },
// }));

// const viewsMap: Record<InsightStep, DatePickerView[]> = {
//   [InsightStep.DAILY]: ["date"],
//   [InsightStep.WEEKLY]: ["date"],
//   [InsightStep.MONTHLY]: ["year", "month"],
//   [InsightStep.YEARLY]: ["year"],
// };

// export const InsightHeader: FC = memo(function InsightHeader() {
//   const classes = useStyles();
//   const dispatch = useAppDispatch();

//   const [applicationId, step, rangeFrom, rangeTo] = useAppSelector<
//     [string, InsightStep, number, number]
//   >((state) => [
//     state.insight.applicationId,
//     state.insight.step,
//     state.insight.rangeFrom,
//     state.insight.rangeTo,
//   ]);

//   const selectedApp = useAppSelector<Application.AsObject | null>(
//     (state) => selectById(state.applications, applicationId) || null
//   );
//   const applications = useAppSelector<Application.AsObject[]>((state) =>
//     selectAll(state.applications)
//   );

//   const views = viewsMap[step];
//   const Picker = step === InsightStep.WEEKLY ? WeekPicker : DatePicker;

//   const handleApplicationChange = useCallback(
//     (_, newValue: Application.AsObject | null) => {
//       if (newValue) {
//         dispatch(changeApplication(newValue.id));
//       } else {
//         dispatch(changeApplication(""));
//       }
//     },
//     [dispatch]
//   );

//   const handleRangeFromChange = useCallback(
//     (time: MaterialUiPickersDate) => {
//       if (time) {
//         dispatch(changeRangeFrom(time.valueOf()));
//       }
//     },
//     [dispatch]
//   );

//   const handleRangeToChange = useCallback(
//     (time: MaterialUiPickersDate) => {
//       if (time) {
//         dispatch(changeRangeTo(time.valueOf()));
//       }
//     },
//     [dispatch]
//   );

//   const handleStepChange = useCallback(
//     (_, value) => {
//       if (value !== null) {
//         dispatch(changeStep(value));
//       }
//     },
//     [dispatch]
//   );

//   // TODO: Enable fetch chart data on insight filter changes.
//   // useEffect(() => {
//   //   dispatch(fetchDeploymentFrequency());
//   // }, [dispatch, applicationId, step, rangeFrom, rangeTo]);

//   return (
//     <Box display="flex" alignItems="center" justifyContent="flex-end">
//       <Autocomplete
//         id="application"
//         style={{ width: 300 }}
//         value={selectedApp}
//         options={applications}
//         getOptionLabel={(option) => option.name}
//         onChange={handleApplicationChange}
//         renderInput={(params) => (
//           <TextField
//             {...params}
//             label="Application"
//             margin="dense"
//             variant="outlined"
//             required
//           />
//         )}
//       />

//       <ToggleButtonGroup
//         value={step}
//         onChange={handleStepChange}
//         exclusive
//         aria-label="insight step"
//         size="medium"
//         className={classes.headerItemMargin}
//       >
//         <ToggleButton
//           value={InsightStep.DAILY}
//           aria-label={INSIGHT_STEP_TEXT[InsightStep.DAILY]}
//         >
//           {INSIGHT_STEP_TEXT[InsightStep.DAILY]}
//         </ToggleButton>
//         <ToggleButton
//           value={InsightStep.WEEKLY}
//           aria-label={INSIGHT_STEP_TEXT[InsightStep.WEEKLY]}
//         >
//           {INSIGHT_STEP_TEXT[InsightStep.WEEKLY]}
//         </ToggleButton>
//         <ToggleButton
//           value={InsightStep.MONTHLY}
//           aria-label={INSIGHT_STEP_TEXT[InsightStep.MONTHLY]}
//         >
//           {INSIGHT_STEP_TEXT[InsightStep.MONTHLY]}
//         </ToggleButton>
//         <ToggleButton
//           value={InsightStep.YEARLY}
//           aria-label={INSIGHT_STEP_TEXT[InsightStep.YEARLY]}
//         >
//           {INSIGHT_STEP_TEXT[InsightStep.YEARLY]}
//         </ToggleButton>
//       </ToggleButtonGroup>

//       <MuiPickersUtilsProvider utils={DayJSUtils}>
//         <Picker
//           views={views}
//           margin="dense"
//           inputVariant="outlined"
//           variant="dialog"
//           label="From"
//           value={rangeFrom}
//           onChange={handleRangeFromChange}
//           className={classes.headerItemMargin}
//         />
//         <Picker
//           views={views}
//           margin="dense"
//           inputVariant="outlined"
//           variant="dialog"
//           label="To"
//           value={rangeTo}
//           onChange={handleRangeToChange}
//           className={classes.rangeMargin}
//         />
//       </MuiPickersUtilsProvider>
//     </Box>
//   );
// });
