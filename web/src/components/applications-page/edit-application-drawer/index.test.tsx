import { setupServer } from "msw/node";
import {
  listApplicationsHandler,
  updateApplicationHandler,
} from "~/mocks/services/application";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createReduxStore, render, screen } from "~~/test-utils";
import EditApplicationDrawer from ".";

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
  pipeds: {
    ids: [dummyPiped.id],
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    registeredPiped: null,
    updating: false,
    releasedVersions: [],
    breakingChangesNote: "",
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
  expect(
    screen.getByText(`${dummyPiped.name} (${dummyPiped.id})`)
  ).toBeInTheDocument();
});
