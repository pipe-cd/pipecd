import userEvent from "@testing-library/user-event";
import { setupServer } from "msw/node";
import { createReduxStore, render, screen, waitFor } from "test-utils";
import { listDeploymentConfigTemplatesHandler } from "../../mocks/services/deployment-config";
import { dummyApplication } from "../../__fixtures__/dummy-application";
import { DeploymentConfigForm } from "./";

const server = setupServer(listDeploymentConfigTemplatesHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const preloadedState = {
  deploymentConfigs: {
    targetApplicationId: dummyApplication.id,
    templates: {},
  },
};

test("Change template", async () => {
  const store = createReduxStore(preloadedState);

  render(<DeploymentConfigForm onSkip={() => null} />, {
    store,
  });

  await waitFor(() => expect(screen.getByRole("button", { name: /simple/i })));
  userEvent.click(screen.getByRole("button", { name: /simple/i }));
  userEvent.click(screen.getByRole("option", { name: /canary/i }));
  await waitFor(() =>
    expect(screen.getByText(/deploy progressively with canary strategy/i))
  );
});

test("Skip", () => {
  const onSkip = jest.fn();
  const store = createReduxStore(preloadedState);
  render(<DeploymentConfigForm onSkip={onSkip} />, {
    store,
  });

  userEvent.click(screen.getByRole("button", { name: /skip/i }));
  expect(onSkip).toHaveBeenCalled();
});
