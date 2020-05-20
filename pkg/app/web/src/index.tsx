import { ThemeProvider } from "@material-ui/core";
import CssBaseline from "@material-ui/core/CssBaseline";
import React from "react";
import { render } from "react-dom";
import { Header } from "./components/header";
import { theme } from "./theme";

render(
  <ThemeProvider theme={theme}>
    <>
      <CssBaseline />
      <Header />
    </>
  </ThemeProvider>,
  document.getElementById("root")
);
