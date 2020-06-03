import React from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  makeStyles,
  Avatar,
  Link
} from "@material-ui/core";
import {
  APP_NAME,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS
} from "../constants";

const useStyles = makeStyles(theme => ({
  title: {
    flexGrow: 1
  },
  appIcon: {
    marginRight: theme.spacing(2),
    width: theme.spacing(4),
    height: theme.spacing(4)
  },
  link: {
    marginRight: theme.spacing(2)
  }
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
        <Link
          className={classes.link}
          color="inherit"
          href={PAGE_PATH_APPLICATIONS}
        >
          Applications
        </Link>
        <Link
          className={classes.link}
          color="inherit"
          href={PAGE_PATH_DEPLOYMENTS}
        >
          Deployments
        </Link>
        <Link className={classes.link} color="inherit">
          Insights
        </Link>
        <Link className={classes.link} color="inherit">
          Settings
        </Link>
        <Button color="inherit">Login</Button>
      </Toolbar>
    </AppBar>
  );
};
