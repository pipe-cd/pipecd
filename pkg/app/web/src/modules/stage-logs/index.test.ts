import {
  stageLogsSlice,
  createActiveStageKey,
  fetchStageLog,
  StageLog,
  LogSeverity,
} from "./";

test("createActiveStageKey", () => {
  expect(
    createActiveStageKey({ deploymentId: "deployment-1", stageId: "stage-1" })
  ).toBe("deployment-1stage-1");
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
          completed: false,
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
        completed: false,
      };
      expect(
        stageLogsSlice.reducer(
          {
            "deployment-1stage-1": {
              stageId: "stage-1",
              deploymentId: "deployment-1",
              logBlocks: [],
              completed: false,
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
              completed: false,
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
          completed: true,
        },
      });
    });
  });
});
