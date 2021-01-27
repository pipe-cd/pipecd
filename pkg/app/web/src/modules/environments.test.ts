import {
  ListEnvironmentsResponse,
  AddEnvironmentResponse,
} from "pipe/pkg/app/web/api_client/service_pb";
import { createStore } from "../../test-utils";
import { createHandler } from "../mocks/create-handler";
import { server } from "../mocks/server";
import {
  createEnvFromObject,
  dummyEnv,
} from "../__fixtures__/dummy-environment";
import {
  environmentsSlice,
  fetchEnvironments,
  addEnvironment,
} from "./environments";

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
    const store = createStore();

    server.use(
      createHandler<ListEnvironmentsResponse>("/ListEnvironments", () => {
        const res = new ListEnvironmentsResponse();
        res.setEnvironmentsList([createEnvFromObject(dummyEnv)]);
        return res;
      })
    );

    await store.dispatch(fetchEnvironments());
    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ type: fetchEnvironments.pending.type }),
        expect.objectContaining({
          type: fetchEnvironments.fulfilled.type,
          payload: [dummyEnv],
        }),
      ])
    );
  });

  test("addEnvironment", async () => {
    const store = createStore();

    server.use(
      createHandler<AddEnvironmentResponse>("/AddEnvironment", () => {
        return new AddEnvironmentResponse();
      })
    );

    await store.dispatch(addEnvironment({ name: "env", desc: "description" }));
    expect(store.getActions()).toEqual(
      expect.arrayContaining([
        expect.objectContaining({ type: addEnvironment.pending.type }),
        expect.objectContaining({
          type: addEnvironment.fulfilled.type,
        }),
      ])
    );
  });
});
