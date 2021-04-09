import * as React from "react";
import { Header } from "./";
import { createStore } from "../../../test-utils";
import { Provider } from "react-redux";
import { Role } from "../../../../../../bazel-bin/pkg/app/web/model/role_pb";

export default {
  title: "COMMON/Header",
  component: Header,
};

export const overview: React.FC = () => {
  const store = createStore({});
  return (
    <Provider store={store}>
      <Header />
    </Provider>
  );
};

export const loggedIn: React.FC = () => {
  const store = createStore({
    me: {
      avatarUrl: "",
      isLogin: true,
      projectId: "pipecd",
      projectRole: Role.ProjectRole.ADMIN,
      subject: "",
    },
  });
  return (
    <Provider store={store}>
      <Header />
    </Provider>
  );
};
