import React from "react";
import { LogViewer } from "./log-viewer";

export default {
  title: "LogViewer",
  component: LogViewer
};

export const overview: React.FC = () => (
  <LogViewer />
);