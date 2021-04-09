import * as React from "react";
import { DiffView } from "./";

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
