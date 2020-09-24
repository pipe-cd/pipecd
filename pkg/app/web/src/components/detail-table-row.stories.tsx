import React from "react";
import { DetailTableRow } from "./detail-table-row";

export default {
  title: "COMMON/DetailTableRow",
  component: DetailTableRow,
};

export const overview: React.FC = () => (
  <DetailTableRow label="piped" value="hello-world" />
);
