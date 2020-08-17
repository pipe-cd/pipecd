import { createMuiTheme } from "@material-ui/core/styles";
import cyan from "@material-ui/core/colors/cyan";

export const theme = createMuiTheme({
  props: {
    MuiButtonBase: {
      disableRipple: true,
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
    primary: { main: "#2C387E" },
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
  },
});
