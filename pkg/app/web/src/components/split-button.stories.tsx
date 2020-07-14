import React from "react";
import { SplitButton } from "./split-button";
import CancelIcon from "@material-ui/icons/Cancel";
import { action } from "@storybook/addon-actions";

export default {
  title: "COMMON|SplitButton",
  component: SplitButton,
};

export const overview: React.FC = () => (
  <SplitButton
    startIcon={<CancelIcon />}
    options={["Cancel", "Cancel Without Rollback"]}
    onClick={action("onClick")}
    loading={false}
  />
);

export const loading: React.FC = () => (
  <SplitButton
    startIcon={<CancelIcon />}
    options={["Cancel", "Cancel Without Rollback"]}
    onClick={action("onClick")}
    loading
  />
);
