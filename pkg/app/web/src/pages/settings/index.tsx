import {
  Drawer,
  List,
  ListItem,
  ListItemText,
  makeStyles,
  Toolbar,
} from "@material-ui/core";
import React, { FC, memo } from "react";
import { NavLink, Redirect, Route, Switch } from "react-router-dom";
import {
  PAGE_PATH_SETTINGS,
  PAGE_PATH_SETTINGS_PIPED,
  PAGE_PATH_SETTINGS_ENV,
  PAGE_PATH_SETTINGS_PROJECT,
} from "../../constants/path";
import { SettingsPipedPage } from "./piped";
import { SettingsEnvironmentPage } from "./environment";
import { SettingsProjectPage } from "./project";

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    flex: 1,
    display: "flex",
    overflow: "hidden",
  },
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerContainer: {
    overflow: "auto",
  },
  content: {
    display: "flex",
    flexDirection: "column",
    flexGrow: 1,
  },
  activeNav: {
    backgroundColor: theme.palette.action.selected,
  },
}));

const MENU_ITEMS = [
  ["Piped", PAGE_PATH_SETTINGS_PIPED],
  ["Environment", PAGE_PATH_SETTINGS_ENV],
  ["Project", PAGE_PATH_SETTINGS_PROJECT],
];

export const SettingsIndexPage: FC = memo(function SettingsIndexPage() {
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{ paper: classes.drawerPaper }}
      >
        <Toolbar variant="dense" />
        <div className={classes.drawerContainer}>
          <List>
            {MENU_ITEMS.map(([text, link]) => (
              <ListItem
                key={`menu-item-${text}`}
                button
                component={NavLink}
                to={link}
                activeClassName={classes.activeNav}
              >
                <ListItemText primary={text} />
              </ListItem>
            ))}
          </List>
        </div>
      </Drawer>
      <main className={classes.content}>
        <Switch>
          <Route
            exact
            path={PAGE_PATH_SETTINGS}
            component={() => <Redirect to={PAGE_PATH_SETTINGS_PIPED} />}
          />
          <Route
            exact
            path={PAGE_PATH_SETTINGS_PIPED}
            component={SettingsPipedPage}
          />
          <Route
            exact
            path={PAGE_PATH_SETTINGS_ENV}
            component={SettingsEnvironmentPage}
          />
          <Route
            exact
            path={PAGE_PATH_SETTINGS_PROJECT}
            component={SettingsProjectPage}
          />
        </Switch>
      </main>
    </div>
  );
});
