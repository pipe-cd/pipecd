import { Story } from "@storybook/react";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { DeleteApplicationDialog, DeleteApplicationDialogProps } from ".";

export default {
  title: "DeleteApplicationDialog",
  component: DeleteApplicationDialog,
  decorators: [
    createDecoratorRedux({
      applications: {
        entities: {
          [dummyApplication.id]: dummyApplication,
        },
        ids: [dummyApplication.id],
      },
      deleteApplication: {
        applicationId: dummyApplication.id,
        deleting: false,
      },
    }),
  ],
  argTypes: {
    onDeleted: {
      action: "onDeleted",
    },
  },
};

const Template: Story<DeleteApplicationDialogProps> = (args) => (
  <DeleteApplicationDialog {...args} />
);

export const Overview = Template.bind({});
Overview.args = {};
