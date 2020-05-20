import React from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  makeStyles,
  Avatar,
} from "@material-ui/core";
import { APP_NAME } from "../constants";

const useStyles = makeStyles((theme) => ({
  title: {
    flexGrow: 1,
  },
  appIcon: {
    marginRight: theme.spacing(2),
    width: theme.spacing(4),
    height: theme.spacing(4),
  },
}));

export const Header: React.FC = () => {
  const classes = useStyles();
  return (
    <AppBar position="static">
      <Toolbar variant="dense">
        <Avatar className={classes.appIcon}>P</Avatar>
        <Typography variant="h6" className={classes.title}>
          {APP_NAME}
        </Typography>
        <Button color="inherit">Login</Button>
      </Toolbar>
    </AppBar>
  );
};
