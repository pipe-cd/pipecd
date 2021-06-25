import { Box, makeStyles } from "@material-ui/core";
import { Error } from "@material-ui/icons";
import { FC } from "react";
import {
  DEFAULT_BACKGROUND_COLOR,
  SELECTED_BACKGROUND_COLOR,
  TERMINAL_LINE_NUMBER_COLOR,
  DEFAULT_TEXT_COLOR,
  TERM_COLORS,
} from "~/constants/term-colors";
import { LogSeverity } from "~/modules/stage-logs";
import { parseLog } from "~/utils/parse-log";
import dayjs from "dayjs";

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
  timestamp: {
    color: DEFAULT_TEXT_COLOR,
    paddingRight: theme.spacing(1),
    opacity: 0.8,
  },
}));

export interface LogLineProps {
  lineNumber: number;
  body: string;
  severity: LogSeverity;
  createdAt: number;
}

const TIMESTAMP_FORMAT = "YYYY-MM-DD HH:mm:ss Z";

export const LogLine: FC<LogLineProps> = ({
  body,
  lineNumber,
  severity,
  createdAt,
}) => {
  const classes = useStyles();

  return (
    <div className={classes.container}>
      {severity === LogSeverity.ERROR && (
        <Error color="error" fontSize="small" className={classes.icon} />
      )}
      <span className={classes.lineNumber}>{lineNumber}</span>
      <span className={classes.timestamp}>{`[${dayjs(createdAt * 1000).format(
        TIMESTAMP_FORMAT
      )}]`}</span>
      <Box pr={2} flex={1} style={{ wordBreak: "break-all" }}>
        {parseLog(body).map((cell, i) => (
          <span
            key={`log-cell-${i}`}
            style={{
              color: TERM_COLORS[cell.fg],
              backgroundColor: cell.bg !== 0 ? TERM_COLORS[cell.bg] : undefined,
              fontWeight: cell.bold ? "bold" : undefined,
              textDecoration: cell.underline ? "underline" : undefined,
              whiteSpace: "pre-wrap",
            }}
          >
            {cell.content.split("\\n").join("\n")}
          </span>
        ))}
      </Box>
    </div>
  );
};
