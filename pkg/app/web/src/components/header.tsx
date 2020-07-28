import React, { FC, memo } from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  makeStyles,
  Avatar,
  Link,
  Button,
} from "@material-ui/core";
import {
  APP_NAME,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_LOGIN,
} from "../constants";
import { Link as RouterLink } from "react-router-dom";
import { useMe } from "../modules/me";
import ArrowDownIcon from "@material-ui/icons/ArrowDropDown";

const useStyles = makeStyles((theme) => ({
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
  },
  left: {
    flexGrow: 1,
    display: "flex",
    alignItems: "center",
  },
  appIcon: {
    marginRight: theme.spacing(2),
    width: theme.spacing(4),
    height: theme.spacing(4),
  },
  link: {
    marginRight: theme.spacing(2),
  },
  userAvatar: {
    width: theme.spacing(4),
    height: theme.spacing(4),
  },
  projectName: {
    marginLeft: theme.spacing(1),
    textTransform: "none",
  },
}));

export const Header: FC = memo(function Header() {
  const classes = useStyles();
  const me = useMe();

  return (
    <AppBar position="static" className={classes.appBar}>
      <Toolbar variant="dense">
        <div className={classes.left}>
          <Avatar className={classes.appIcon}>P</Avatar>
          <Typography variant="h6">{APP_NAME}</Typography>
          {me?.isLogin && (
            <Button
              color="inherit"
              className={classes.projectName}
              endIcon={<ArrowDownIcon />}
            >
              {me.projectId}
            </Button>
          )}
        </div>
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
        {me?.isLogin ? (
          <Avatar className={classes.userAvatar} src={me.avatarUrl} />
        ) : (
          <Link color="inherit" component={RouterLink} to={PAGE_PATH_LOGIN}>
            <Typography variant="body2">Login</Typography>
          </Link>
        )}
      </Toolbar>
    </AppBar>
  );
});
