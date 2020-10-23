import React from "react";
import { createStore, render, screen } from "../../test-utils";
import { fireEvent } from "@testing-library/react";
import { ApplicationStateView } from "./application-state-view";
import { dummyApplicationLiveState } from "../__fixtures__/dummy-application-live-state";
import { clearError } from "../modules/applications-live-state";

it("shows refresh button if live state fetching has error", () => {
  const store = createStore({
    applicationLiveState: {
      entities: {
        [dummyApplicationLiveState.applicationId]: dummyApplicationLiveState,
      },
      ids: [dummyApplicationLiveState.applicationId],
      hasError: true,
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

  expect(screen.getByText("Sorry, something went wrong.")).toBeInTheDocument();

  expect(store.getActions()).toEqual([]);

  fireEvent.click(screen.getByRole("button", { name: "REFRESH" }));

  expect(store.getActions()).toMatchObject([{ type: clearError.type }]);
});
