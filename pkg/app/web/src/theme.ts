import { createMuiTheme } from "@material-ui/core/styles";
import cyan from "@material-ui/core/colors/cyan";

export const theme = createMuiTheme({
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
        },
      },
    },
  },
  props: {
    MuiTypography: {
      variantMapping: {
        subtitle2: "span",
      },
    },
  },
  palette: {
    primary: { main: "#1a73e8" },
    secondary: cyan,
  },
  typography: {
    subtitle2: {
      fontWeight: 600,
    },
  },
});
