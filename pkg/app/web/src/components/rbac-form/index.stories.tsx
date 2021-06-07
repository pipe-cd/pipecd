import { RBACForm } from "./";
import { Provider } from "react-redux";
import { createStore } from "../../../test-utils";
import { Story } from "@storybook/react";

export default {
  title: "RBACForm",
  component: RBACForm,
};

export const Overview: Story = () => (
  <Provider
    store={createStore({
      project: {
        teams: { admin: "admin", editor: "editor", viewer: "viewer" },
      },
    })}
  >
    <RBACForm />
  </Provider>
);

export const NotConfigured: Story = () => (
  <Provider
    store={createStore({
      project: {
        teams: null,
      },
    })}
  >
    <RBACForm />
  </Provider>
);
