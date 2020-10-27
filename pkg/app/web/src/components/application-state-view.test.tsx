import React from "react";
import { createStore, render, screen } from "../../test-utils";
import { fireEvent } from "@testing-library/react";
import { ApplicationStateView } from "./application-state-view";
import { dummyApplicationLiveState } from "../__fixtures__/dummy-application-live-state";
import { clearError } from "../modules/applications-live-state";
import { UI_TEXT_REFRESH } from "../constants/ui-text";

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

  fireEvent.click(screen.getByRole("button", { name: UI_TEXT_REFRESH }));

  expect(store.getActions()).toMatchObject([{ type: clearError.type }]);
});
