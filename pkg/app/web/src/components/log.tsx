import React, { FC } from "react";
import { makeStyles, CircularProgress, Box } from "@material-ui/core";
import { LogLine } from "./log-line";
import { DEFAULT_BACKGROUND_COLOR } from "../constants/term-colors";

const useStyles = makeStyles((theme) => ({
  container: {
    fontFamily: "Roboto Mono",
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    paddingTop: theme.spacing(1),
    paddingBottom: theme.spacing(1),
  },
}));

interface Props {
  logs: string[];
  loading: boolean;
}

export const Log: FC<Props> = ({ logs, loading }) => {
  const classes = useStyles();
  return (
    <div className={classes.container}>
      {logs.map((body, i) => (
        <LogLine body={body} lineNumber={i + 1} />
      ))}
      {loading && (
        <Box display="flex" justifyContent="center" p={1}>
          <CircularProgress />
        </Box>
      )}
    </div>
  );
};
