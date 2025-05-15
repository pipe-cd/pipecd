import { FC, memo } from "react";
import { TextField, Button, Typography, Box } from "@mui/material";
import {
  STATIC_LOGIN_ENDPOINT,
  LOGIN_ENDPOINT,
  PAGE_PATH_LOGIN,
} from "~/constants/path";
import { MarkGithubIcon } from "@primer/octicons-react";
import { LOGGING_IN_PROJECT } from "~/constants/localstorage";

export interface LoginFormProps {
  projectName: string;
}

export const LoginForm: FC<LoginFormProps> = memo(function LoginForm({
  projectName,
}) {
  const handleOnBack = (): void => {
    localStorage.removeItem(LOGGING_IN_PROJECT);
    setTimeout(() => {
      window.location.href = PAGE_PATH_LOGIN;
    }, 300);
  };

  return (
    <Box
      sx={{
        display: "flex",
        alignItems: "center",
        flexDirection: "column",
        flex: 1,
      }}
    >
      <Typography variant="h4">Sign in to {projectName}</Typography>
      <Box sx={{ width: 320 }}>
        <Box
          component="form"
          method="POST"
          action={LOGIN_ENDPOINT}
          sx={{
            display: "flex",
            flexDirection: "column",
            textAlign: "center",
            marginTop: 4,
            width: "100%",
          }}
        >
          <input
            type="hidden"
            id="project-gh"
            name="project"
            value={projectName || undefined}
          />
          <Button
            type="submit"
            color="primary"
            variant="contained"
            sx={{ bgcolor: "#24292E" }}
          >
            <Box mr={1}>
              <MarkGithubIcon />
            </Box>
            LOGIN WITH GITHUB
          </Button>

          <Button
            type="submit"
            color="primary"
            variant="contained"
            sx={{
              background: "#4A90E2",
              marginTop: 1,
            }}
          >
            LOGIN WITH OIDC
          </Button>
        </Box>

        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            marginTop: 3,
            marginBottom: 3,
          }}
        >
          <Box
            sx={{
              flex: 1,
              border: "none",
              borderTop: "1px solid #ddd",
            }}
          />
          <Box sx={{ marginLeft: 2, marginRight: 2 }}>OR</Box>
          <Box
            sx={{
              flex: 1,
              border: "none",
              borderTop: "1px solid #ddd",
            }}
          />
        </Box>

        <Box
          component="form"
          method="POST"
          action={STATIC_LOGIN_ENDPOINT}
          sx={{
            display: "flex",
            flexDirection: "column",
            textAlign: "center",
            marginTop: 4,
            width: "100%",
            gap: 2,
          }}
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
            size="small"
            required
          />
          <TextField
            id="password"
            name="password"
            label="Password"
            type="password"
            variant="outlined"
            size="small"
            required
          />
          <Box
            sx={{
              display: "flex",
              justifyContent: "flex-end",
              marginTop: 3,
            }}
          >
            <Button type="reset" color="primary" onClick={handleOnBack}>
              back
            </Button>
            <Button type="submit" color="primary" variant="contained">
              login
            </Button>
          </Box>
        </Box>
      </Box>
    </Box>
  );
});
