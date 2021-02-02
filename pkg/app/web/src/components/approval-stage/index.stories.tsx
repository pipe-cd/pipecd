import { action } from "@storybook/addon-actions";
import React from "react";
import { ApprovalStage } from "./";

export default {
  title: "DEPLOYMENT/Pipeline/ApprovalStage",
  component: ApprovalStage,
};

export const overview: React.FC = () => (
  <ApprovalStage
    id="stage-1"
    name="K8S_CANARY_ROLLOUT"
    onClick={action("onClick")}
    active={false}
  />
);
