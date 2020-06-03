import React from "react";
import { theme } from "../src/theme";
import CssBaseline from "@material-ui/core/CssBaseline";
import { makeStyles, ThemeProvider } from "@material-ui/core";

// NOTE: To prevent difference of screenshot that is caused by mouse hover.
const useStyles = makeStyles({
  wrapper: {
    padding: 10
  }
});

export const ThemeDecorator = (fn: () => JSX.Element) => {
  const classes = useStyles();
  return (
    <ThemeProvider theme={theme}>
      <div className={classes.wrapper}>
        <CssBaseline />
        {fn()}
      </div>
    </ThemeProvider>
  );
};
