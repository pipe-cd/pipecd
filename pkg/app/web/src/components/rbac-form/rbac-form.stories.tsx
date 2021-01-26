import React from "react";
import { RBACForm } from "./rbac-form";
import { Provider } from "react-redux";
import { createStore } from "../../../test-utils";

export default {
  title: "RBACForm",
  component: RBACForm,
};

export const overview: React.FC = () => (
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

export const notConfigured: React.FC = () => (
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
