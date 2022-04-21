import {
  stageLogsSlice,
  createActiveStageKey,
  fetchStageLog,
  StageLog,
  LogSeverity,
} from "./";
import { setupServer } from "msw/node";
import { createReduxStore } from "~~/test-utils";
import { dummyDeployment } from "~/__fixtures__/dummy-deployment";
import {
  getStageLogHandler,
  getStageLogNotFoundHandler,
  getStageLogInternalErrorHandler,
} from "~/mocks/services/stage-log";
import { dummyLogBlocks } from "~/__fixtures__/dummy-stage-log";
import { StageStatus } from "../deployments";

const server = setupServer();

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

test("createActiveStageKey", () => {
  expect(
    createActiveStageKey({ deploymentId: "deployment-1", stageId: "stage-1" })
  ).toBe("deployment-1stage-1");
});

describe("async actions", () => {
  describe("fetchStageLog", () => {
    it("should store fetched stage data to state ", async () => {
      server.use(getStageLogHandler);
      const store = createReduxStore({
        deployments: {
          canceling: {},
          cursor: "",
          hasMore: false,
          minUpdatedAt: 0,
          loading: {},
          status: "idle",
          ids: [dummyDeployment.id],
          entities: { [dummyDeployment.id]: dummyDeployment },
          skippable: {},
        },
      });

      await store.dispatch(
        fetchStageLog({
          deploymentId: dummyDeployment.id,
          offsetIndex: 0,
          stageId: dummyDeployment.stagesList[0].id,
          retriedCount: 0,
        })
      );

      expect(store.getState().stageLogs).toEqual({
        [`${dummyDeployment.id}${dummyDeployment.stagesList[0].id}`]: {
          deploymentId: dummyDeployment.id,
          stageId: dummyDeployment.stagesList[0].id,
          logBlocks: dummyLogBlocks,
        },
      });
    });

    it("should return initialState if API returns NOT_FOUND when stage is running", async () => {
      server.use(getStageLogNotFoundHandler);
      const deployment = {
        ...dummyDeployment,
        stagesList: [
          {
            ...dummyDeployment.stagesList[0],
            status: StageStatus.STAGE_RUNNING,
          },
        ],
      };

      const store = createReduxStore({
        deployments: {
          canceling: {},
          cursor: "",
          hasMore: false,
          minUpdatedAt: 0,
          loading: {},
          status: "idle",
          ids: [deployment.id],
          entities: { [deployment.id]: deployment },
          skippable: {},
        },
      });

      await store.dispatch(
        fetchStageLog({
          deploymentId: deployment.id,
          offsetIndex: 0,
          stageId: deployment.stagesList[0].id,
          retriedCount: 0,
        })
      );

      expect(store.getState().stageLogs).toEqual({
        [`${deployment.id}${deployment.stagesList[0].id}`]: {
          deploymentId: deployment.id,
          stageId: deployment.stagesList[0].id,
          logBlocks: [],
        },
      });
      expect(store.getState().toasts).toEqual({ ids: [], entities: {} });
    });

    it("should add error toast if API return error code", async () => {
      server.use(getStageLogInternalErrorHandler);
      const store = createReduxStore({
        deployments: {
          canceling: {},
          cursor: "",
          hasMore: false,
          minUpdatedAt: 0,
          loading: {},
          status: "idle",
          ids: [dummyDeployment.id],
          entities: { [dummyDeployment.id]: dummyDeployment },
          skippable: {},
        },
      });

      await store.dispatch(
        fetchStageLog({
          deploymentId: dummyDeployment.id,
          offsetIndex: 0,
          stageId: dummyDeployment.stagesList[0].id,
          retriedCount: 0,
        })
      );

      expect(store.getState().stageLogs).toEqual({
        [`${dummyDeployment.id}${dummyDeployment.stagesList[0].id}`]: {
          deploymentId: dummyDeployment.id,
          stageId: dummyDeployment.stagesList[0].id,
          logBlocks: [],
        },
      });
      expect(store.getState().toasts).not.toEqual({ ids: [], entities: {} });
    });
  });
});

describe("stageLogsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      stageLogsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({});
  });

  describe("fetchStageLog", () => {
    it(`should handle ${fetchStageLog.pending.type}`, () => {
      expect(
        stageLogsSlice.reducer(undefined, {
          type: fetchStageLog.pending.type,
          meta: {
            arg: { deploymentId: "deployment-1", stageId: "stage-1" },
          },
        })
      ).toEqual({
        "deployment-1stage-1": {
          stageId: "stage-1",
          deploymentId: "deployment-1",
          logBlocks: [],
        },
      });
    });

    it(`should handle ${fetchStageLog.fulfilled.type}`, () => {
      const payload: StageLog = {
        stageId: "stage-1",
        deploymentId: "deployment-1",
        logBlocks: [
          { createdAt: 0, index: 0, log: "log", severity: LogSeverity.SUCCESS },
        ],
      };
      expect(
        stageLogsSlice.reducer(
          {
            "deployment-1stage-1": {
              stageId: "stage-1",
              deploymentId: "deployment-1",
              logBlocks: [],
            },
          },
          {
            type: fetchStageLog.fulfilled.type,
            meta: {
              arg: { deploymentId: "deployment-1", stageId: "stage-1" },
            },
            payload,
          }
        )
      ).toEqual({
        "deployment-1stage-1": payload,
      });
    });

    it(`should handle ${fetchStageLog.rejected.type}`, () => {
      expect(
        stageLogsSlice.reducer(
          {
            "deployment-1stage-1": {
              stageId: "stage-1",
              deploymentId: "deployment-1",
              logBlocks: [],
            },
          },
          {
            type: fetchStageLog.rejected.type,
            meta: {
              arg: { deploymentId: "deployment-1", stageId: "stage-1" },
            },
          }
        )
      ).toEqual({
        "deployment-1stage-1": {
          stageId: "stage-1",
          deploymentId: "deployment-1",
          logBlocks: [],
        },
      });
    });
  });
});
