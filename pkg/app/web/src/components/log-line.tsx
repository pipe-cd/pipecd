import { Box, makeStyles } from "@material-ui/core";
import { Error } from "@material-ui/icons";
import React, { FC } from "react";
import {
  DEFAULT_BACKGROUND_COLOR,
  SELECTED_BACKGROUND_COLOR,
  TERMINAL_LINE_NUMBER_COLOR,
  TERM_COLORS,
} from "../constants/term-colors";
import { LogSeverity } from "../modules/stage-logs";
import { parseLog } from "../utils/parse-log";

const useStyles = makeStyles((theme) => ({
  container: {
    display: "flex",
    alignItems: "flex-start",
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    "&:hover": {
      backgroundColor: SELECTED_BACKGROUND_COLOR,
    },
    position: "relative",
  },
  lineNumber: {
    color: TERMINAL_LINE_NUMBER_COLOR,
    width: "5rem",
    textAlign: "center",
    flexShrink: 0,
    userSelect: "none",
    cursor: "pointer",
  },
  icon: {
    position: "absolute",
    marginLeft: theme.spacing(1),
  },
}));

interface Props {
  lineNumber: number;
  body: string;
  severity: LogSeverity;
}

export const LogLine: FC<Props> = ({ body, lineNumber, severity }) => {
  const classes = useStyles();

  return (
    <div className={classes.container}>
      {severity === LogSeverity.ERROR && (
        <Error color="error" fontSize="small" className={classes.icon} />
      )}
      <span className={classes.lineNumber}>{lineNumber}</span>
      <Box pr={2} flex={1} style={{ wordBreak: "break-word" }}>
        {parseLog(body).map((cell) => (
          <span
            style={{
              color: TERM_COLORS[cell.fg],
              backgroundColor: cell.bg !== 0 ? TERM_COLORS[cell.bg] : undefined,
              fontWeight: cell.bold ? "bold" : undefined,
              textDecoration: cell.underline ? "underline" : undefined,
            }}
          >
            {cell.content}
          </span>
        ))}
      </Box>
    </div>
  );
};
