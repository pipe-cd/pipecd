import { Box } from "@material-ui/core";
import { Story } from "@storybook/react/types-6-0";
import React from "react";
import { ChangeFailureRateChart, ChangeFailureRateChartProps } from ".";

export default {
  title: "insights/ChangeFailureRateChart",
  component: ChangeFailureRateChart,
};

const randData = Array.from(new Array(20)).map((_, v) => ({
  value: Math.floor(Math.random() * 20 + 10),
  timestamp: new Date(`2020/10/${v + 5}`).getTime(),
}));

const Template: Story<ChangeFailureRateChartProps> = (args) => (
  <Box width={800}>
    <ChangeFailureRateChart data={args.data} />
  </Box>
);
export const Overview = Template.bind({});
Overview.args = {
  data: [{ name: "Average Change Failure Rate", points: randData }],
};

export const NoData = Template.bind({});
NoData.args = {
  data: [],
};
