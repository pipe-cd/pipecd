import { Header } from "./";
import { createStore } from "test-utils";
import { Provider } from "react-redux";
import { Role } from "../../../../../../bazel-bin/pkg/app/web/model/role_pb";
import { Story } from "@storybook/react";

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
