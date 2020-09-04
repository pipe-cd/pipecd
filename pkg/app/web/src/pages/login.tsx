import {
  Button,
  Card,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import ArrowRightAltIcon from "@material-ui/icons/ArrowRightAlt";
import React, { FC, memo, useState } from "react";
import { useDispatch } from "react-redux";
import { Redirect } from "react-router-dom";
import { LoginForm } from "../components/login-form";
import { PAGE_PATH_APPLICATIONS } from "../constants";
import { setProjectName, useProjectName } from "../modules/login";
import { useMe } from "../modules/me";

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    flex: 1,
  },
  content: {
    display: "flex",
    flexDirection: "column",
    padding: theme.spacing(3),
    width: 500,
    textAlign: "center",
  },
  fields: {
    display: "flex",
    flexDirection: "column",
    marginTop: theme.spacing(4),
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
    marginTop: theme.spacing(4),
  },
}));

export const LoginPage: FC = memo(function LoginPage() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const me = useMe();
  const projectName = useProjectName();
  const [name, setName] = useState<string>("");

  const handleOnContinue = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    dispatch(setProjectName(name));
  };

  return (
    <div className={classes.root}>
      {me && me.isLogin && <Redirect to={PAGE_PATH_APPLICATIONS} />}
      <Card className={classes.content}>
        {projectName === null ? (
          <form onSubmit={handleOnContinue}>
            <Typography variant="h4">Sign in to your project</Typography>
            <div className={classes.fields}>
              <TextField
                id="project-name"
                name="project-name"
                label="Project Name"
                variant="outlined"
                margin="dense"
                required
                value={name}
                onChange={(e) => setName(e.currentTarget.value)}
              />
            </div>
            <div className={classes.buttons}>
              <Button
                type="submit"
                color="primary"
                variant="contained"
                endIcon={<ArrowRightAltIcon />}
                disabled={name === ""}
              >
                CONTINUE
              </Button>
            </div>
          </form>
        ) : (
          <LoginForm />
        )}
      </Card>
    </div>
  );
});
