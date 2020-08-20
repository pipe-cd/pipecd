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
import React, { FC, memo, useState } from "react";
import { useStyles as useButtonStyles } from "../styles/button";

const useStyles = makeStyles((theme) => ({
  main: {
    overflow: "auto",
    padding: theme.spacing(3),
    background: theme.palette.background.paper,
  },
  group: {
    padding: theme.spacing(1),
  },
  titleMargin: {
    marginTop: theme.spacing(2),
  },
}));

interface Props {
  username: string | null;
  isUpdatingUsername: boolean;
  isUpdatingPassword: boolean;
  staticAdminDisabled: boolean;
  onUpdateUsername: (username: string) => void;
  onUpdatePassword: (password: string) => Promise<unknown>;
  onToggleAvailability: () => void;
}

export const StaticAdminForm: FC<Props> = memo(function StaticAdminForm({
  username,
  isUpdatingUsername,
  isUpdatingPassword,
  staticAdminDisabled,
  onUpdatePassword,
  onUpdateUsername,
  onToggleAvailability,
}) {
  const classes = useStyles();
  const buttonClasses = useButtonStyles();
  const [usernameState, setUsernameState] = useState(username);
  const [password, setPassword] = useState("");

  return (
    <>
      <Typography variant="h5">Static Admin User</Typography>
      <div className={classes.group}>
        <Typography variant="subtitle1">Status: Enabled</Typography>
        <Button
          color="primary"
          variant="contained"
          onClick={onToggleAvailability}
        >
          {staticAdminDisabled ? "Enable" : "Disable"}
        </Button>

        <Typography variant="h6" className={classes.titleMargin}>
          Change username
        </Typography>

        <Typography variant="body2">Current username: {username}</Typography>

        <FormControl variant="outlined" margin="dense">
          <InputLabel htmlFor="outlined-adornment-username">
            Username
          </InputLabel>
          <OutlinedInput
            id="outlined-adornment-username"
            type="text"
            labelWidth={70}
            value={usernameState || undefined}
            onChange={(e) => setUsernameState(e.target.value)}
            endAdornment={
              <InputAdornment position="end">
                <Button
                  color="primary"
                  disabled={
                    !usernameState ||
                    usernameState === username ||
                    isUpdatingUsername
                  }
                  onClick={() => {
                    if (usernameState) {
                      onUpdateUsername(usernameState);
                    }
                  }}
                >
                  Update
                  {isUpdatingUsername && (
                    <CircularProgress
                      size={24}
                      className={buttonClasses.progress}
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
                  disabled={!password || isUpdatingPassword}
                  onClick={() => {
                    onUpdatePassword(password).then(() => {
                      setPassword("");
                    });
                  }}
                >
                  Update
                  {isUpdatingPassword && (
                    <CircularProgress
                      size={24}
                      className={buttonClasses.progress}
                    />
                  )}
                </Button>
              </InputAdornment>
            }
          />
        </FormControl>
      </div>
    </>
  );
});
