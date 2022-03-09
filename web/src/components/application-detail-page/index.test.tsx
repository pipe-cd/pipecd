import { MemoryRouter, Route } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import { fetchApplication } from "~/modules/applications";
import { createStore, render } from "~~/test-utils";
import { ApplicationDetailPage } from ".";

describe("ApplicationDetailPage", () => {
  it("should dispatch actions that fetch application and live state when render", () => {
    const store = createStore({});
    render(
      <MemoryRouter
        initialEntries={[`${PAGE_PATH_APPLICATIONS}/application-1`]}
        initialIndex={0}
      >
        <Route
          exact
          path={`${PAGE_PATH_APPLICATIONS}/:applicationId`}
          component={ApplicationDetailPage}
        />
      </MemoryRouter>,
      { store }
    );

    expect(store.getActions()).toMatchObject([
      { type: fetchApplication.pending.type },
    ]);
  });
});
