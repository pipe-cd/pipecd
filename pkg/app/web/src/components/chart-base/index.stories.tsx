import { ChartBase, ChartBaseProps } from "./";
import { Story } from "@storybook/react/types-6-0";
import chartColor from "@material-ui/core/colors/indigo";

export default {
  title: "insights/ChartBase",
  component: ChartBase,
};

const randData = Array.from(new Array(20)).map((_, v) => ({
  value: Math.floor(Math.random() * 20 + 10),
  timestamp: new Date(`2020/10/${v + 5}`).getTime(),
}));

const Template: Story<ChartBaseProps> = (args) => <ChartBase {...args} />;

export const Overview = Template.bind({});
Overview.args = {
  data: [{ name: "data name", points: randData }],
  title: "Title",
  xName: "xName",
  yName: "yName",
  lineColor: chartColor[500],
  areaColor: chartColor[200],
};
