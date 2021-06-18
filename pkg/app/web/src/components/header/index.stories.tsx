import { Header } from "./";
import { createStore } from "~~/test-utils";
import { Provider } from "react-redux";
import { Story } from "@storybook/react";
import { Role } from "~/modules/me";

export default {
  title: "COMMON/Header",
  component: Header,
};

export const Overview: Story = () => (
  <Provider store={createStore()}>
    <Header />
  </Provider>
);

export const LoggedIn: Story = () => (
  <Provider
    store={createStore({
      me: {
        avatarUrl: "",
        isLogin: true,
        projectId: "pipecd",
        projectRole: Role.ProjectRole.ADMIN,
        subject: "",
      },
    })}
  >
    <Header />
  </Provider>
);
