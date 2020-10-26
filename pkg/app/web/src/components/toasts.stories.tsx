import React from "react";
import { Toasts } from "./toasts";
import { Provider } from "react-redux";
import { createStore } from "../../test-utils";

export default {
  title: "COMMON/Toasts",
  component: Toasts,
};

export const overview: React.FC = () => (
  <Provider
    store={createStore({
      toasts: {
        entities: {
          "1": {
            id: "1",
            message: "toast message",
            severity: "error",
          },
        },
        ids: ["1"],
      },
    })}
  >
    <Toasts />
  </Provider>
);

export const url: React.FC = () => (
  <Provider
    store={createStore({
      toasts: {
        entities: {
          "1": {
            id: "1",
            message: "toast message",
            severity: "success",
            to: "/",
          },
        },
        ids: ["1"],
      },
    })}
  >
    <Toasts />
  </Provider>
);
