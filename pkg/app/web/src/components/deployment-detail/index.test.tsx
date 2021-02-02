import React from "react";
import { createStore, render, screen } from "../../../test-utils";
import { DeploymentDetail } from "./";
import { dummyDeployment } from "../../__fixtures__/dummy-deployment";
import { dummyEnv } from "../../__fixtures__/dummy-environment";
import { dummyPiped } from "../../__fixtures__/dummy-piped";
import { MemoryRouter } from "react-router-dom";
import { cancelDeployment, DeploymentStatus } from "../../modules/deployments";
import userEvent from "@testing-library/user-event";

const baseState = {
  deployments: {
    entities: {
      [dummyDeployment.id]: dummyDeployment,
    },
    ids: [dummyDeployment.id],
    canceling: {},
  },
  environments: {
    entities: {
      [dummyEnv.id]: dummyEnv,
    },
    ids: [dummyEnv.id],
  },
  pipeds: {
    entities: {
      [dummyPiped.id]: dummyPiped,
    },
    ids: [dummyPiped.id],
  },
};

describe("DeploymentDetail", () => {
  it("shows deployment detail", () => {
    const store = createStore(baseState);
    render(
      <MemoryRouter>
        <DeploymentDetail deploymentId={dummyDeployment.id} />
      </MemoryRouter>,
      {
        store,
      }
    );

    expect(screen.getByText("SUCCESS")).toBeInTheDocument();
    expect(
      screen.getByText(dummyDeployment.applicationName)
    ).toBeInTheDocument();
    expect(screen.getByText(dummyDeployment.summary)).toBeInTheDocument();
  });

  describe("status: RUNNING", () => {
    const store = createStore({
      ...baseState,
      deployments: {
        entities: {
          [dummyDeployment.id]: {
            ...dummyDeployment,
            status: DeploymentStatus.DEPLOYMENT_RUNNING,
          },
        },
        ids: [dummyDeployment.id],
        canceling: {},
      },
    });

    beforeEach(() => {
      render(
        <MemoryRouter>
          <DeploymentDetail deploymentId={dummyDeployment.id} />
        </MemoryRouter>,
        {
          store,
        }
      );
    });

    it("shows cancel button if deployment is running", () => {
      expect(screen.getByText("RUNNING")).toBeInTheDocument();
    });

    it("dispatch cancelDeployment action if click cancel button", () => {
      userEvent.click(screen.getByRole("button", { name: "Cancel" }));

      expect(store.getActions()).toMatchObject([
        {
          type: cancelDeployment.pending.type,
          meta: {
            arg: {
              deploymentId: dummyDeployment.id,
              forceRollback: false,
              forceNoRollback: false,
            },
          },
        },
      ]);
    });
  });
});
