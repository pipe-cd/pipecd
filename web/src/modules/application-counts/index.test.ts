import { setupServer } from "msw/node";
import { createReduxStore, createStore } from "~~/test-utils";
import {
  getInsightApplicationCountHandler,
  getInsightApplicationCountNotFound,
} from "~/mocks/services/insight";
import { applicationCountsSlice, fetchApplicationCount } from "./";

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

describe("applicationCountsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      applicationCountsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      counts: {
        CLOUDRUN: {
          DISABLED: 0,
          ENABLED: 0,
        },
        ECS: {
          DISABLED: 0,
          ENABLED: 0,
        },
        KUBERNETES: {
          DISABLED: 0,
          ENABLED: 0,
        },
        LAMBDA: {
          DISABLED: 0,
          ENABLED: 0,
        },
        TERRAFORM: {
          DISABLED: 0,
          ENABLED: 0,
        },
      },
      updatedAt: 0,
      summary: { total: 0, enabled: 0, disabled: 0 },
    });
  });
});

test("fetchApplicationCount", async () => {
  server.use(getInsightApplicationCountHandler);
  const store = createReduxStore();

  await store.dispatch(fetchApplicationCount());
  expect(store.getState().applicationCounts).toEqual(
    expect.objectContaining({
      counts: {
        CLOUDRUN: {
          DISABLED: 0,
          ENABLED: 0,
        },
        ECS: {
          DISABLED: 0,
          ENABLED: 0,
        },
        KUBERNETES: {
          DISABLED: 8,
          ENABLED: 123,
        },
        LAMBDA: {
          DISABLED: 0,
          ENABLED: 0,
        },
        TERRAFORM: {
          DISABLED: 2,
          ENABLED: 75,
        },
      },
    })
  );
});

it("should be fulfilled if the API error is Not Found", async () => {
  server.use(getInsightApplicationCountNotFound);
  const store = createStore();

  await store.dispatch(fetchApplicationCount());
  expect(store.getActions()).toEqual([
    expect.objectContaining({
      type: fetchApplicationCount.pending.type,
    }),
    expect.objectContaining({
      type: fetchApplicationCount.fulfilled.type,
    }),
  ]);
});
