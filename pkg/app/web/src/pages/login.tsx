import { makeStyles, TextField, Button } from "@material-ui/core";
import React, { FC, memo } from "react";
import { STATIC_LOGIN_ENDPOINT } from "../constants";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
  },
  form: {
    display: "flex",
    flexDirection: "column",
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
  },
}));

export const LoginPage: FC = memo(function LoginPage() {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <form
        method="POST"
        action={STATIC_LOGIN_ENDPOINT}
        className={classes.form}
      >
        <TextField
          id="project"
          name="project"
          label="Project"
          variant="outlined"
          margin="dense"
          required
        />
        <TextField
          id="username"
          name="username"
          label="Username"
          variant="outlined"
          margin="dense"
          required
        />
        <TextField
          id="password"
          name="password"
          label="Password"
          type="password"
          variant="outlined"
          margin="dense"
          required
        />
        <div className={classes.buttons}>
          <Button type="submit" color="primary">
            Log In
          </Button>
        </div>
      </form>
    </div>
  );
});
