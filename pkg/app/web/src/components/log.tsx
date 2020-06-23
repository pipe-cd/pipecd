import React, { FC, memo } from "react";
import { makeStyles, CircularProgress, Box } from "@material-ui/core";
import { LogLine } from "./log-line";
import { DEFAULT_BACKGROUND_COLOR } from "../constants/term-colors";
import { LogBlock } from "../modules/stage-logs";

const useStyles = makeStyles((theme) => ({
  container: {
    fontFamily: "Roboto Mono",
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
    height: (props: { height?: number }) =>
      props.height ? props.height : undefined,
  },
}));

interface Props {
  logs: LogBlock[];
  loading: boolean;
  height?: number;
}

export const Log: FC<Props> = memo(function Log({ logs, loading, height }) {
  const classes = useStyles({ height });
  return (
    <div className={classes.container}>
      {logs.map((log, i) => (
        <LogLine
          key={`log-${log.index}`}
          severity={log.severity}
          body={log.log}
          lineNumber={i + 1}
        />
      ))}
      {loading && (
        <Box display="flex" justifyContent="center" p={1}>
          <CircularProgress />
        </Box>
      )}
    </div>
  );
});
