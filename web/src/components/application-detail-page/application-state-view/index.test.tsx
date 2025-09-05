import userEvent from "@testing-library/user-event";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen } from "~~/test-utils";
import { ApplicationStateView } from ".";

it("shows refresh button if live state fetching has error", () => {
  const mockFn = jest.fn();
  render(
    <ApplicationStateView
      app={dummyApplication}
      hasError={true}
      refetchLiveState={mockFn}
    />
  );

  expect(
    screen.getByText("It was unable to fetch the latest state of application.")
  ).toBeInTheDocument();
  expect(
    screen.getByRole("button", { name: UI_TEXT_REFRESH })
  ).toBeInTheDocument();
  userEvent.click(screen.getByRole("button", { name: UI_TEXT_REFRESH }));
  expect(mockFn).toHaveBeenCalledTimes(1);
});
