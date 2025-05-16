import { Box } from "@mui/material";
import { Error } from "@mui/icons-material";
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
  return (
    <Box
      sx={(theme) => ({
        display: "flex",
        alignItems: "flex-start",
        backgroundColor: DEFAULT_BACKGROUND_COLOR,
        "&:hover": {
          backgroundColor: SELECTED_BACKGROUND_COLOR,
        },
        position: "relative",
        fontSize: theme.typography.body2.fontSize,
      })}
    >
      {severity === LogSeverity.ERROR && (
        <Error
          color="error"
          fontSize="small"
          sx={{
            position: "absolute",
            marginLeft: 1,
          }}
        />
      )}
      <Box
        sx={{
          color: TERMINAL_LINE_NUMBER_COLOR,
          width: "5rem",
          textAlign: "center",
          flexShrink: 0,
          userSelect: "none",
          cursor: "pointer",
        }}
      >
        {lineNumber}
      </Box>
      <Box
        sx={(theme) => ({
          color: DEFAULT_TEXT_COLOR,
          paddingRight: theme.spacing(1),
          opacity: 0.8,
        })}
      >{`[${dayjs(createdAt * 1000).format(TIMESTAMP_FORMAT)}]`}</Box>
      <Box
        sx={{
          wordBreak: "break-all",
          pr: 2,
          flex: 1,
        }}
      >
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
    </Box>
  );
};
