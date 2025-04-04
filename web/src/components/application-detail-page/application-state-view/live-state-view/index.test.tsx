import userEvent from "@testing-library/user-event";
import { UI_TEXT_FILTER } from "~/constants/ui-text";
import { resourcesApplicationList } from "~/__fixtures__/dummy-application-live-state";
import { render, screen } from "~~/test-utils";
import { LiveStateView } from ".";

test("render resources on tab local", () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  userEvent.click(screen.getByRole("tab", { name: "local" }));

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(2);
});

test("render resources on tab kubernetes", () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  userEvent.click(screen.getByRole("tab", { name: "kubernetes" }));

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(3);
});

test("filter resources on tab kubernetes", () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  userEvent.click(screen.getByRole("tab", { name: "kubernetes" }));
  userEvent.click(screen.getByRole("button", { name: UI_TEXT_FILTER }));
  userEvent.click(screen.getByRole("checkbox", { name: "ReplicaSet" }));
  userEvent.click(screen.getByRole("button", { name: "APPLY" }));

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(2);
});
