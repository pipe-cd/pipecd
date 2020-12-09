import { action } from "@storybook/addon-actions";
import React from "react";
import { WeekPicker } from "./week-picker";

export default {
  title: "COMMON/WeekPicker",
  component: WeekPicker,
};

export const overview: React.FC = () => (
  <WeekPicker
    value={new Date()}
    label="Week Picker"
    onChange={action("onChange")}
  />
);
