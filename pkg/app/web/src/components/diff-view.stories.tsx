import React from "react";
import { DiffView } from "./diff-view";

export default {
  title: "DiffView",
  component: DiffView,
};

const content = `
+ added line
- deleted line
normal
  indent
`;

export const overview: React.FC = () => <DiffView content={content} />;
