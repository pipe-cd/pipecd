import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import {
  act,
  render,
  screen,
  waitFor,
  waitForElementToBeRemoved,
} from "~~/test-utils";
import { StaticAdminForm } from ".";
import { UPDATE_STATIC_ADMIN_INFO_SUCCESS } from "~/constants/toast-text";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

it("should shows current username", async () => {
  await act(async () => {
    await render(<StaticAdminForm />);
  });

  await waitFor(() => {
    expect(screen.getByText("static-admin-user")).toBeInTheDocument();
  });
});

it("should show success message when update static admin", async () => {
  render(<StaticAdminForm />);

  await waitFor(() => screen.getByText("static-admin-user"));

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

  await waitFor(() =>
    expect(
      screen.getByText(UPDATE_STATIC_ADMIN_INFO_SUCCESS)
    ).toBeInTheDocument()
  );
});
