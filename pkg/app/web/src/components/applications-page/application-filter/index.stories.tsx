import { Story } from "@storybook/react";
import { createDecoratorRedux } from "~~/.storybook/redux-decorator";
import { ApplicationFilter, ApplicationFilterProps } from ".";

export default {
  title: "APPLICATION/ApplicationFilter",
  component: ApplicationFilter,
  decorators: [
    createDecoratorRedux({
      environments: {
        entities: {
          "env-1": {
            createdAt: 0,
            desc: "env-1",
            id: "env-1",
            name: "stg",
            projectId: "1",
            updatedAt: 0,
            deletedAt: 0,
            deleted: false,
            disabled: false,
          },
        },
        ids: ["env-1"],
      },
    }),
  ],
  argTypes: {
    onChange: { action: "onChange" },
    onClear: { action: "onClear" },
  },
};

const Template: Story<ApplicationFilterProps> = (args) => (
  <ApplicationFilter {...args} />
);
export const Overview = Template.bind({});
Overview.args = { options: {} };
