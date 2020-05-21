import { ThemeProvider } from "@material-ui/core";
import CssBaseline from "@material-ui/core/CssBaseline";
import React from "react";
import { render } from "react-dom";
import { Header } from "./components/header";
import { theme } from "./theme";
import { Provider } from "react-redux";
import { store } from "./store";

render(
  <Provider store={store}>
    <ThemeProvider theme={theme}>
      <>
        <CssBaseline />
        <Header />
      </>
    </ThemeProvider>
  </Provider>,
  document.getElementById("root")
);
