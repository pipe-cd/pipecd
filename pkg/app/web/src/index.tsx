import { ThemeProvider } from "@material-ui/core";
import CssBaseline from "@material-ui/core/CssBaseline";
import React from "react";
import { render } from "react-dom";
import { theme } from "./theme";
import { Provider } from "react-redux";
import { store } from "./store";
import { Pages as App } from "./pages/index";
import { Router } from "react-router-dom";
import { history } from "./history";

render(
  <Provider store={store}>
    <ThemeProvider theme={theme}>
      <Router history={history}>
        <CssBaseline />
        <App />
      </Router>
    </ThemeProvider>
  </Provider>,
  document.getElementById("root")
);
