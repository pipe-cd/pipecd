import * as React from "react";
import { SplitButton } from "./";
import CancelIcon from "@material-ui/icons/Cancel";
import { action } from "@storybook/addon-actions";

export default {
  title: "COMMON/SplitButton",
  component: SplitButton,
};

export const overview: React.FC = () => (
  <SplitButton
    label="select option"
    startIcon={<CancelIcon />}
    options={["Cancel", "Cancel Without Rollback"]}
    onClick={action("onClick")}
    disabled={false}
    loading={false}
  />
);

export const loading: React.FC = () => (
  <SplitButton
    label="select option"
    startIcon={<CancelIcon />}
    options={["Cancel", "Cancel Without Rollback"]}
    onClick={action("onClick")}
    loading
    disabled
  />
);
