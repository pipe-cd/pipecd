//import { waitFor } from "@testing-library/react";
//import userEvent from "@testing-library/user-event";
import { setupServer } from "msw/node";
//import { UI_TEXT_SAVE } from "~/constants/ui-text";
import {
  listApplicationsHandler,
  updateApplicationHandler,
} from "~/mocks/services/application";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createReduxStore, render, screen } from "~~/test-utils";
import { EditApplicationDrawer } from ".";

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

const initialState = {
  updateApplication: {
    targetId: dummyApplication.id,
    updating: false,
  },
  environments: {
    ids: [dummyEnv.id],
    entities: { [dummyEnv.id]: dummyEnv },
  },
  pipeds: {
    ids: [dummyPiped.id],
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    registeredPiped: null,
    updating: false,
  },
  applications: {
    loading: false,
    adding: false,
    entities: {
      [dummyApplication.id]: dummyApplication,
    },
    fetchApplicationError: null,
    addedApplicationId: null,
    ids: [dummyApplication.id],
    syncing: {},
    disabling: {},
  },
};

test("Show target application info ", () => {
  const store = createReduxStore(initialState);
  render(<EditApplicationDrawer onUpdated={() => null} />, {
    store,
  });

  expect(screen.getByDisplayValue(dummyApplication.name)).toBeInTheDocument();
  expect(screen.getByText(dummyEnv.name)).toBeInTheDocument();
  expect(
    screen.getByText(`${dummyPiped.name} (${dummyPiped.id})`)
  ).toBeInTheDocument();
});

// TODO: Uncomment out after it terns out why pointer-events set to "none"
/*
test("Edit an application ", async () => {
  const store = createReduxStore(initialState);
  render(<EditApplicationDrawer onUpdated={() => null} />, {
    store,
  });

  expect(
    screen.getByRole("heading", {
      name: `Edit "${dummyApplication.name}"`,
    })
  );
  userEvent.type(screen.getByRole("textbox", { name: /^name$/i }), "new-name");
  userEvent.click(screen.getByRole("button", { name: UI_TEXT_SAVE }));

  await waitFor(() =>
    expect(
      screen.queryByRole("heading", {
        name: `Edit "${dummyApplication.name}"`,
      })
    ).not.toBeInTheDocument()
  );
});
*/
