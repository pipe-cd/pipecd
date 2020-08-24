import React from "react";
import { LabeledText } from "./labeled-text";

export default {
  title: "COMMON/LabeledText",
  component: LabeledText,
};

export const overview: React.FC = () => (
  <LabeledText label="piped" value="hello-world" />
);
