import { configureStore } from "@reduxjs/toolkit";
import { render, waitFor } from "@testing-library/react";
import React from "react";
import { Provider } from "react-redux";
import { MemoryRouter, Route } from "react-router";
import * as deploymentsApi from "../../api/deployments";
import { reducers } from "../../modules";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { DeploymentDetailPage } from "./detail";
import { dummyPiped } from "../../__fixtures__/dummy-piped";

jest.mock("../../api/deployments");

describe("DeploymentDetailPage", () => {
  test("fetch a deployment data and show that data", async () => {
    jest.spyOn(deploymentsApi, "getDeployment").mockReturnValue(
      Promise.resolve({
        deployment: dummyDeployment,
      })
    );

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
      <Provider store={store}>
        <MemoryRouter
          initialEntries={[`/deployments/${dummyDeployment.id}`]}
          initialIndex={0}
        >
          <Route path="/deployments/:deploymentId">
            <DeploymentDetailPage />
          </Route>
        </MemoryRouter>
      </Provider>
    );

    expect(deploymentsApi.getDeployment).toHaveBeenCalledWith({
      deploymentId: dummyDeployment.id,
    });

    await waitFor(() => getByText("SUCCESS"));
  });
});
