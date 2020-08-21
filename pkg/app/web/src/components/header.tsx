import React, { FC, memo, useState } from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  makeStyles,
  Avatar,
  Link,
  Button,
  IconButton,
  MenuItem,
  Menu,
} from "@material-ui/core";
import {
  APP_NAME,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_LOGIN,
  LOGOUT_ENDPOINT,
} from "../constants";
import { NavLink as RouterLink } from "react-router-dom";
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
  userAvatar: {
    width: theme.spacing(4),
    height: theme.spacing(4),
  },
  projectName: {
    marginLeft: theme.spacing(1),
    textTransform: "none",
  },
  link: {
    marginRight: theme.spacing(2),
    display: "inline-flex",
    height: "100%",
    alignItems: "center",
  },
  activeLink: {
    fontWeight: "bold",
    borderBottom: `2px solid ${theme.palette.background.paper}`,
  },
}));

export const Header: FC = memo(function Header() {
  const classes = useStyles();
  const me = useMe();
  const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);

  const handleClose = (): void => {
    setAnchorEl(null);
  };

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
          activeClassName={classes.activeLink}
          color="inherit"
          to={PAGE_PATH_APPLICATIONS}
        >
          Applications
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          activeClassName={classes.activeLink}
          color="inherit"
          to={PAGE_PATH_DEPLOYMENTS}
        >
          Deployments
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          activeClassName={classes.activeLink}
          color="inherit"
          to={PAGE_PATH_INSIGHTS}
        >
          Insights
        </Link>
        <Link
          component={RouterLink}
          className={classes.link}
          activeClassName={classes.activeLink}
          color="inherit"
          to={PAGE_PATH_SETTINGS}
        >
          Settings
        </Link>
        {me?.isLogin ? (
          <IconButton
            aria-controls="user-menu"
            aria-haspopup="true"
            onClick={(e) => setAnchorEl(e.currentTarget)}
          >
            <Avatar className={classes.userAvatar} src={me.avatarUrl} />
          </IconButton>
        ) : (
          <Link color="inherit" component={RouterLink} to={PAGE_PATH_LOGIN}>
            <Typography variant="body2">Login</Typography>
          </Link>
        )}
      </Toolbar>

      <Menu
        id="user-menu"
        anchorEl={anchorEl}
        keepMounted
        open={Boolean(anchorEl)}
        onClose={handleClose}
      >
        <MenuItem component={Link} href={LOGOUT_ENDPOINT}>
          Logout
        </MenuItem>
      </Menu>
    </AppBar>
  );
});
