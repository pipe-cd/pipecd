import {
  Button,
  Card,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import ArrowRightAltIcon from "@material-ui/icons/ArrowRightAlt";
import MuiAlert from "@material-ui/lab/Alert";
import { FC, memo, useState } from "react";
import { useCookies } from "react-cookie";
import { Navigate } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS, PAGE_PATH_LOGIN } from "~/constants/path";
import { getQueryStringValue } from "~/hooks/use-query-string";
import { useAppSelector } from "~/hooks/redux";
import { LoginForm } from "./login-form";
import { LOGGING_IN_PROJECT } from "~/constants/localstorage";

const CONTENT_WIDTH = 500;

const useStyles = makeStyles((theme) => ({
  root: {
    padding: theme.spacing(2),
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    flexDirection: "column",
    flex: 1,
  },
  content: {
    display: "flex",
    flexDirection: "column",
    padding: theme.spacing(3),
    width: CONTENT_WIDTH,
    textAlign: "center",
  },
  fields: {
    display: "flex",
    flexDirection: "column",
    marginTop: theme.spacing(4),
  },
  note: {
    color: "orange",
    textAlign: "right",
  },
  buttons: {
    display: "flex",
    justifyContent: "flex-end",
    marginTop: theme.spacing(2),
  },
  loginError: {
    width: CONTENT_WIDTH,
    marginBottom: theme.spacing(2),
  },
}));

export const LoginPage: FC = memo(function LoginPage() {
  const classes = useStyles();
  const me = useAppSelector((state) => state.me);
  const [name, setName] = useState<string>("");
  const [cookies, , removeCookie] = useCookies(["error"]);
  const queryProject = getQueryStringValue("project") as string;
  const project = queryProject
    ? queryProject
    : localStorage.getItem(LOGGING_IN_PROJECT) || "";

  const handleCloseErrorAlert = (): void => {
    removeCookie("error");
  };

  const handleOnContinue = (): void => {
    window.location.href = `${PAGE_PATH_LOGIN}?project=${name}`;
  };

  const isPlayEnvironment = window.location.hostname.includes("play.");

  return (
    <div className={classes.root}>
      {me && me.isLogin && <Navigate to={PAGE_PATH_APPLICATIONS} replace />}
      {cookies.error && (
        <MuiAlert
          severity="error"
          className={classes.loginError}
          onClose={handleCloseErrorAlert}
        >
          {cookies.error}
        </MuiAlert>
      )}
      <Card className={classes.content}>
        {project ? (
          <LoginForm projectName={project} />
        ) : (
          <div>
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
              {isPlayEnvironment && (
                <div className={classes.note}>
                  Input <strong>play</strong> if you want to join the playground
                  environment
                </div>
              )}
            </div>
            <div className={classes.buttons}>
              <Button
                type="submit"
                color="primary"
                variant="contained"
                endIcon={<ArrowRightAltIcon />}
                disabled={name === ""}
                onClick={handleOnContinue}
              >
                CONTINUE
              </Button>
            </div>
          </div>
        )}
      </Card>
    </div>
  );
});
