import { setupServer } from "msw/node";
import {
  listApplicationsHandler,
  updateApplicationHandler,
} from "~/mocks/services/application";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen } from "~~/test-utils";
import EditApplicationDrawer from ".";
import { UI_TEXT_SAVE } from "~/constants/ui-text";

const server = setupServer(updateApplicationHandler, listApplicationsHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("Show target application info ", async () => {
  render(
    <EditApplicationDrawer
      onUpdated={() => null}
      application={dummyApplication}
      open={true}
      onClose={() => null}
    />
  );
  expect(
    screen.getByText(`Edit "${dummyApplication.name}"`)
  ).toBeInTheDocument();

  const button = screen.getByRole("button", { name: UI_TEXT_SAVE });
  expect(button).toBeDisabled();
});
