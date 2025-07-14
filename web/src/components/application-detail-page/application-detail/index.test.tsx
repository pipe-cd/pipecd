import userEvent from "@testing-library/user-event";
import { server } from "~/mocks/server";
import * as applicationsAPI from "~/api/applications";
import { dummyApplication } from "~/__fixtures__/dummy-application";
import { dummyPiped } from "~/__fixtures__/dummy-piped";
import { render, screen, waitFor, MemoryRouter, act } from "~~/test-utils";
import { ApplicationDetail } from ".";
import { dummyLiveStates } from "~/__fixtures__/dummy-application-live-state";
import { Application, ApplicationSyncStatus } from "~/types/applications";
import { SyncStrategy } from "~~/model/common_pb";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

const deployingApp: Application.AsObject = {
  ...dummyApplication,
  deploying: true,
  id: "deploying-app",
  syncState: {
    status: ApplicationSyncStatus.DEPLOYING,
    headDeploymentId: "",
    reason: "",
    shortReason: "",
    timestamp: 0,
  },
};

describe("ApplicationDetail", () => {
  it("shows application detail and live state", async () => {
    render(
      <MemoryRouter>
        <ApplicationDetail
          app={dummyApplication}
          hasError={false}
          liveStateLoading={false}
          liveState={dummyLiveStates[dummyApplication.kind]}
          refetchApp={() => {}}
        />
      </MemoryRouter>
    );

    expect(screen.getByText(dummyApplication.name)).toBeInTheDocument();
    await waitFor(() => {
      expect(screen.getByText(dummyPiped.name)).toBeInTheDocument();
    });
    expect(screen.getByText("Healthy")).toBeInTheDocument();
    expect(screen.getByText("Synced")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /sync$/i })).toBeInTheDocument();
  });

  it("shows application sync state as deploying if application is deploying", () => {
    render(
      <MemoryRouter>
        <ApplicationDetail
          app={deployingApp}
          hasError={false}
          liveStateLoading={false}
          liveState={dummyLiveStates[deployingApp.kind]}
          refetchApp={() => {}}
        />
      </MemoryRouter>
    );

    expect(screen.getByText("Deploying")).toBeInTheDocument();
  });

  describe("sync", () => {
    it("dispatch sync action if click sync button", async () => {
      jest.spyOn(applicationsAPI, "syncApplication");
      render(
        <MemoryRouter>
          <ApplicationDetail
            app={dummyApplication}
            hasError={false}
            liveStateLoading={false}
            liveState={dummyLiveStates[dummyApplication.kind]}
            refetchApp={() => {}}
          />
        </MemoryRouter>
      );

      await act(async () => {
        await userEvent.click(screen.getByRole("button", { name: /sync$/i }));
      });
      await waitFor(() => {
        expect(applicationsAPI.syncApplication).toHaveBeenCalledWith({
          applicationId: dummyApplication.id,
          syncStrategy: SyncStrategy.AUTO,
        });
      });
    });

    it("dispatch sync action with selected sync strategy if changed strategy and click the sync button", async () => {
      jest.spyOn(applicationsAPI, "syncApplication");
      render(
        <MemoryRouter>
          <ApplicationDetail
            app={dummyApplication}
            hasError={false}
            liveStateLoading={false}
            liveState={dummyLiveStates[dummyApplication.kind]}
            refetchApp={() => {}}
          />
        </MemoryRouter>
      );
      act(() => {
        userEvent.click(
          screen.getByRole("button", { name: /select sync strategy/i })
        );
      });
      act(() => {
        userEvent.click(
          screen.getByRole("menuitem", { name: /pipeline sync/i })
        );
      });
      await act(async () => {
        await userEvent.click(
          screen.getByRole("button", { name: /pipeline sync/i })
        );
      });
      await waitFor(() => {
        expect(applicationsAPI.syncApplication).toHaveBeenCalledWith({
          applicationId: dummyApplication.id,
          syncStrategy: SyncStrategy.PIPELINE,
        });
      });
    });
  });
});
