import { Box, Button, Card, TextField, Typography } from "@mui/material";
import ArrowRightAltIcon from "@mui/icons-material/ArrowRightAlt";
import MuiAlert from "@mui/material/Alert";
import { FC, memo, useState } from "react";
import { useCookies } from "react-cookie";
import { Navigate } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS, PAGE_PATH_LOGIN } from "~/constants/path";
import { getQueryStringValue } from "~/hooks/use-query-string";
import { useAppSelector } from "~/hooks/redux";
import { LoginForm } from "./login-form";
import { LOGGING_IN_PROJECT } from "~/constants/localstorage";

const CONTENT_WIDTH = 500;

export const LoginPage: FC = memo(function LoginPage() {
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
    <Box
      sx={{
        padding: 2,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        flexDirection: "column",
        flex: 1,
      }}
    >
      {me && me.isLogin && <Navigate to={PAGE_PATH_APPLICATIONS} replace />}
      {cookies.error && (
        <MuiAlert
          severity="error"
          sx={{
            width: CONTENT_WIDTH,
            marginBottom: 2,
          }}
          onClose={handleCloseErrorAlert}
        >
          {cookies.error}
        </MuiAlert>
      )}
      <Card
        sx={{
          display: "flex",
          flexDirection: "column",
          padding: 3,
          width: CONTENT_WIDTH,
          textAlign: "center",
        }}
      >
        {project ? (
          <LoginForm projectName={project} />
        ) : (
          <div>
            <Typography variant="h4">Sign in to your project</Typography>
            <Box
              sx={{
                display: "flex",
                flexDirection: "column",
                marginTop: 4,
              }}
            >
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
                <Box sx={{ color: "orange", textAlign: "right" }}>
                  Input <strong>play</strong> if you want to join the playground
                  environment
                </Box>
              )}
            </Box>
            <Box
              sx={{
                display: "flex",
                justifyContent: "flex-end",
                marginTop: 2,
              }}
            >
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
            </Box>
          </div>
        )}
      </Card>
    </Box>
  );
});
