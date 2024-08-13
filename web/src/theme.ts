import { createTheme } from "@material-ui/core/styles";
import cyan from "@material-ui/core/colors/cyan";

declare module "@material-ui/core/styles/createTypography" {
  interface FontStyle {
    fontFamilyMono: string;
  }

  interface FontStyleOptions {
    fontFamilyMono: string;
  }
}

export const theme = createTheme({
  props: {
    MuiButtonBase: {
      disableRipple: true,
    },
    MuiTypography: {
      variantMapping: {
        body1: "div",
        body2: "div",
      },
    },
  },
  overrides: {
    MuiCssBaseline: {
      "@global": {
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
      paper: {
        borderRadius: 16,
      },
    },
    MuiDialogActions: {
      spacing: {
        padding: 16,
      },
    },
  },
  palette: {
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
  },
  typography: {
    subtitle2: {
      fontWeight: 600,
    },
    fontFamilyMono:
      '"Roboto Mono",SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono","Courier New",monospace',
  },
});
