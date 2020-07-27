import { Button, makeStyles, TextField } from "@material-ui/core";
import React, { FC, memo } from "react";
import { Redirect } from "react-router";
import { PAGE_PATH_APPLICATIONS, STATIC_LOGIN_ENDPOINT } from "../constants";
import { useMe } from "../modules/me";

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
  const me = useMe();

  return (
    <div className={classes.root}>
      {me && me.isLogin && <Redirect to={PAGE_PATH_APPLICATIONS} />}
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
