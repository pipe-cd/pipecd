import userEvent from "@testing-library/user-event";
import { UI_TEXT_CLEAR } from "~/constants/ui-text";
import { ApplicationKind, ApplicationSyncStatus } from "~/modules/applications";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen, waitFor } from "~~/test-utils";
import { ApplicationFilter } from ".";
import { listApplicationsHandler } from "~/mocks/services/application";
import { setupServer } from "msw/node";

const server = setupServer(listApplicationsHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const initialState = {
  applications: {
    ids: [dummyApplication.id],
    entities: { [dummyApplication.id]: dummyApplication },
  },
};

test("Change filter values", async () => {
  const onChange = jest.fn();
  render(
    <ApplicationFilter onChange={onChange} onClear={() => null} options={{}} />,
    {
      initialState,
    }
  );

  userEvent.type(
    screen.getByRole("combobox", { name: "Application Name" }),
    dummyApplication.name
  );
  await waitFor(() => {
    userEvent.click(
      screen.getByRole("option", { name: dummyApplication.name })
    );
  });

  expect(onChange).toHaveBeenCalledWith({ name: dummyApplication.name });
  onChange.mockClear();

  userEvent.click(screen.getByRole("combobox", { name: /kind/i }));
  userEvent.click(screen.getByRole("option", { name: /kubernetes/i }));

  expect(onChange).toHaveBeenCalledWith({ kind: ApplicationKind.KUBERNETES });
  onChange.mockClear();

  userEvent.click(screen.getByRole("combobox", { name: /sync status/i }));
  userEvent.click(screen.getByRole("option", { name: /synced/i }));

  expect(onChange).toHaveBeenCalledWith({
    syncStatus: ApplicationSyncStatus.SYNCED,
  });
  onChange.mockClear();

  userEvent.click(screen.getByRole("combobox", { name: /active status/i }));
  userEvent.click(screen.getByRole("option", { name: /enabled/i }));

  expect(onChange).toHaveBeenCalledWith({ activeStatus: "enabled" });
  onChange.mockClear();
});

test("Click clear filter", () => {
  const onClear = jest.fn();
  render(
    <ApplicationFilter onChange={() => null} onClear={onClear} options={{}} />,
    {
      initialState,
    }
  );

  userEvent.click(screen.getByRole("button", { name: UI_TEXT_CLEAR }));

  expect(onClear).toHaveBeenCalled();
});
