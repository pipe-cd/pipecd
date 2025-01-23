import {
  Drawer,
  List,
  ListItem,
  ListItemText,
  makeStyles,
} from "@material-ui/core";
import { FC, memo } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";
import {
  PAGE_PATH_SETTINGS_API_KEY,
  PAGE_PATH_SETTINGS_PIPED,
  PAGE_PATH_SETTINGS_PROJECT,
} from "~/constants/path";

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
    top: "auto",
    width: drawerWidth,
  },
  drawerContainer: {
    overflow: "auto",
  },
  listGroup: {
    paddingTop: 0,
  },
  content: {
    display: "flex",
    flexDirection: "column",
    flexGrow: 1,
  },
  activeNav: {
    backgroundColor: theme.palette.action.selected,
  },
  listItemIcon: {
    minWidth: 110,
  },
}));

const MENU_ITEMS = [
  ["Piped", PAGE_PATH_SETTINGS_PIPED],
  ["Project", PAGE_PATH_SETTINGS_PROJECT],
  ["API Key", PAGE_PATH_SETTINGS_API_KEY],
];

export const SettingsIndexPage: FC = memo(function SettingsIndexPage() {
  const classes = useStyles();
  const location = useLocation();
  return (
    <div className={classes.root}>
      <Drawer
        className={classes.drawer}
        variant="permanent"
        classes={{ paper: classes.drawerPaper }}
      >
        <div className={classes.drawerContainer}>
          <List className={classes.listGroup}>
            {MENU_ITEMS.map(([text, link]) => (
              <ListItem
                key={`menu-item-${text}`}
                button
                component={NavLink}
                to={link}
                className={link === location.pathname ? classes.activeNav : ""}
                selected={link === location.pathname}
              >
                <ListItemText primary={text} />
              </ListItem>
            ))}
          </List>
        </div>
      </Drawer>
      <main className={classes.content}>
        <Outlet />
      </main>
    </div>
  );
});
