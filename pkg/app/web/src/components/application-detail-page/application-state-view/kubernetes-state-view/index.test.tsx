import userEvent from "@testing-library/user-event";
import { UI_TEXT_FILTER } from "~/constants/ui-text";
import { resourcesList } from "~/__fixtures__/dummy-application-live-state";
import { render, screen } from "~~/test-utils";
import { KubernetesStateView } from ".";

test("render resources", () => {
  render(<KubernetesStateView resources={resourcesList} />, {});

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(3);
});

test("filter resources", () => {
  render(<KubernetesStateView resources={resourcesList} />, {});

  userEvent.click(screen.getByRole("button", { name: UI_TEXT_FILTER }));
  userEvent.click(screen.getByRole("checkbox", { name: "ReplicaSet" }));
  userEvent.click(screen.getByRole("button", { name: "APPLY" }));

  expect(screen.queryAllByTestId("kubernetes-resource")).toHaveLength(2);
});
