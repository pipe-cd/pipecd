import { Story } from "@storybook/react";
import { dummyApplicationSyncState } from "../../__fixtures__/dummy-application";

import { AppSyncStatus, AppSyncStatusProps } from "./";

export default {
  title: "application/AppSyncStatus",
  component: AppSyncStatus,
};

const Template: Story<AppSyncStatusProps> = (args) => (
  <AppSyncStatus {...args} />
);
export const Overview = Template.bind({});
Overview.args = { deploying: false, syncState: dummyApplicationSyncState };

export const Large = Template.bind({});
Large.args = {
  deploying: false,
  size: "large",
  syncState: dummyApplicationSyncState,
};

export const Deploying = Template.bind({});
Large.args = {
  deploying: true,
  syncState: dummyApplicationSyncState,
};
