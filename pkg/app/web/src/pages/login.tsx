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
import { LoginForm } from "../components/login-form";
import { setProjectName, useProjectName } from "../modules/login";
import { useMe } from "../modules/me";
import { PAGE_PATH_APPLICATIONS } from "../constants";
import { Redirect } from "react-router-dom";

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

  const handleOnContinue = (): void => {
    dispatch(setProjectName(name));
  };

  return (
    <div className={classes.root}>
      {me && me.isLogin && <Redirect to={PAGE_PATH_APPLICATIONS} />}
      <Card className={classes.content}>
        {projectName === null ? (
          <>
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
                color="primary"
                variant="contained"
                endIcon={<ArrowRightAltIcon />}
                onClick={handleOnContinue}
                disabled={name === ""}
              >
                Continue
              </Button>
            </div>
          </>
        ) : (
          <LoginForm />
        )}
      </Card>
    </div>
  );
});
