import userEvent from "@testing-library/user-event";
import { UI_TEXT_FILTER } from "~/constants/ui-text";
import { resourcesApplicationList } from "~/__fixtures__/dummy-application-live-state";
import { act, render, screen } from "~~/test-utils";
import { LiveStateView } from ".";

test("render resources on tab local", async () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  await userEvent.click(screen.getByRole("tab", { name: "local" }));

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(2);
});

test("render resources on tab kubernetes", async () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  await userEvent.click(screen.getByRole("tab", { name: "kubernetes" }));

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(3);
});

test("filter resources on tab kubernetes", async () => {
  render(<LiveStateView resources={resourcesApplicationList} />, {});

  await userEvent.click(screen.getByRole("tab", { name: "kubernetes" }));
  await userEvent.click(screen.getByRole("button", { name: UI_TEXT_FILTER }));
  await userEvent.click(screen.getByRole("checkbox", { name: "ReplicaSet" }));

  await act(async () => {
    await userEvent.click(screen.getByRole("button", { name: "APPLY" }));
  });

  expect(screen.queryAllByTestId("application-resource")).toHaveLength(2);
});

test("show a remaining deploy target when the selected one disappears", () => {
  const { rerender } = render(
    <LiveStateView resources={resourcesApplicationList} />,
    {}
  );
  expect(screen.getByRole("tab", { name: "kubernetes" })).toHaveAttribute(
    "aria-selected",
    "true"
  );

  rerender(
    <LiveStateView
      resources={resourcesApplicationList.filter(
        (r) => r.deployTarget !== "kubernetes"
      )}
    />
  );

  expect(screen.getByRole("tab", { name: "local" })).toHaveAttribute(
    "aria-selected",
    "true"
  );
  expect(screen.queryAllByTestId("application-resource")).toHaveLength(2);
});

test("render nothing when every deploy target disappears", () => {
  const { rerender } = render(
    <LiveStateView resources={resourcesApplicationList} />,
    {}
  );

  rerender(<LiveStateView resources={[]} />);

  expect(screen.queryAllByRole("tab")).toHaveLength(0);
  expect(screen.queryAllByTestId("application-resource")).toHaveLength(0);
});

test("keep the shown deploy target when an earlier-sorted one appears", () => {
  const { rerender } = render(
    <LiveStateView resources={resourcesApplicationList} />,
    {}
  );
  expect(screen.getByRole("tab", { name: "kubernetes" })).toHaveAttribute(
    "aria-selected",
    "true"
  );

  rerender(
    <LiveStateView
      resources={[
        ...resourcesApplicationList,
        {
          ...resourcesApplicationList[0],
          id: "alpha-1",
          deployTarget: "alpha",
        },
      ]}
    />
  );

  expect(screen.getByRole("tab", { name: "kubernetes" })).toHaveAttribute(
    "aria-selected",
    "true"
  );
});
