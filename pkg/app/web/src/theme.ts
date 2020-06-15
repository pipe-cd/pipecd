import { createMuiTheme } from "@material-ui/core/styles";
import cyan from "@material-ui/core/colors/cyan";

export const theme = createMuiTheme({
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
