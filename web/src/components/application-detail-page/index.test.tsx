import { MemoryRouter, render, screen, waitFor } from "~~/test-utils";
import { ApplicationDetailPage } from ".";
import { server } from "~/mocks/server";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { Routes as ReactRoutes } from "react-router-dom";
import { Route } from "react-router-dom";
import { PAGE_PATH_APPLICATIONS } from "~/constants/path";
import userEvent from "@testing-library/user-event";

beforeAll(() => {
  server.listen();
});
afterEach(() => {
  server.resetHandlers();
});
afterAll(() => {
  server.close();
});

describe("ApplicationDetailPage", () => {
  it("should have Icon menu and option disabled, encrypt", async () => {
    server.use();
    render(
      <MemoryRouter
        initialEntries={[`${PAGE_PATH_APPLICATIONS}/${dummyApplication.id}`]}
        initialIndex={0}
      >
        <ReactRoutes>
          <Route
            path={`${PAGE_PATH_APPLICATIONS}/:applicationId`}
            element={<ApplicationDetailPage />}
          />
        </ReactRoutes>
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
    });

    expect(
      screen.getByRole("button", { name: "Open menu" })
    ).toBeInTheDocument();
    userEvent.click(screen.getByRole("button", { name: "Open menu" }));
    expect(
      screen.getByRole("menuitem", { name: "Disable" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("menuitem", { name: "Encrypt Secret" })
    ).toBeInTheDocument();
    expect(
      screen.getByRole("menuitem", { name: "Delete" })
    ).toBeInTheDocument();
  });
});
