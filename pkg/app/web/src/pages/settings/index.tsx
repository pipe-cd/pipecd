import {
  Drawer,
  List,
  ListItem,
  ListItemText,
  makeStyles,
  Toolbar,
} from "@material-ui/core";
import React, { FC, memo } from "react";
import { Link as RouterLink, Route, Switch } from "react-router-dom";
import { PAGE_PATH_SETTINGS_PIPED } from "../../constants";
import { SettingsPipedPage } from "./piped";

const drawerWidth = 240;

const useStyles = makeStyles(() => ({
  root: {
    flex: 1,
    display: "flex",
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
    flexGrow: 1,
  },
}));

const MENU_ITEMS = [["Piped", PAGE_PATH_SETTINGS_PIPED]];

export const SettingsIndexPage: FC = memo(() => {
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
                component={RouterLink}
                to={link}
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
            path={PAGE_PATH_SETTINGS_PIPED}
            component={SettingsPipedPage}
          />
        </Switch>
      </main>
    </div>
  );
});
