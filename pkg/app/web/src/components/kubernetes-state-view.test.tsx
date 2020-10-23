import { fireEvent } from "@testing-library/react";
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

  fireEvent.click(screen.getByRole("button", { name: "FILTER" }));
  fireEvent.click(screen.getByRole("checkbox", { name: "ReplicaSet" }));
  fireEvent.click(screen.getByRole("button", { name: "APPLY" }));

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(2);
});
