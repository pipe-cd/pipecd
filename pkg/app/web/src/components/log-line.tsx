import { makeStyles, Box } from "@material-ui/core";
import React, { FC } from "react";
import { parseLog } from "../utils/parse-log";
import {
  TERM_COLORS,
  DEFAULT_BACKGROUND_COLOR,
  SELECTED_BACKGROUND_COLOR,
} from "../constants/term-colors";

const useStyles = makeStyles((theme) => ({
  container: {
    display: "flex",
    alignItems: "flex-start",
    backgroundColor: DEFAULT_BACKGROUND_COLOR,
    "&:hover": {
      backgroundColor: SELECTED_BACKGROUND_COLOR,
    },
  },
  lineNumber: {
    color: theme.palette.primary.main,
    width: "5rem",
    textAlign: "center",
    flexShrink: 0,
    userSelect: "none",
    cursor: "pointer",
  },
}));

interface Props {
  lineNumber: number;
  body: string;
}

export const LogLine: FC<Props> = ({ body, lineNumber }) => {
  const classes = useStyles();

  return (
    <div className={classes.container}>
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
