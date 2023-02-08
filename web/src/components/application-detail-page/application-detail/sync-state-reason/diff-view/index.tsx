import { makeStyles } from "@material-ui/core";
import green from "@material-ui/core/colors/green";
import red from "@material-ui/core/colors/red";
import yellow from "@material-ui/core/colors/yellow";
import { FC, memo } from "react";

const useStyles = makeStyles((theme) => ({
  root: {
    fontFamily: theme.typography.fontFamilyMono,
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
  change: {
    color: yellow[800],
    backgroundColor: yellow[50],
  },
  line: {
    minHeight: `${theme.typography.body2.lineHeight}em`,
  },
}));

export interface DiffViewProps {
  content: string;
}

export const DiffView: FC<DiffViewProps> = memo(function DiffView({ content }) {
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
          case "~":
            return (
              <div key={i} className={classes.line} data-testid="changed-line">
                <span className={classes.change}>{line}</span>
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
