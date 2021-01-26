import React from "react";
import { ProjectSettingLabeledText } from "./project-setting-labeled-text";

export default {
  title: "ProjectSettingLabeledText",
  component: ProjectSettingLabeledText,
};

export const overview: React.FC = () => (
  <ProjectSettingLabeledText label="Label" value="value" />
);
