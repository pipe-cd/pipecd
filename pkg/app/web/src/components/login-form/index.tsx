import React, { FC, memo } from "react";
import { makeStyles, TextField, Button, Typography } from "@material-ui/core";
import { STATIC_LOGIN_ENDPOINT, LOGIN_ENDPOINT } from "../../constants/path";
import { useProjectName, clearProjectName } from "../../modules/login";
import { useDispatch } from "react-redux";
import { MarkGithubIcon } from "@primer/octicons-react";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    alignItems: "center",
    flexDirection: "column",
    flex: 1,
  },
  content: {
    width: 320,
  },
  form: {
    display: "flex",
    flexDirection: "column",
    textAlign: "center",
    marginTop: theme.spacing(4),
    width: "100%",
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
  githubMark: {
    marginRight: theme.spacing(1),
  },
  githubLoginButton: {
    background: "#24292E",
  },
  divider: {
    display: "flex",
    alignItems: "center",
    marginTop: theme.spacing(3),
    marginBottom: theme.spacing(3),
  },
  dividerText: {
    marginLeft: theme.spacing(2),
    marginRight: theme.spacing(2),
  },
  line: {
    flex: 1,
    border: "none",
    borderTop: "1px solid #ddd",
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
      <div className={classes.content}>
        <form method="POST" action={LOGIN_ENDPOINT} className={classes.form}>
          <input
            type="hidden"
            id="project"
            name="project"
            value={projectName || undefined}
          />
          <Button
            type="submit"
            color="primary"
            variant="contained"
            className={classes.githubLoginButton}
          >
            <MarkGithubIcon className={classes.githubMark} />
            LOGIN WITH GITHUB
          </Button>
        </form>

        <div className={classes.divider}>
          <span className={classes.line} />
          <div className={classes.dividerText}>OR</div>
          <span className={classes.line} />
        </div>

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
    </div>
  );
});
