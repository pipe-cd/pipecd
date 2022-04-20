global.Date.now = jest.fn(() => 0);

import { LoadingStatus } from "~/types/module";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import { Command, fetchCommand, CommandStatus } from "../commands";
import {
  deploymentsSlice,
  isDeploymentRunning,
  DeploymentStatus,
  StageStatus,
  isStageRunning,
  fetchDeploymentById,
  fetchDeployments,
  fetchMoreDeployments,
  cancelDeployment,
  updateSkippableState,
} from ".";

const initialState = {
  canceling: {},
  entities: {},
  hasMore: true,
  minUpdatedAt: 0,
  ids: [],
  status: "idle" as LoadingStatus,
  loading: {},
  cursor: "",
  skippable: {},
};

test("isDeploymentRunning", () => {
  expect(isDeploymentRunning(undefined)).toBeFalsy();

  expect(isDeploymentRunning(DeploymentStatus.DEPLOYMENT_PENDING)).toBeTruthy();
  expect(isDeploymentRunning(DeploymentStatus.DEPLOYMENT_PLANNED)).toBeTruthy();
  expect(isDeploymentRunning(DeploymentStatus.DEPLOYMENT_RUNNING)).toBeTruthy();
  expect(
    isDeploymentRunning(DeploymentStatus.DEPLOYMENT_ROLLING_BACK)
  ).toBeTruthy();

  expect(
    isDeploymentRunning(DeploymentStatus.DEPLOYMENT_CANCELLED)
  ).toBeFalsy();
  expect(isDeploymentRunning(DeploymentStatus.DEPLOYMENT_FAILURE)).toBeFalsy();
  expect(isDeploymentRunning(DeploymentStatus.DEPLOYMENT_SUCCESS)).toBeFalsy();
});

test("isStageRunning", () => {
  expect(isStageRunning(StageStatus.STAGE_CANCELLED)).toBeFalsy();
  expect(isStageRunning(StageStatus.STAGE_FAILURE)).toBeFalsy();
  expect(isStageRunning(StageStatus.STAGE_SUCCESS)).toBeFalsy();
  expect(isStageRunning(StageStatus.STAGE_NOT_STARTED_YET)).toBeTruthy();
  expect(isStageRunning(StageStatus.STAGE_RUNNING)).toBeTruthy();
});

describe("deploymentsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      deploymentsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      ...initialState,
      minUpdatedAt: -2592000,
    });
  });

  describe("fetchDeploymentById", () => {
    it(`should handle ${fetchDeploymentById.pending.type}`, () => {
      expect(
        deploymentsSlice.reducer(initialState, {
          type: fetchDeploymentById.pending.type,
          meta: {
            arg: "deployment-1",
          },
        })
      ).toEqual({
        ...initialState,
        loading: {
          "deployment-1": true,
        },
      });
    });

    it(`should handle ${fetchDeploymentById.rejected.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          {
            ...initialState,
            loading: {
              "deployment-1": true,
            },
          },
          {
            type: fetchDeploymentById.rejected.type,
            meta: {
              arg: "deployment-1",
            },
          }
        )
      ).toEqual({
        ...initialState,
        loading: {
          "deployment-1": false,
        },
      });
    });

    it(`should handle ${fetchDeploymentById.fulfilled.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          {
            ...initialState,
            loading: {
              "deployment-1": false,
            },
          },
          {
            type: fetchDeploymentById.fulfilled.type,
            meta: {
              arg: "deployment-1",
            },
            payload: dummyDeployment,
          }
        )
      ).toEqual({
        ...initialState,
        entities: { [dummyDeployment.id]: dummyDeployment },
        ids: [dummyDeployment.id],
        loading: {
          "deployment-1": false,
        },
      });
    });
  });

  describe("fetchDeployments", () => {
    it(`should handle ${fetchDeployments.pending.type}`, () => {
      expect(
        deploymentsSlice.reducer(initialState, {
          type: fetchDeployments.pending.type,
        })
      ).toEqual({
        ...initialState,
        status: "loading",
      });
    });

    it(`should handle ${fetchDeployments.rejected.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          {
            ...initialState,
            status: "loading",
          },
          {
            type: fetchDeployments.rejected.type,
          }
        )
      ).toEqual({
        ...initialState,
        status: "failed",
      });
    });

    it(`should handle ${fetchDeployments.fulfilled.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          {
            ...initialState,
            status: "loading",
          },
          {
            type: fetchDeployments.fulfilled.type,
            payload: { deployments: [dummyDeployment], cursor: "next cursor" },
          }
        )
      ).toEqual({
        ...initialState,
        entities: { [dummyDeployment.id]: dummyDeployment },
        hasMore: false,
        ids: [dummyDeployment.id],
        status: "succeeded",
        cursor: "next cursor",
      });
    });
  });

  describe("fetchMoreDeployments", () => {
    it(`should handle ${fetchMoreDeployments.pending.type}`, () => {
      expect(
        deploymentsSlice.reducer(initialState, {
          type: fetchMoreDeployments.pending.type,
        })
      ).toEqual({ ...initialState, status: "loading" });
    });

    it(`should handle ${fetchMoreDeployments.rejected.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchMoreDeployments.rejected.type,
          }
        )
      ).toEqual({ ...initialState, status: "failed" });
    });

    it(`should handle ${fetchMoreDeployments.fulfilled.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          { ...initialState, status: "loading" },
          {
            type: fetchMoreDeployments.fulfilled.type,
            payload: { deployments: [dummyDeployment], cursor: "next cursor" },
          }
        )
      ).toEqual({
        ...initialState,
        hasMore: false,
        ids: [dummyDeployment.id],
        entities: { [dummyDeployment.id]: dummyDeployment },
        status: "succeeded",
        cursor: "next cursor",
        minUpdatedAt: -2592000,
      });
    });
  });

  describe("cancelDeployment", () => {
    it(`should handle ${cancelDeployment.pending.type}`, () => {
      expect(
        deploymentsSlice.reducer(initialState, {
          type: cancelDeployment.pending.type,
          meta: {
            arg: {
              deploymentId: "deployment-1",
            },
          },
        })
      ).toEqual({
        ...initialState,
        canceling: {
          "deployment-1": true,
        },
      });
    });
  });

  describe("fetchCommand", () => {
    it(`should handle ${fetchCommand.fulfilled.type}`, () => {
      expect(
        deploymentsSlice.reducer(
          {
            ...initialState,
            canceling: {
              "deployment-1": true,
            },
          },
          {
            type: fetchCommand.fulfilled.type,
            payload: {
              deploymentId: "deployment-1",
              type: Command.Type.CANCEL_DEPLOYMENT,
              status: CommandStatus.COMMAND_SUCCEEDED,
            },
          }
        )
      ).toEqual({
        ...initialState,
        canceling: {
          "deployment-1": false,
        },
      });
    });
  });

  describe("updateSkippableState", () => {
    it(`should handle ${updateSkippableState.fulfilled.type}`, () => {
      expect(
        deploymentsSlice.reducer(initialState, {
          type: updateSkippableState.fulfilled.type,
          meta: {
            arg: {
              stageId: "stage-id",
            },
          },
        })
      ).toEqual({
        ...initialState,
        skippable: {
          "stage-id": true,
        },
      });
    });
  });
});
