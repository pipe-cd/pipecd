import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import { disableApplication, enableApplication } from "~/modules/applications";
import { setDeletingAppId } from "~/modules/delete-application";
import { generateSealedSecret } from "~/modules/sealed-secret";
import { setUpdateTargetId } from "~/modules/update-application";
import { dummyApplication } from "~/__fixtures__/dummy-application";
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

test("delete", () => {
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

  expect(store.getActions()).toEqual(
    expect.arrayContaining([
      expect.objectContaining({
        type: setDeletingAppId.type,
        payload: dummyApplication.id,
      }),
    ])
  );
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

  await waitFor(() =>
    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          type: disableApplication.pending.type,
          meta: expect.objectContaining({
            arg: {
              applicationId: dummyApplication.id,
            },
          }),
        }),
      ])
    )
  );

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

  await waitFor(() =>
    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({
          type: generateSealedSecret.pending.type,
          meta: expect.objectContaining({
            arg: {
              base64Encoding: false,
              data: "secret data",
              pipedId: dummyApplication.pipedId,
            },
          }),
        }),
      ])
    )
  );
});
