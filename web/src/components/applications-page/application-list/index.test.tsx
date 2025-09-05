import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import { enableApplication } from "~/modules/applications";
import { setUpdateTargetId } from "~/modules/update-application";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import * as applicationApi from "~/api/applications";
import * as pipedApi from "~/api/piped";
import {
  createStore,
  render,
  screen,
  waitFor,
  MemoryRouter,
} from "~~/test-utils";
import { ApplicationList } from ".";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const state = {
  applications: {
    entities: {
      [dummyApplication.id]: dummyApplication,
    },
    ids: [dummyApplication.id],
  },
};

test("delete", async () => {
  jest.spyOn(applicationApi, "deleteApplication");

  const store = createStore(state);
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} />
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Delete" }));

  await waitFor(() =>
    expect(
      screen.getByRole("heading", { name: "Delete Application" })
    ).toBeInTheDocument()
  );

  userEvent.click(screen.getByRole("button", { name: "Delete" }));

  await waitFor(() =>
    expect(
      screen.queryByRole("heading", { name: "Delete Application" })
    ).not.toBeInTheDocument()
  );
  expect(applicationApi.deleteApplication).toHaveBeenCalledWith({
    applicationId: dummyApplication.id,
  });
});

test("show specific page", async () => {
  const apps = [...new Array(30)].map((_, i) => ({
    ...dummyApplication,
    id: `${dummyApplication.id}${i}`,
  }));
  const store = createStore({
    applications: {
      entities: apps.reduce((prev, current) => {
        return { ...prev, [current.id]: current };
      }, {}),
      ids: apps.map((app) => app.id),
    },
  });
  render(
    <MemoryRouter>
      <ApplicationList currentPage={2} />
    </MemoryRouter>,
    {
      store,
    }
  );

  const items = await screen.findAllByText(dummyApplication.name);
  expect(items).toHaveLength(10);
});

test("edit", () => {
  const store = createStore(state);
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} />
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Edit" }));

  expect(store.getActions()).toEqual(
    expect.arrayContaining([
      expect.objectContaining({
        type: setUpdateTargetId.type,
        payload: dummyApplication.id,
      }),
    ])
  );
});

test("disable", async () => {
  jest.spyOn(applicationApi, "disableApplication");
  const store = createStore(state);
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} />
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Disable" }));

  userEvent.click(screen.getByRole("button", { name: "Disable" }));

  await waitFor(() => {
    expect(applicationApi.disableApplication).toHaveBeenCalledWith({
      applicationId: dummyApplication.id,
    });
  });

  await waitFor(() => {
    expect(
      screen.queryByRole("button", { name: "Disable" })
    ).not.toBeInTheDocument();
  });
});

test("enable", async () => {
  const store = createStore({
    applications: {
      entities: {
        [dummyApplication.id]: { ...dummyApplication, disabled: true },
      },
      ids: [dummyApplication.id],
    },
  });
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} />
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Enable" }));

  await waitFor(() =>
    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          type: enableApplication.pending.type,
          meta: expect.objectContaining({
            arg: {
              applicationId: dummyApplication.id,
            },
          }),
        }),
      ])
    )
  );
});

test("Encrypt Secret", async () => {
  jest.spyOn(pipedApi, "generateApplicationSealedSecret");
  const store = createStore(state);
  render(
    <MemoryRouter>
      <ApplicationList currentPage={1} />
    </MemoryRouter>,
    {
      store,
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Open menu" }));
  userEvent.click(screen.getByRole("menuitem", { name: "Encrypt Secret" }));

  userEvent.type(
    screen.getByRole("textbox", { name: "Secret Data" }),
    "secret data"
  );

  userEvent.click(screen.getByRole("button", { name: "Encrypt" }));

  await waitFor(() => {
    expect(pipedApi.generateApplicationSealedSecret).toHaveBeenCalledWith({
      base64Encoding: false,
      data: "secret data",
      pipedId: dummyApplication.pipedId,
    });
  });

  await waitFor(() => {
    expect(screen.getByText("Encrypted secret data")).toBeInTheDocument();
  });
});
