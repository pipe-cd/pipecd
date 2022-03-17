// import { DateType } from "@date-io/type";
// import { IconButton, makeStyles } from "@material-ui/core";
// import { DatePicker } from "@material-ui/pickers";
// import { MaterialUiPickersDate } from "@material-ui/pickers/typings/date";
// import clsx from "clsx";
// import { FC, memo, useCallback } from "react";

// const useStyles = makeStyles((theme) => ({
//   ayWrapper: {
//     position: "relative",
//   },
//   day: {
//     width: 36,
//     height: 36,
//     fontSize: theme.typography.caption.fontSize,
//     margin: "0 2px",
//     color: "inherit",
//   },
//   customDayHighlight: {
//     position: "absolute",
//     top: 0,
//     bottom: 0,
//     left: "2px",
//     right: "2px",
//     border: `1px solid ${theme.palette.secondary.main}`,
//     borderRadius: "50%",
//   },
//   nonCurrentMonthDay: {
//     color: theme.palette.text.disabled,
//   },
//   highlightNonCurrentMonthDay: {
//     color: "#999999",
//   },
//   highlight: {
//     background: theme.palette.primary.main,
//     color: theme.palette.common.white,
//   },
//   firstHighlight: {
//     extend: "highlight",
//     borderTopLeftRadius: "50%",
//     borderBottomLeftRadius: "50%",
//   },
//   endHighlight: {
//     extend: "highlight",
//     borderTopRightRadius: "50%",
//     borderBottomRightRadius: "50%",
//   },
// }));

// const formatWeekSelectLabel = (
//   date: MaterialUiPickersDate,
//   invalidLabel: string
// ): string => {
//   return date ? `Week of ${date.day(0).format("MMM Do")}` : invalidLabel;
// };

// export interface WeekPickerProps {
//   value: Date | number | null;
//   label: string;
//   onChange: (date: MaterialUiPickersDate) => void;
//   className?: string;
// }

// export const WeekPicker: FC<WeekPickerProps> = memo(function WeekPicker({
//   value,
//   label,
//   onChange,
//   className,
// }) {
//   const classes = useStyles();

//   const renderWeekDay = useCallback(
//     (date: DateType, selected: DateType, dayInCurrentMonth: boolean) => {
//       const start = selected.day(0);
//       const end = selected.day(6);

//       const dayIsBetween = date.isBetween(start, end, "day", "[]");
//       const isFirstDay = date.isSame(start, "day");
//       const isLastDay = date.isSame(end, "day");

//       const wrapperClassName = clsx({
//         [classes.highlight]: dayIsBetween,
//         [classes.firstHighlight]: isFirstDay,
//         [classes.endHighlight]: isLastDay,
//       });

//       const dayClasses = clsx(classes.day, {
//         [classes.nonCurrentMonthDay]: !dayInCurrentMonth,
//         [classes.highlightNonCurrentMonthDay]:
//           !dayInCurrentMonth && dayIsBetween,
//       });

//       return (
//         <div className={wrapperClassName}>
//           <IconButton className={dayClasses}>
//             <span>{date.format("D")}</span>
//           </IconButton>
//         </div>
//       );
//     },
//     [classes]
//   );

//   return (
//     <DatePicker
//       margin="dense"
//       inputVariant="outlined"
//       variant="dialog"
//       label={label}
//       onChange={onChange}
//       value={value}
//       renderDay={renderWeekDay as any}
//       labelFunc={formatWeekSelectLabel}
//       className={className}
//     />
//   );
// });
