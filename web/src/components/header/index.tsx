import {
  AppBar,
  Avatar,
  Button,
  IconButton,
  Link,
  Menu,
  MenuItem,
  Toolbar,
  Typography,
  makeStyles,
} from "@material-ui/core";
import { ExitToApp, MoreVert, OpenInNew } from "@material-ui/icons";
import ArrowDownIcon from "@material-ui/icons/ArrowDropDown";
import clsx from "clsx";
import { FC, memo, useEffect, useState } from "react";
import { NavLink as RouterLink, useLocation } from "react-router-dom";
import { APP_NAME } from "~/constants/common";
import { LOGGING_IN_PROJECT, USER_PROJECTS } from "~/constants/localstorage";
import {
  LOGOUT_ENDPOINT,
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_DEPLOYMENT_CHAINS,
  PAGE_PATH_DEPLOYMENT_TRACE,
  PAGE_PATH_EVENTS,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_LOGIN,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_TOP,
} from "~/constants/path";
import { useAppSelector } from "~/hooks/redux";
import logo from "~~/assets/logo.svg";

export const APP_HEADER_HEIGHT = 56;

const useStyles = makeStyles((theme) => ({
  root: {
    zIndex: theme.zIndex.drawer + 1,
    height: APP_HEADER_HEIGHT,
  },
  logo: {
    height: APP_HEADER_HEIGHT,
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
  right: {
    height: "100%",
    overflow: "hidden",
    display: "flex",
    alignItems: "center",
    "&:hover": {
      color: theme.palette.grey[400],
    },
  },
  link: {
    marginRight: theme.spacing(2),
    display: "inline-flex",
    height: "100%",
    alignItems: "center",
    "&:hover": {
      color: theme.palette.grey[100],
      textDecoration: "none",
    },
  },
  activeLink: {
    borderBottom: `4px solid ${theme.palette.background.paper}`,
  },
  iconButton: {
    padding: 0,
  },
  iconOpenInNew: {
    fontSize: "0.95rem",
    marginLeft: "5px",
    marginBottom: "-2px",
    color: "rgba(0, 0, 0, 0.5)",
  },
  menuItem: {
    color: "#283778",
    textDecorationLine: "none",
  },
}));

export const Header: FC = memo(function Header() {
  const classes = useStyles();
  const location = useLocation();
  const me = useAppSelector((state) => state.me);
  const [userAnchorEl, setUserAnchorEl] = useState<HTMLButtonElement | null>(
    null
  );
  const [moreAnchorEl, setMoreAnchorEl] = useState<HTMLButtonElement | null>(
    null
  );
  const [
    projectAnchorEl,
    setProjectAnchorEl,
  ] = useState<HTMLButtonElement | null>(null);

  const [selectableProjects, setSelectableProjects] = useState<string[]>([]);

  useEffect(() => {
    if (me?.isLogin) {
      const projects =
        localStorage
          .getItem(USER_PROJECTS)
          ?.split(",")
          .filter((e) => e !== me.projectId) || [];
      setSelectableProjects(projects);
    }
  }, [me]);

  const handleSwitchProject = (proj: string): void => {
    localStorage.setItem(LOGGING_IN_PROJECT, proj);
    window.location.href = LOGOUT_ENDPOINT;
  };

  return (
    <AppBar position="static" className={classes.root}>
      <Toolbar variant="dense">
        <div className={classes.left}>
          <Link
            component={RouterLink}
            to={PAGE_PATH_TOP}
            className={classes.logo}
          >
            <img className={classes.logo} src={logo} alt={APP_NAME}></img>
          </Link>
          {me?.isLogin && (
            <Button
              color="inherit"
              className={classes.projectName}
              endIcon={<ArrowDownIcon />}
              onClick={(e) => setProjectAnchorEl(e.currentTarget)}
            >
              {me.projectId}
            </Button>
          )}
        </div>
        <div className={classes.right}>
          {me?.isLogin ? (
            <>
              <Link
                component={RouterLink}
                className={clsx(classes.link, {
                  [classes.activeLink]:
                    location.pathname === PAGE_PATH_APPLICATIONS,
                })}
                color="inherit"
                to={PAGE_PATH_APPLICATIONS}
              >
                Applications
              </Link>
              <Link
                component={RouterLink}
                className={clsx(classes.link, {
                  [classes.activeLink]:
                    location.pathname === PAGE_PATH_DEPLOYMENTS,
                })}
                color="inherit"
                to={PAGE_PATH_DEPLOYMENTS}
              >
                Deployments
              </Link>
              <Link
                component={RouterLink}
                className={clsx(classes.link, {
                  [classes.activeLink]:
                    location.pathname === PAGE_PATH_DEPLOYMENT_TRACE,
                })}
                color="inherit"
                to={PAGE_PATH_DEPLOYMENT_TRACE}
              >
                Traces
              </Link>
              <Link
                component={RouterLink}
                className={clsx(classes.link, {
                  [classes.activeLink]:
                    location.pathname === PAGE_PATH_DEPLOYMENT_CHAINS,
                })}
                color="inherit"
                to={PAGE_PATH_DEPLOYMENT_CHAINS}
              >
                Chains
              </Link>
              <IconButton
                color="inherit"
                className={classes.iconButton}
                aria-label="More Menu"
                aria-controls="more-menu"
                aria-haspopup="true"
                size="small"
                onClick={(e) => setMoreAnchorEl(e.currentTarget)}
              >
                <MoreVert />
              </IconButton>
              <Button
                color="inherit"
                aria-label="User Menu"
                aria-controls="user-menu"
                aria-haspopup="true"
                onClick={(e) => setUserAnchorEl(e.currentTarget)}
                style={{
                  display: "flex",
                  alignItems: "center",
                  textTransform: "none",
                }}
              >
                <Avatar
                  className={classes.userAvatar}
                  src={me.avatarUrl}
                  style={{ marginRight: 8 }}
                />
                {me.subject}
              </Button>
            </>
          ) : (
            <Link
              color="inherit"
              component={RouterLink}
              to={PAGE_PATH_LOGIN}
              className={classes.link}
            >
              <Typography variant="body2">Login</Typography>
            </Link>
          )}
        </div>
      </Toolbar>

      <Menu
        id="project-selection"
        anchorEl={projectAnchorEl}
        open={Boolean(projectAnchorEl) && selectableProjects.length !== 0}
        onClose={(): void => {
          setProjectAnchorEl(null);
        }}
      >
        {selectableProjects.map((p) => (
          <MenuItem key={p} onClick={() => handleSwitchProject(p)}>
            {p}
          </MenuItem>
        ))}
      </Menu>

      <Menu
        id="user-menu"
        anchorEl={userAnchorEl}
        open={Boolean(userAnchorEl)}
        getContentAnchorEl={null}
        anchorOrigin={{
          vertical: "bottom",
          horizontal: "right",
        }}
        onClose={(): void => {
          setUserAnchorEl(null);
        }}
      >
        <MenuItem component={Link} href={LOGOUT_ENDPOINT}>
          <ExitToApp style={{ marginRight: 8 }} /> Logout
        </MenuItem>
      </Menu>

      <Menu
        id="more-menu"
        anchorEl={moreAnchorEl}
        open={Boolean(moreAnchorEl)}
        getContentAnchorEl={null}
        anchorOrigin={{ vertical: 35, horizontal: "right" }}
        transformOrigin={{ vertical: "top", horizontal: "right" }}
        onClose={(): void => {
          setMoreAnchorEl(null);
        }}
      >
        <MenuItem
          className={classes.menuItem}
          component={RouterLink}
          to={PAGE_PATH_INSIGHTS}
        >
          Insights
        </MenuItem>
        <MenuItem
          className={classes.menuItem}
          component={RouterLink}
          to={PAGE_PATH_EVENTS}
        >
          Events
        </MenuItem>
        <MenuItem
          className={classes.menuItem}
          component={RouterLink}
          to={PAGE_PATH_SETTINGS}
          divider
        >
          Settings
        </MenuItem>
        <MenuItem
          component={Link}
          href="https://pipecd.dev/docs/"
          target="_blank"
          rel="noreferrer"
        >
          Documentation
          <OpenInNew className={classes.iconOpenInNew} />
        </MenuItem>
        <MenuItem
          component={Link}
          href="https://github.com/pipe-cd/pipecd"
          target="_blank"
          rel="noreferrer"
        >
          GitHub
          <OpenInNew className={classes.iconOpenInNew} />
        </MenuItem>
        <MenuItem disabled={true} dense={true} button={false}>
          {process.env.PIPECD_VERSION}
        </MenuItem>
      </Menu>
    </AppBar>
  );
});
