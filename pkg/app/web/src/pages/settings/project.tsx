import {
  Button,
  CircularProgress,
  FormControl,
  InputAdornment,
  InputLabel,
  makeStyles,
  OutlinedInput,
  Typography,
} from "@material-ui/core";
import React, { FC, memo, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../../modules";
import {
  fetchProject,
  ProjectState,
  updatePassword,
  updateUsername,
  toggleAvailability,
} from "../../modules/project";
import { AppDispatch } from "../../store";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
    padding: theme.spacing(3),
    background: theme.palette.background.paper,
  },
  listItem: {
    backgroundColor: theme.palette.background.paper,
  },
  group: {
    padding: theme.spacing(1),
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
  titleMargin: {
    marginTop: theme.spacing(2),
  },
}));

export const SettingsProjectPage: FC = memo(function SettingsProjectPage() {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const project = useSelector<AppState, ProjectState>((state) => state.project);
  const [password, setPassword] = useState("");
  const [username, setUsername] = useState("");

  const handleUpdateUsername = (): void => {
    dispatch(updateUsername({ username })).then(() => {
      dispatch(fetchProject());
    });
  };
  const handleUpdatePassword = (): void => {
    dispatch(updatePassword({ password })).then(() => {
      setPassword("");
    });
  };
  const handleToggleAvailability = (): void => {
    dispatch(toggleAvailability());
  };

  useEffect(() => {
    dispatch(fetchProject());
  }, [dispatch]);

  if (!project) {
    return <CircularProgress />;
  }

  return (
    <div className={classes.main}>
      <Typography variant="h5">Static Admin User</Typography>
      <div className={classes.group}>
        <Typography variant="subtitle1">Status: Enabled</Typography>
        <Button
          color="primary"
          variant="contained"
          onClick={handleToggleAvailability}
        >
          {project.staticAdminDisabled ? "Enable" : "Disable"}
        </Button>

        <Typography variant="h6" className={classes.titleMargin}>
          Change username
        </Typography>

        <Typography variant="body2">
          Current username: {project.username}
        </Typography>

        <FormControl variant="outlined" margin="dense">
          <InputLabel htmlFor="outlined-adornment-username">
            Username
          </InputLabel>
          <OutlinedInput
            id="outlined-adornment-username"
            type="text"
            labelWidth={70}
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            endAdornment={
              <InputAdornment position="end">
                <Button
                  color="primary"
                  disabled={
                    !username ||
                    project.username === username ||
                    project.isUpdatingUsername
                  }
                  onClick={handleUpdateUsername}
                >
                  Update
                  {project.isUpdatingUsername && (
                    <CircularProgress
                      size={24}
                      className={classes.buttonProgress}
                    />
                  )}
                </Button>
              </InputAdornment>
            }
          />
        </FormControl>

        <Typography variant="h6" className={classes.titleMargin}>
          Change password
        </Typography>

        <FormControl variant="outlined" margin="dense">
          <InputLabel htmlFor="outlined-adornment-password">
            Password
          </InputLabel>
          <OutlinedInput
            id="outlined-adornment-password"
            type="password"
            labelWidth={70}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            endAdornment={
              <InputAdornment position="end">
                <Button
                  color="primary"
                  disabled={!password || project.isUpdatingPassword}
                  onClick={handleUpdatePassword}
                >
                  Update
                  {project.isUpdatingPassword && (
                    <CircularProgress
                      size={24}
                      className={classes.buttonProgress}
                    />
                  )}
                </Button>
              </InputAdornment>
            }
          />
        </FormControl>
      </div>
      <Typography variant="h5">SSO</Typography>
      WIP
    </div>
  );
});
