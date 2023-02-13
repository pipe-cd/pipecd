import { MemoryRouter, Route } from "react-router-dom";
import { server } from "~/mocks/server";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createReduxStore, render, waitFor } from "~~/test-utils";
import { DeploymentDetailPage } from ".";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("DeploymentDetailPage", () => {
  test("fetch a deployment data and show that data", async () => {
    const store = createReduxStore({
      pipeds: {
        entities: { [dummyPiped.id]: dummyPiped },
        ids: [dummyPiped.id],
        registeredPiped: null,
        updating: false,
        releasedVersions: [],
      },
    });

    const { getByText } = render(
      <MemoryRouter
        initialEntries={[`/deployments/${dummyDeployment.id}`]}
        initialIndex={0}
      >
        <Route path="/deployments/:deploymentId">
          <DeploymentDetailPage />
        </Route>
      </MemoryRouter>,
      { store }
    );

    await waitFor(() => getByText("SUCCESS"));
  });
});
