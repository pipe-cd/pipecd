import { FC, memo, useEffect, useState } from "react";
import {
  AppBar,
  Toolbar,
  Typography,
  Avatar,
  Link,
  Button,
  IconButton,
  MenuItem,
  Menu,
  Box,
} from "@mui/material";
import { MoreVert } from "@mui/icons-material";
import {
  PAGE_PATH_APPLICATIONS,
  PAGE_PATH_DEPLOYMENTS,
  PAGE_PATH_SETTINGS,
  PAGE_PATH_LOGIN,
  LOGOUT_ENDPOINT,
  PAGE_PATH_TOP,
  PAGE_PATH_INSIGHTS,
  PAGE_PATH_DEPLOYMENT_CHAINS,
  PAGE_PATH_EVENTS,
  PAGE_PATH_DEPLOYMENT_TRACE,
} from "~/constants/path";
import { APP_NAME } from "~/constants/common";
import { LOGGING_IN_PROJECT, USER_PROJECTS } from "~/constants/localstorage";
import { NavLink as RouterLink } from "react-router-dom";
import ArrowDownIcon from "@mui/icons-material/ArrowDropDown";
import logo from "~~/assets/logo.svg";
import { useAppSelector } from "~/hooks/redux";
import NavLink from "./NavLink";
import { IconOpenNewTab, LogoImage } from "./styles";

export const APP_HEADER_HEIGHT = 56;

export const Header: FC = memo(function Header() {
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
    <AppBar
      position="static"
      sx={{
        zIndex: (theme) => theme.zIndex.drawer - 1,
        height: APP_HEADER_HEIGHT,
      }}
    >
      <Toolbar variant="dense">
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            flexGrow: 1,
          }}
        >
          <Link
            component={RouterLink}
            to={PAGE_PATH_TOP}
            sx={{
              height: APP_HEADER_HEIGHT,
            }}
          >
            <LogoImage src={logo} alt={APP_NAME} />
          </Link>
          {me?.isLogin && (
            <Button
              color="inherit"
              sx={{ ml: 2, textTransform: "none" }}
              endIcon={<ArrowDownIcon />}
              onClick={(e) => setProjectAnchorEl(e.currentTarget)}
            >
              {me.projectId}
            </Button>
          )}
        </Box>
        <Box
          sx={{
            height: "100%",
            overflow: "hidden",
            display: "flex",
            alignItems: "center",
          }}
        >
          {me?.isLogin ? (
            <>
              <NavLink href={PAGE_PATH_APPLICATIONS}>Applications</NavLink>
              <NavLink href={PAGE_PATH_DEPLOYMENTS}>Deployments</NavLink>
              <NavLink href={PAGE_PATH_DEPLOYMENT_TRACE}>Traces</NavLink>
              <NavLink href={PAGE_PATH_DEPLOYMENT_CHAINS}>Chains</NavLink>
              <IconButton
                color="inherit"
                aria-label="More Menu"
                aria-controls="more-menu"
                aria-haspopup="true"
                size="small"
                onClick={(e) => setMoreAnchorEl(e.currentTarget)}
                sx={{ p: 0 }}
              >
                <MoreVert />
              </IconButton>
              <IconButton
                aria-label="User Menu"
                aria-controls="user-menu"
                aria-haspopup="true"
                onClick={(e) => setUserAnchorEl(e.currentTarget)}
                size="large"
              >
                <Avatar sx={{ width: 32, height: 32 }} src={me.avatarUrl} />
              </IconButton>
              <Typography variant="body2">{me.subject}</Typography>
            </>
          ) : (
            <NavLink href={PAGE_PATH_LOGIN} active={false}>
              <Typography variant="body2">Login</Typography>
            </NavLink>
          )}
        </Box>
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
        onClose={(): void => {
          setUserAnchorEl(null);
        }}
      >
        <MenuItem component={Link} href={LOGOUT_ENDPOINT}>
          Logout
        </MenuItem>
      </Menu>
      <Menu
        id="more-menu"
        anchorEl={moreAnchorEl}
        open={Boolean(moreAnchorEl)}
        anchorOrigin={{ vertical: 35, horizontal: "right" }}
        transformOrigin={{ vertical: "top", horizontal: "right" }}
        onClose={(): void => {
          setMoreAnchorEl(null);
        }}
      >
        <MenuItem component={RouterLink} to={PAGE_PATH_INSIGHTS}>
          Insights
        </MenuItem>
        <MenuItem component={RouterLink} to={PAGE_PATH_EVENTS}>
          Events
        </MenuItem>
        <MenuItem component={RouterLink} to={PAGE_PATH_SETTINGS} divider>
          Settings
        </MenuItem>
        <MenuItem
          component={Link}
          href="https://pipecd.dev/docs/"
          target="_blank"
          rel="noreferrer"
          sx={{ "&:hover": { textDecorationLine: "underline" } }}
        >
          Documentation
          <IconOpenNewTab />
        </MenuItem>
        <MenuItem
          component={Link}
          href="https://github.com/pipe-cd/pipecd"
          target="_blank"
          rel="noreferrer"
          sx={{ "&:hover": { textDecorationLine: "underline" } }}
        >
          GitHub
          <IconOpenNewTab />
        </MenuItem>
        <MenuItem disabled={true} dense={true}>
          {process.env.PIPECD_VERSION}
        </MenuItem>
      </Menu>
    </AppBar>
  );
});
