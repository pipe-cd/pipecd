import userEvent from "@testing-library/user-event";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { render, screen, MemoryRouter, waitFor, act } from "~~/test-utils";
import { DeploymentDetail } from ".";
import { DeploymentStatus } from "~~/model/deployment_pb";
import * as deploymentsApi from "~/api/deployments";
import { server } from "~/mocks/server";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("DeploymentDetail", () => {
  it("shows deployment detail", async () => {
    render(
      <MemoryRouter>
        <DeploymentDetail
          deploymentId={dummyDeployment.id}
          deployment={dummyDeployment}
        />
      </MemoryRouter>
    );

    await waitFor(() => {
      expect(screen.getByText("SUCCESS")).toBeInTheDocument();
    });
    expect(
      screen.getByText(dummyDeployment.applicationName)
    ).toBeInTheDocument();
    expect(screen.getByText(dummyDeployment.summary)).toBeInTheDocument();
  });

  describe("status: RUNNING", () => {
    beforeEach(() => {
      render(
        <MemoryRouter>
          <DeploymentDetail
            deploymentId={dummyDeployment.id}
            deployment={{
              ...dummyDeployment,
              status: DeploymentStatus.DEPLOYMENT_RUNNING,
            }}
          />
        </MemoryRouter>
      );
    });
    it("shows cancel button if deployment is running", async () => {
      await waitFor(() =>
        expect(screen.getByText("RUNNING")).toBeInTheDocument()
      );
    });

    it("calls cancelDeployment when cancel button is clicked", async () => {
      jest.spyOn(deploymentsApi, "cancelDeployment");

      await waitFor(() =>
        expect(screen.getByText("RUNNING")).toBeInTheDocument()
      );

      const cancelButton = screen.getByRole("button", { name: "Cancel" });
      expect(cancelButton).toBeInTheDocument();

      await act(async () => {
        userEvent.click(cancelButton);
      });

      expect(deploymentsApi.cancelDeployment).toHaveBeenCalledWith({
        deploymentId: dummyDeployment.id,
        forceRollback: false,
        forceNoRollback: false,
      });
    });
  });
});
