import userEvent from "@testing-library/user-event";
import { UI_TEXT_REFRESH } from "~/constants/ui-text";
import { fetchApplicationStateById } from "~/modules/applications-live-state";
import { dummyApplicationLiveState } from "~/__fixtures__/dummy-application-live-state";
import { createStore, render, screen } from "~~/test-utils";
import { ApplicationStateView } from ".";

it("shows refresh button if live state fetching has error", () => {
  const store = createStore({
    applicationLiveState: {
      entities: {
        [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
      },
      ids: [dummyApplicationLiveState.applicationId],
      hasError: {
        [dummyApplicationLiveState.applicationId]: true,
      },
    },
  });
  render(
    <ApplicationStateView
      applicationId={dummyApplicationLiveState.applicationId}
    />,
    {
      store,
    }
  );

  expect(
    screen.getByText("It was unable to fetch the latest state of application.")
  ).toBeInTheDocument();

  expect(store.getActions()).toEqual([]);

  userEvent.click(screen.getByRole("button", { name: UI_TEXT_REFRESH }));

  expect(store.getActions()).toMatchObject([
    { type: fetchApplicationStateById.pending.type },
  ]);
});
