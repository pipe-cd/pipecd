import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { dummyApplicationLiveState } from "~/__fixtures__/dummy-application-live-state";
import { ApplicationStateView, ApplicationStateViewProps } from "./";

export default {
  title: "APPLICATION/ApplicationStateView",
  component: ApplicationStateView,
  decorators: [
    createDecoratorRedux({
      applicationLiveState: {
        entities: {
          [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
        },
        ids: [dummyApplicationLiveState.applicationId],
      },
    }),
  ],
};

const Template: Story<ApplicationStateViewProps> = (args) => (
  <ApplicationStateView {...args} />
);
export const Overview = Template.bind({});
Overview.args = { applicationId: dummyApplicationLiveState.applicationId };
