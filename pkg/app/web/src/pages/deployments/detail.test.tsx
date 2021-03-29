import React from "react";
import { MemoryRouter, Route } from "react-router-dom";
import { createReduxStore, render, waitFor } from "../../../test-utils";
import { server } from "../../mocks/server";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { DeploymentDetailPage } from "./detail";

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
      environments: {
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      },
      pipeds: {
        entities: { [dummyPiped.id]: dummyPiped },
        ids: [dummyPiped.id],
        registeredPiped: null,
        updating: false,
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
