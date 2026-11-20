import { waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { setupServer } from "msw/node";
import { generateApplicationSealedSecretHandler } from "~/mocks/services/piped";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { render, screen } from "~~/test-utils/index";
import { SealedSecretDialog } from ".";

const server = setupServer();

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("render", () => {
  render(
    <SealedSecretDialog
      application={dummyApplication}
      onClose={() => null}
      open
    />
  );

  expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
});

test("cancel", async () => {
  const onClose = jest.fn();
  render(
    <SealedSecretDialog application={dummyApplication} onClose={onClose} open />
  );

  await userEvent.click(screen.getByRole("button", { name: "Cancel" }));
  expect(onClose).toHaveBeenCalled();
});

test("Generate sealed secret", async () => {
  server.use(generateApplicationSealedSecretHandler);
  const handleClose = jest.fn();
  render(
    <SealedSecretDialog
      application={dummyApplication}
      onClose={handleClose}
      open
    />
  );

  await userEvent.type(
    screen.getByRole("textbox", { name: "Secret Data" }),
    "secret"
  );
  await userEvent.click(screen.getByRole("button", { name: "Encrypt" }));

  await waitFor(() =>
    expect(screen.getByText("Encrypted secret data")).toBeInTheDocument()
  );

  await userEvent.click(screen.getByRole("button", { name: "Close" }));

  await waitFor(() => expect(handleClose).toHaveBeenCalled());
});
