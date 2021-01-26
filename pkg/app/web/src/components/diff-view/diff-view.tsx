import React, { FC, memo } from "react";
import { makeStyles } from "@material-ui/core";
import red from "@material-ui/core/colors/red";
import green from "@material-ui/core/colors/green";

const useStyles = makeStyles((theme) => ({
  root: {
    fontFamily: "Roboto Mono",
    wordBreak: "break-all",
    whiteSpace: "pre-wrap",
  },
  add: {
    color: green[800],
    backgroundColor: green[50],
  },
  del: {
    color: red[800],
    backgroundColor: red[50],
  },
  line: {
    minHeight: `${theme.typography.body2.lineHeight}em`,
  },
}));

interface Props {
  content: string;
}

export const DiffView: FC<Props> = memo(function DiffView({ content }) {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {content.split("\n").map((line, i) => {
        switch (line[0]) {
          case "+":
            return (
              <div key={i} className={classes.line} data-testid="added-line">
                <span key={i} className={classes.add}>
                  {line}
                </span>
              </div>
            );
          case "-":
            return (
              <div key={i} className={classes.line} data-testid="deleted-line">
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
});
