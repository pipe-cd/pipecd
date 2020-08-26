import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles((theme) => ({
  root: {
    fontFamily: "Roboto Mono",
    wordBreak: "break-all",
    whiteSpace: "pre-wrap",
  },
  add: {
    color: "#22863a",
    backgroundColor: "#f0fff4",
  },
  del: {
    color: "#b31d28",
    backgroundColor: "#ffeef0",
  },
  line: {
    minHeight: `${theme.typography.body2.lineHeight}em`,
  },
}));

interface Props {
  content: string;
}

export const DiffView: FC<Props> = ({ content }) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {content.split("\n").map((line, i) => {
        switch (line[0]) {
          case "+":
            return (
              <div key={i} className={classes.line}>
                <span key={i} className={classes.add}>
                  {line}
                </span>
              </div>
            );
          case "-":
            return (
              <div key={i} className={classes.line}>
                <span className={classes.del}>{line}</span>
              </div>
            );
          default:
            return (
              <div key={i} className={classes.line}>
                {line}
              </div>
            );
        }
      })}
    </div>
  );
};
