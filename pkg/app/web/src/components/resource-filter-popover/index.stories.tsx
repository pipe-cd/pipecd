import { action } from "@storybook/addon-actions";
import * as React from "react";
import { ResourceFilterPopover } from "./";

export default {
  title: "ResourceFilterPopover",
  component: ResourceFilterPopover,
};

export const overview: React.FC = () => (
  <ResourceFilterPopover
    enables={{ Pod: true, ReplicaSet: false }}
    onChange={action("onChange")}
  />
);
