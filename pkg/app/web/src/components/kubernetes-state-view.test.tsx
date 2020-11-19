import userEvent from "@testing-library/user-event";
import React from "react";
import { render, screen } from "../../test-utils";
import { resourcesList } from "../__fixtures__/dummy-application-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";

test("render resources", () => {
  render(<KubernetesStateView resources={resourcesList} />, {});

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(3);
});

test("filter resources", () => {
  render(<KubernetesStateView resources={resourcesList} />, {});

  userEvent.click(screen.getByRole("button", { name: "FILTER" }));
  userEvent.click(screen.getByRole("checkbox", { name: "ReplicaSet" }));
  userEvent.click(screen.getByRole("button", { name: "APPLY" }));

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(2);
});
