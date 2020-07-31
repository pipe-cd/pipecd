import React, { FC, memo } from "react";
import { makeStyles, TextField, Button, Typography } from "@material-ui/core";
import { STATIC_LOGIN_ENDPOINT } from "../constants";
import { useProjectName, clearProjectName } from "../modules/login";
import { useDispatch } from "react-redux";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    alignItems: "center",
    flexDirection: "column",
    flex: 1,
  },
  form: {
    display: "flex",
    flexDirection: "column",
    textAlign: "center",
    marginTop: theme.spacing(4),
    width: 320,
  },
  fields: {
    display: "flex",
    flexDirection: "column",
    marginTop: theme.spacing(4),
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
    marginTop: theme.spacing(3),
  },
}));

export const LoginForm: FC = memo(function LoginForm() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const projectName = useProjectName();

  const handleReset = (): void => {
    dispatch(clearProjectName());
  };

  return (
    <div className={classes.root}>
      <Typography variant="h4">Sign in to {projectName}</Typography>
      <form
        method="POST"
        action={STATIC_LOGIN_ENDPOINT}
        className={classes.form}
      >
        <input
          type="hidden"
          id="project"
          name="project"
          value={projectName || undefined}
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
          <Button type="reset" color="primary" onClick={handleReset}>
            back
          </Button>
          <Button type="submit" color="primary" variant="contained">
            login
          </Button>
        </div>
      </form>
    </div>
  );
});
