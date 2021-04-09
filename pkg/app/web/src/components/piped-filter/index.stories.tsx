import { action } from "@storybook/addon-actions";
import * as React from "react";
import { PipedFilter } from "./";

export default {
  title: "SETTINGS/Piped/PipedFilter",
  component: PipedFilter,
};

export const overview: React.FC = () => (
  <PipedFilter values={{ enabled: false }} onChange={action("onChange")} />
);
