import { Story } from "@storybook/react";
import { Provider } from "react-redux";
import { createStore } from "~~/test-utils";
import { RBACForm } from ".";

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
