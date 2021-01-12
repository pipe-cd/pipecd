import { configureStore } from "@reduxjs/toolkit";
import React from "react";
import { MemoryRouter, Route } from "react-router-dom";
import { render, waitFor } from "../../../test-utils";
import { server } from "../../mocks/server";
import { reducers } from "../../modules";
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
    const store = configureStore({
      reducer: reducers,
      preloadedState: {
        environments: {
          entities: { [dummyEnv.id]: dummyEnv },
          ids: [dummyEnv.id],
        },
        pipeds: {
          entities: { [dummyPiped.id]: dummyPiped },
          ids: [dummyPiped.id],
        },
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
