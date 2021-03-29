import { setupServer } from "msw/node";
import { createReduxStore } from "../../test-utils";
import { getInsightApplicationCountHandler } from "../mocks/services/insight";
import {
  applicationCountsSlice,
  fetchApplicationCount,
} from "./application-counts";

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
    ).toMatchInlineSnapshot(`
      Object {
        "counts": Object {
          "CLOUDRUN": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "CROSSPLANE": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "ECS": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "KUBERNETES": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "LAMBDA": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
          "TERRAFORM": Object {
            "DISABLED": 0,
            "ENABLED": 0,
          },
        },
        "updatedAt": 0,
      }
    `);
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
        CROSSPLANE: {
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
