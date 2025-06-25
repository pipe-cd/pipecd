import { FC, memo, useEffect, useRef } from "react";
import { CircularProgress, Box } from "@mui/material";
import { LogLine } from "../log-line";
import { DEFAULT_BACKGROUND_COLOR } from "~/constants/term-colors";
import { LogBlock } from "~~/model/logblock_pb";

export interface LogProps {
  logs: LogBlock.AsObject[];
  loading: boolean;
}

export const Log: FC<LogProps> = memo(function Log({ logs, loading }) {
  const bottomRef = useRef<null | HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <Box
      sx={(theme) => ({
        fontFamily: theme.typography.fontFamilyMono,
        backgroundColor: DEFAULT_BACKGROUND_COLOR,
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(1),
        height: "100%",
      })}
    >
      {logs.map((log, i) => (
        <LogLine
          key={`log-${log.index}`}
          severity={log.severity}
          body={log.log}
          lineNumber={i + 1}
          createdAt={log.createdAt}
        />
      ))}
      {loading && (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            p: 1,
          }}
        >
          <CircularProgress color="secondary" />
        </Box>
      )}
      <Box
        sx={(theme) => ({
          height: theme.spacing(1),
          backgroundColor: DEFAULT_BACKGROUND_COLOR,
        })}
        ref={bottomRef}
      />
    </Box>
  );
});
