import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import { updateStaticAdmin } from "~/modules/project";
import {
  act,
  createStore,
  render,
  screen,
  waitFor,
  waitForElementToBeRemoved,
} from "~~/test-utils";
import { StaticAdminForm } from ".";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

it("should shows current username", () => {
  render(<StaticAdminForm />, {
    initialState: {
      project: {
        username: "pipe-user",
        staticAdminDisabled: false,
      },
    },
  });

  expect(screen.getByText("pipe-user")).toBeInTheDocument();
});

it("should dispatch action that update static admin when input fields and click submit button", async () => {
  const store = createStore({
    project: {
      username: "pipe-user",
      staticAdminDisabled: false,
    },
  });

  render(<StaticAdminForm />, {
    store,
  });

  userEvent.click(
    screen.getByRole("button", { name: "edit static admin user" })
  );

  await waitFor(() => screen.getByText("Edit Static Admin"));

  userEvent.type(screen.getByRole("textbox", { name: /username/i }), "-new");
  userEvent.type(screen.getByLabelText(/password/i), "new-password");

  act(() => {
    userEvent.click(screen.getByRole("button", { name: /save/i }));
  });

  await waitForElementToBeRemoved(() => screen.getByText("Edit Static Admin"));

  expect(store.getActions()).toEqual(
    expect.arrayContaining([
      expect.objectContaining({
        type: updateStaticAdmin.pending.type,
        meta: expect.objectContaining({
          arg: {
            username: "pipe-user-new",
            password: "new-password",
          },
        }),
      }),
      expect.objectContaining({
        type: updateStaticAdmin.fulfilled.type,
      }),
    ])
  );
});
