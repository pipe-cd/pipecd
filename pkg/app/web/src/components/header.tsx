import React from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  makeStyles,
  Avatar,
  Link,
} from "@material-ui/core";
import {
  APP_NAME,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_SETTINGS,
} from "../constants";
import { Link as RouterLink } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  title: {
    flexGrow: 1,
  },
  appIcon: {
    marginRight: theme.spacing(2),
    width: theme.spacing(4),
    height: theme.spacing(4),
  },
  link: {
    marginRight: theme.spacing(2),
  },
}));

export const Header: React.FC = () => {
  const classes = useStyles();
  return (
    <AppBar position="static" className={classes.appBar}>
      <Toolbar variant="dense">
        <Avatar className={classes.appIcon}>P</Avatar>
        <Typography variant="h6" className={classes.title}>
          {APP_NAME}
        </Typography>
        <Link
          component={RouterLink}
          className={classes.link}
          color="inherit"
          to={PAGE_PATH_APPLICATIONS}
        >
          Applications
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          color="inherit"
          to={PAGE_PATH_DEPLOYMENTS}
        >
          Deployments
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          color="inherit"
          to={PAGE_PATH_INSIGHTS}
        >
          Insights
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          color="inherit"
          to={PAGE_PATH_SETTINGS}
        >
          Settings
        </Link>
        <Button color="inherit">Login</Button>
      </Toolbar>
    </AppBar>
  );
};
