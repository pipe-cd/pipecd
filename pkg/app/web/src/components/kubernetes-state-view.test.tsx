import React from "react";
import { render, screen } from "../../test-utils";
import { resourcesList } from "../__fixtures__/dummy-application-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";

test("render resources", () => {
  render(
    <KubernetesStateView
      resources={resourcesList}
      showKinds={["Pod", "ReplicaSet"]}
    />,
    {}
  );

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(3);
});

test("filter resources", () => {
  render(
    <KubernetesStateView resources={resourcesList} showKinds={["Pod"]} />,
    {}
  );

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(2);
  expect(screen.queryByText("ReplicaSet")).not.toBeInTheDocument();
});
