import { Box } from "@material-ui/core";
import { Story } from "@storybook/react/types-6-0";
import React from "react";
import { DeploymentFrequencyChart, DeploymentFrequencyChartProps } from "./";

export default {
  title: "insights/DeploymentFrequencyChart",
  component: DeploymentFrequencyChart,
};

const randData = Array.from(new Array(20)).map((_, v) => ({
  value: Math.floor(Math.random() * 20 + 10),
  timestamp: new Date(`2020/10/${v + 5}`).getTime(),
}));

const Template: Story<DeploymentFrequencyChartProps> = (args) => (
  <Box width={800}>
    <DeploymentFrequencyChart data={args.data} />
  </Box>
);
export const Overview = Template.bind({});
Overview.args = {
  data: [{ name: "application-1", points: randData }],
};

export const MultipleApplication = Template.bind({});
MultipleApplication.args = {
  data: [
    { name: "application-1", points: randData },
    {
      name: "application-2",
      points: randData.map((v) => ({
        ...v,
        value: Math.floor(Math.random() * 30 + 5),
      })),
    },
    {
      name: "application-3",
      points: randData.map((v) => ({
        ...v,
        value: Math.floor(Math.random() * 10 + 15),
      })),
    },
  ],
};

export const NoData = Template.bind({});
NoData.args = {
  data: [],
};
