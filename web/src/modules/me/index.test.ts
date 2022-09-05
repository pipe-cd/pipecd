import { setupServer } from "msw/node";
import { createStore } from "~~/test-utils";
import { getMeHandler } from "~/mocks/services/me";
import { dummyMe } from "~/__fixtures__/dummy-me";
import { fetchMe, meSlice, selectProjectName } from "./";

const server = setupServer(getMeHandler);

beforeAll(() => {
  server.listen();
});

afterEach(() => {
  server.resetHandlers();
});

afterAll(() => {
  server.close();
});

describe("fetchMe", () => {
  it(`creates ${fetchMe.fulfilled.type} when fetching me has been done`, async () => {
    const store = createStore();
    await store.dispatch(fetchMe());

    expect(store.getActions()).toMatchObject([
      {
        type: fetchMe.pending.type,
      },
      {
        type: fetchMe.fulfilled.type,
        payload: {
          ...dummyMe,
          isLogin: true,
        },
      },
    ]);
  });
});

describe("meSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      meSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toBeNull();
  });

  describe("fetchMe", () => {
    it(`should handle ${fetchMe.fulfilled.type}`, () => {
      expect(
        meSlice.reducer(null, {
          type: fetchMe.fulfilled.type,
          payload: {
            subject: "userName",
            avatarUrl: "avatar-url",
            projectId: "pipecd",
            isLogin: true,
          },
        })
      ).toEqual({
        subject: "userName",
        avatarUrl: "avatar-url",
        projectId: "pipecd",
        isLogin: true,
      });
    });

    it(`should handle ${fetchMe.rejected.type}`, () => {
      expect(
        meSlice.reducer(null, {
          type: fetchMe.rejected.type,
        })
      ).toEqual({
        isLogin: false,
      });
    });
  });
});

describe("selectProjectName", () => {
  it("should returns projectId", () => {
    expect(
      selectProjectName({
        me: {
          subject: "userName",
          avatarUrl: "avatar-url",
          projectId: "pipecd",
          isLogin: true,
        },
      })
    ).toBe("pipecd");
  });

  it("should returns empty string if user is not logged in", () => {
    expect(
      selectProjectName({
        me: {
          isLogin: false,
        },
      })
    ).toBe("");
  });

  it("should returns empty string if MeState is null", () => {
    expect(selectProjectName({ me: null })).toBe("");
  });
});
