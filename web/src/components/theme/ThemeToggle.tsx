import React from "react";
import { IconButton, Tooltip } from "@mui/material";
import { DarkMode, LightMode } from "@mui/icons-material";
import { useTheme } from "../../contexts/ThemeContext";

/**
 * ThemeToggle component for switching between light and dark modes
 * Can be placed in the header or navigation
 */
export const ThemeToggle: React.FC = () => {
  const { mode, toggleTheme } = useTheme();

  return (
    <Tooltip title={`Switch to ${mode === "light" ? "dark" : "light"} mode`}>
      <IconButton onClick={toggleTheme} size="small" color="inherit">
        {mode === "light" ? <DarkMode /> : <LightMode />}
      </IconButton>
    </Tooltip>
  );
};

export default ThemeToggle;