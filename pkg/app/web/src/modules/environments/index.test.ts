import { createReduxStore } from "~~/test-utils";
import { server } from "~/mocks/server";
import {
  deleteEnvironmentHandler,
  listEnvironmentHandler,
} from "~/mocks/services/environment";
import { dummyEnv } from "~/__fixtures__/dummy-environment";
import { environmentsSlice, fetchEnvironments, deleteEnvironment } from ".";

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("environmentsSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      environmentsSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      entities: {},
      ids: [],
    });
  });

  describe("fetchEnvironments", () => {
    it(`should handle ${fetchEnvironments.fulfilled.type}`, () => {
      expect(
        environmentsSlice.reducer(undefined, {
          type: fetchEnvironments.fulfilled.type,
          payload: [dummyEnv],
        })
      ).toEqual({
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      });
    });
  });
});

describe("async actions", () => {
  test("fetchEnvironments", async () => {
    const store = createReduxStore();

    server.use(listEnvironmentHandler);

    await store.dispatch(fetchEnvironments());
    expect(store.getState().environments).toEqual({
      entities: { [dummyEnv.id]: dummyEnv },
      ids: [dummyEnv.id],
    });
  });

  test("deleteEnvironment", async () => {
    const store = createReduxStore({
      environments: {
        entities: { [dummyEnv.id]: dummyEnv },
        ids: [dummyEnv.id],
      },
    });

    server.use(deleteEnvironmentHandler);

    await store.dispatch(deleteEnvironment({ environmentId: dummyEnv.id }));
    expect(store.getState().environments).toEqual({ entities: {}, ids: [] });
  });
});
