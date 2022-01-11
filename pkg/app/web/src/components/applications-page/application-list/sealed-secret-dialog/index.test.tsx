import { waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { setupServer } from "msw/node";
import { generateApplicationSealedSecretHandler } from "~/mocks/services/piped";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { createReduxStore, render, screen } from "~~/test-utils/index";
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
      applicationId={dummyApplication.id}
      onClose={() => null}
      open
    />,
    {
      initialState: {
        applications: {
          entities: {
            [dummyApplication.id]: dummyApplication,
          },
          ids: [dummyApplication.id],
        },
      },
    }
  );

  expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
});

test("cancel", () => {
  const onClose = jest.fn();
  render(
    <SealedSecretDialog
      applicationId={dummyApplication.id}
      onClose={onClose}
      open
    />,
    {
      initialState: {
        applications: {
          entities: {
            [dummyApplication.id]: dummyApplication,
          },
          ids: [dummyApplication.id],
        },
      },
    }
  );

  userEvent.click(screen.getByRole("button", { name: "Cancel" }));
  expect(onClose).toHaveBeenCalled();
});

test("Generate sealed secret", async () => {
  const store = createReduxStore({
    applications: {
      entities: {
        [dummyApplication.id]: dummyApplication,
      },
      ids: [dummyApplication.id],
      adding: false,
      disabling: {},
      loading: false,
      syncing: {},
      addedApplicationId: null,
      fetchApplicationError: null,
      allLabels: {},
    },
  });

  server.use(generateApplicationSealedSecretHandler);
  render(
    <SealedSecretDialog
      applicationId={dummyApplication.id}
      onClose={jest.fn()}
      open
    />,
    {
      store,
    }
  );

  userEvent.type(
    screen.getByRole("textbox", { name: "Secret Data" }),
    "secret"
  );
  userEvent.click(screen.getByRole("button", { name: "Encrypt" }));

  await waitFor(() =>
    expect(screen.getByText("Encrypted secret data")).toBeInTheDocument()
  );

  userEvent.click(screen.getByRole("button", { name: "Close" }));

  await waitFor(() =>
    expect(screen.queryByText("Encrypted secret data")).not.toBeInTheDocument()
  );
});
