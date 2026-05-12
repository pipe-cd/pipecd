import { createTheme, Theme } from "@mui/material/styles";
import { cyan } from "@mui/material/colors";

declare module "@mui/material/styles/createTypography" {
  interface FontStyle {
    fontFamilyMono: string;
  }

  interface FontStyleOptions {
    fontFamilyMono: string;
  }
}

const fontFamilyMono =
  '"Roboto Mono",SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace';

const commonComponents = {
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
};

export const lightTheme = createTheme({
  palette: {
    mode: "light",
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
      default: "#fafafa",
      paper: "#ffffff",
    },
    text: {
      primary: "rgba(0, 0, 0, 0.87)",
      secondary: "rgba(0, 0, 0, 0.6)",
    },
  },
  components: commonComponents,
  typography: {
    subtitle2: {
      fontWeight: 600,
    },
    fontFamilyMono,
  },
});

export const darkTheme = createTheme({
  palette: {
    mode: "dark",
    primary: { main: "#6fa3ff" },
    success: {
      main: "#66bb6a",
      light: "#81c784",
      dark: "#2e7d32",
    },
    error: {
      main: "#ef5350",
      light: "#e57373",
      dark: "#c62828",
    },
    secondary: cyan,
    background: {
      default: "#121212",
      paper: "#1e1e1e",
    },
    text: {
      primary: "#ffffff",
      secondary: "rgba(255, 255, 255, 0.7)",
    },
  },
  components: commonComponents,
  typography: {
    subtitle2: {
      fontWeight: 600,
    },
    fontFamilyMono,
  },
});

export const getTheme = (mode: "light" | "dark"): Theme => {
  return mode === "dark" ? darkTheme : lightTheme;
};

// Default export for backward compatibility
export const theme = lightTheme;
