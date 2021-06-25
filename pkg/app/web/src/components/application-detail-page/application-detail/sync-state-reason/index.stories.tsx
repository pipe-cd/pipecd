import { Story } from "@storybook/react";
import { SyncStateReason, SyncStateReasonProps } from ".";

export default {
  title: "APPLICATION/SyncStateReason",
  component: SyncStateReason,
};

const Template: Story<SyncStateReasonProps> = (args) => (
  <SyncStateReason {...args} />
);

export const Overview = Template.bind({});
Overview.args = {
  summary: "There are 1 missing manifests and 2 redundant manifests.",
  detail: `The following 1 manifests are defined in Git, but NOT appearing in the cluster:
  - apiVersion=v1, kind=Service, namespace=default, name=wait-approvalThe following 2 manifests are NOT defined in Git, but appearing in the cluster:
  - apiVersion=apps/v1, kind=Deployment, namespace=default, name=wait-approval-canary- apiVersion=v1, kind=Service, namespace=default, name=wait-approval`,
};

export const Diff = Template.bind({});
Diff.args = {
  summary: "Summary message",
  detail: `message\n\n+ added-line\n- deleted-line`,
};
