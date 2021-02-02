import React from "react";
import { DetailTableRow } from "./";

export default {
  title: "COMMON/DetailTableRow",
  component: DetailTableRow,
};

export const overview: React.FC = () => (
  <DetailTableRow label="piped" value="hello-world" />
);
