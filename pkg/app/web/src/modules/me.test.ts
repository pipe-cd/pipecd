import { createStore } from "../../test-utils";
import * as meAPI from "../api/me";
import { meSlice, fetchMe, Role, selectProjectName } from "./me";

describe("fetchMe", () => {
  it(`creates ${fetchMe.fulfilled.type} when fetching me has been done`, async () => {
    const me = {
      subject: "userName",
      avatarUrl: "avatar-url",
      projectId: "pipecd",
      projectRole: Role.ProjectRole.ADMIN,
    };
    jest.spyOn(meAPI, "getMe").mockImplementation(() => Promise.resolve(me));

    const store = createStore();

    await store.dispatch(fetchMe());

    expect(store.getActions()).toMatchObject([
      {
        type: fetchMe.pending.type,
      },
      {
        type: fetchMe.fulfilled.type,
        payload: {
          ...me,
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
            projectRole: Role.ProjectRole.ADMIN,
            isLogin: true,
          },
        })
      ).toEqual({
        subject: "userName",
        avatarUrl: "avatar-url",
        projectId: "pipecd",
        projectRole: Role.ProjectRole.ADMIN,
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
        subject: "userName",
        avatarUrl: "avatar-url",
        projectId: "pipecd",
        projectRole: Role.ProjectRole.ADMIN,
        isLogin: true,
      })
    ).toBe("pipecd");
  });

  it("should returns empty string if user is not logged in", () => {
    expect(
      selectProjectName({
        isLogin: false,
      })
    ).toBe("");
  });

  it("should returns empty string if MeState is null", () => {
    expect(selectProjectName(null)).toBe("");
  });
});
