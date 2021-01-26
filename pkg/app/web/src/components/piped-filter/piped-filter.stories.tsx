import { action } from "@storybook/addon-actions";
import React from "react";
import { PipedFilter } from "./piped-filter";

export default {
  title: "SETTINGS/Piped/PipedFilter",
  component: PipedFilter,
};

export const overview: React.FC = () => (
  <PipedFilter values={{ enabled: false }} onChange={action("onChange")} />
);
