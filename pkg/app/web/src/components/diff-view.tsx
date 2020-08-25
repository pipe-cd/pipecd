import React, { FC } from "react";
import { makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() => ({
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
}));

interface Props {
  content: string;
}

export const DiffView: FC<Props> = ({ content }) => {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      {content.split("\n").map((line, i) => {
        console.log(line, line[0]);
        switch (line[0]) {
          case "+":
            return (
              <div key={i} className={classes.add}>
                {line}
              </div>
            );
          case "-":
            return (
              <div key={i} className={classes.del}>
                {line}
              </div>
            );
          default:
            return <div key={i}>{line}</div>;
        }
      })}
    </div>
  );
};
