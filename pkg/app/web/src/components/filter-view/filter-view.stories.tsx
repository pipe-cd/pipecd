import { action } from "@storybook/addon-actions";
import React from "react";
import { FilterView } from "./filter-view";

export default {
  title: "FilterView",
  component: FilterView,
};

export const overview: React.FC = () => (
  <FilterView onClear={action("onClear")}>filter</FilterView>
);
