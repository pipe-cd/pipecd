import { FC, memo, useEffect, useRef } from "react";
import { makeStyles, CircularProgress, Box } from "@material-ui/core";
import { LogLine } from "../log-line";
import { DEFAULT_BACKGROUND_COLOR } from "~/constants/term-colors";
import { LogBlock } from "~/modules/stage-logs";

const useStyles = makeStyles((theme) => ({
  container: {
    fontFamily: theme.typography.fontFamilyMono,
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
    height: "100%",
  },
  space: {
    height: theme.spacing(1),
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
  },
}));

export interface LogProps {
  logs: LogBlock.AsObject[];
  loading: boolean;
}

export const Log: FC<LogProps> = memo(function Log({ logs, loading }) {
  const classes = useStyles();
  const bottomRef = useRef<null | HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [logs]);

  return (
    <div className={classes.container}>
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
        <Box display="flex" justifyContent="center" p={1}>
          <CircularProgress color="secondary" />
        </Box>
      )}
      <div className={classes.space} ref={bottomRef} />
    </div>
  );
});
