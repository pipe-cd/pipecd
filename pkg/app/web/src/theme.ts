import { createMuiTheme } from "@material-ui/core/styles";
import blue from "@material-ui/core/colors/blue";
import cyan from "@material-ui/core/colors/cyan";

export const theme = createMuiTheme({
  props: {
    MuiTypography: {
      variantMapping: {
        subtitle2: "span"
      }
    }
  },
  palette: {
    primary: blue,
    secondary: cyan
  },
  typography: {
    subtitle2: {
      fontWeight: 600
    }
  }
});
