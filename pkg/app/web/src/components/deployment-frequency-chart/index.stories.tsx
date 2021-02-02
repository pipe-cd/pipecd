import { Box } from "@material-ui/core";
import React from "react";
import { DeploymentFrequencyChart } from "./";

export default {
  title: "insights/DeploymentFrequencyChart",
  component: DeploymentFrequencyChart,
};

const randData = Array.from(new Array(20)).map((_, v) => ({
  value: Math.floor(Math.random() * 20 + 10),
  timestamp: new Date(`2020/10/${v + 5}`).getTime(),
}));

export const overview: React.FC = () => (
  <Box width={800}>
    <DeploymentFrequencyChart data={randData} />
  </Box>
);
export const noData: React.FC = () => (
  <Box width={800}>
    <DeploymentFrequencyChart data={[]} />
  </Box>
);
