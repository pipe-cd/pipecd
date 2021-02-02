import React from "react";
import { ProjectSettingLabeledText } from "./";

export default {
  title: "ProjectSettingLabeledText",
  component: ProjectSettingLabeledText,
};

export const overview: React.FC = () => (
  <ProjectSettingLabeledText label="Label" value="value" />
);
