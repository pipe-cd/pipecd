import {
  Box,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
} from "@mui/material";
import { FC, memo } from "react";
import { NavLink, Outlet, useLocation } from "react-router-dom";
import {
  PAGE_PATH_SETTINGS_API_KEY,
  PAGE_PATH_SETTINGS_PIPED,
  PAGE_PATH_SETTINGS_PROJECT,
} from "~/constants/path";

const drawerWidth = 240;

const MENU_ITEMS = [
  ["Piped", PAGE_PATH_SETTINGS_PIPED],
  ["Project", PAGE_PATH_SETTINGS_PROJECT],
  ["API Key", PAGE_PATH_SETTINGS_API_KEY],
];

export const SettingsIndexPage: FC = memo(function SettingsIndexPage() {
  const location = useLocation();
  return (
    <Box sx={{ flex: 1, display: "flex", overflow: "hidden" }}>
      <Drawer
        sx={{
          width: drawerWidth,
          flexShrink: 0,
        }}
        variant="permanent"
        slotProps={{
          paper: {
            sx: {
              top: "auto",
              width: drawerWidth,
            },
          },
        }}
      >
        <Box sx={{ overflow: "auto" }}>
          <List sx={{ paddingTop: 0 }}>
            {MENU_ITEMS.map(([text, link]) => (
              <ListItem
                key={`menu-item-${text}`}
                component={NavLink}
                to={link}
                disablePadding
              >
                <ListItemButton
                  selected={link === location.pathname}
                  sx={{
                    color: "text.primary",
                    "&.Mui-selected": {
                      backgroundColor: (theme) => theme.palette.action.selected,
                    },
                  }}
                >
                  <ListItemText primary={text} />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </Box>
      </Drawer>
      <Box
        component={"main"}
        sx={{ display: "flex", flexDirection: "column", flexGrow: 1 }}
      >
        <Outlet />
      </Box>
    </Box>
  );
});
