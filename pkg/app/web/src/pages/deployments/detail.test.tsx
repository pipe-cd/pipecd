import { configureStore } from "@reduxjs/toolkit";
import React from "react";
import { MemoryRouter, Route } from "react-router";
import { render, waitFor } from "../../../test-utils";
import * as deploymentsApi from "../../api/deployments";
import { reducers } from "../../modules";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { DeploymentDetailPage } from "./detail";

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

    expect(deploymentsApi.getDeployment).toHaveBeenCalledWith({
      deploymentId: dummyDeployment.id,
    });

    await waitFor(() => getByText("SUCCESS"));
  });
});
