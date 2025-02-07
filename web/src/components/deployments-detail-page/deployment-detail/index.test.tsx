import userEvent from "@testing-library/user-event";
import { cancelDeployment, DeploymentStatus } from "~/modules/deployments";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { createStore, render, screen, MemoryRouter } from "~~/test-utils";
import { DeploymentDetail } from ".";

const baseState = {
  deployments: {
    entities: {
      [dummyDeployment.id]: dummyDeployment,
    },
    ids: [dummyDeployment.id],
    canceling: {},
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
