import { styled } from "@mui/material";

export const StateViewRoot = styled("div")(() => ({
  display: "flex",
  flex: 1,
  justifyContent: "center",
  overflow: "hidden",
}));

export const StateViewWrapper = styled("div")(() => ({
  flex: 1,
  display: "flex",
  justifyContent: "center",
  overflow: "hidden",
}));

export const StateView = styled("div")(() => ({
  position: "relative",
  overflow: "auto",
}));
