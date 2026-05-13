import { createTheme } from "@mui/material/styles";
import { cyan } from "@mui/material/colors";
import { createContext, useContext } from "react";

declare module "@mui/material/styles/createTypography" {
  interface FontStyle {
    fontFamilyMono: string;
  }

  interface FontStyleOptions {
    fontFamilyMono: string;
  }
}

// The two modes the app can be in
export type ThemeMode = "light" | "dark";

// A simple context so any component can read the mode and call toggleTheme
export const ThemeContext = createContext<{
  mode: ThemeMode;
  toggleTheme: () => void;
}>({
  mode: "light",
  toggleTheme: () => {},
});

// Convenience hook — just call useThemeMode() anywhere you need the toggle
export const useThemeMode = () => useContext(ThemeContext);

// Build a MUI theme for whichever mode is requested.
// All existing colors and component tweaks are preserved exactly as before.
export const createAppTheme = (mode: ThemeMode) =>
  createTheme({
    components: {
      MuiButtonBase: {
        defaultProps: {
          disableRipple: true,
        },
      },
      MuiTypography: {
        defaultProps: {
          variantMapping: {
            body1: "div",
            body2: "div",
          },
        },
      },
      MuiCssBaseline: {
        styleOverrides: {
          html: {
            height: "100%",
          },
          body: {
            height: "100%",
          },
          "#root": {
            height: "100%",
            display: "flex",
            flexDirection: "column",
            overflow: "hidden",
          },
        },
      },
      MuiDialog: {
        styleOverrides: {
          paper: {
            borderRadius: 16,
          },
        },
      },
      MuiDialogActions: {
        styleOverrides: {
          spacing: {
            padding: 16,
          },
        },
      },
    },
    palette: {
      mode, // this single line is what makes MUI flip all its colours
      primary: { main: "#283778" },
      success: {
        main: "#539d56",
        light: "#83cf84",
        dark: "#216e2b",
      },
      error: {
        main: "#d6442c",
        light: "#ff7657",
        dark: "#9d0001",
      },
      secondary: cyan,
      background: {
        default: mode === "light" ? "#fafafa" : "#121212",
        paper: mode === "light" ? "#ffffff" : "#1e1e1e",
      },
    },
    typography: {
      subtitle2: {
        fontWeight: 600,
      },
      fontFamilyMono:
        '"Roboto Mono",SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace',
    },
  });

// Keep the old named export so nothing else in the project breaks
export const theme = createAppTheme("light");
