import { Story } from "@storybook/react";
import { ProjectSettingLabeledText, ProjectSettingLabeledTextProps } from ".";

export default {
  title: "ProjectSettingLabeledText",
  component: ProjectSettingLabeledText,
};

const Template: Story<ProjectSettingLabeledTextProps> = (args) => (
  <ProjectSettingLabeledText {...args} />
);
export const Overview = Template.bind({});
Overview.args = { label: "label", value: "value" };
