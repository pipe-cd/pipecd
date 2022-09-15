import { projectSlice, fetchProject, updateStaticAdmin } from "./";

describe("projectSlice reducer", () => {
  it("should return the initial state", () => {
    expect(
      projectSlice.reducer(undefined, {
        type: "TEST_ACTION",
      })
    ).toEqual({
      desc: null,
      id: null,
      isUpdatingGitHubSSO: false,
      isUpdatingStaticAdmin: false,
      sharedSSO: null,
      staticAdminDisabled: false,
      username: null,
      userGroups: [],
      rbacRoles: [],
    });
  });

  it(`should handle ${fetchProject.fulfilled.type}`, () => {
    expect(
      projectSlice.reducer(undefined, {
        type: fetchProject.fulfilled.type,
        payload: {
          id: "id",
          desc: "desc",
          username: "User",
          teams: {
            admin: "admin",
            editor: "editor",
            viewer: "viewer",
          },
          sharedSSO: "shared-sso",
          staticAdminDisabled: false,
          github: {
            clientId: "clientId",
            clientSecret: "clientSecret",
            baseUrl: "base-url",
            uploadUrl: "upload-url",
          },
          userGroups: [
            {
              ssoGroup: "team-a",
              role: "Admin",
            },
          ],
          rbacRoles: [
            {
              name: "Admin",
              isBuiltin: true,
              policiesList: [
                {
                  resourcesList: [
                    {
                      type: 0,
                    },
                  ],
                  actionsList: [
                    {
                      action: 0,
                    },
                  ],
                },
              ],
            },
          ],
        },
      })
    ).toEqual({
      desc: "desc",
      github: {
        baseUrl: "base-url",
        clientId: "clientId",
        clientSecret: "clientSecret",
        uploadUrl: "upload-url",
      },
      id: "id",
      isUpdatingGitHubSSO: false,
      isUpdatingStaticAdmin: false,
      sharedSSO: "shared-sso",
      staticAdminDisabled: false,
      teams: {
        admin: "admin",
        editor: "editor",
        viewer: "viewer",
      },
      username: "User",
      userGroups: [
        {
          ssoGroup: "team-a",
          role: "Admin",
        },
      ],
      rbacRoles: [
        {
          name: "Admin",
          isBuiltin: true,
          policiesList: [
            {
              resourcesList: [
                {
                  type: 0,
                },
              ],
              actionsList: [
                {
                  action: 0,
                },
              ],
            },
          ],
        },
      ],
    });
  });

  it(`should handle ${updateStaticAdmin.pending.type}`, () => {
    expect(
      projectSlice.reducer(undefined, {
        type: updateStaticAdmin.pending.type,
      })
    ).toEqual({
      desc: null,
      id: null,
      isUpdatingGitHubSSO: false,
      isUpdatingStaticAdmin: true,
      sharedSSO: null,
      staticAdminDisabled: false,
      username: null,
      userGroups: [],
      rbacRoles: [],
    });
  });

  it(`should handle ${updateStaticAdmin.fulfilled.type}`, () => {
    expect(
      projectSlice.reducer(
        {
          desc: null,
          id: null,
          isUpdatingGitHubSSO: false,
          isUpdatingStaticAdmin: true,
          sharedSSO: null,
          staticAdminDisabled: false,
          username: null,
          userGroups: [],
          rbacRoles: [],
        },
        {
          type: updateStaticAdmin.fulfilled.type,
        }
      )
    ).toEqual({
      desc: null,
      id: null,
      isUpdatingGitHubSSO: false,
      isUpdatingStaticAdmin: false,
      sharedSSO: null,
      staticAdminDisabled: false,
      username: null,
      userGroups: [],
      rbacRoles: [],
    });
  });

  it(`should handle ${updateStaticAdmin.rejected.type}`, () => {
    expect(
      projectSlice.reducer(
        {
          desc: null,
          id: null,
          isUpdatingGitHubSSO: false,
          isUpdatingStaticAdmin: true,
          sharedSSO: null,
          staticAdminDisabled: false,
          username: null,
          userGroups: [],
          rbacRoles: [],
        },
        {
          type: updateStaticAdmin.rejected.type,
        }
      )
    ).toEqual({
      desc: null,
      id: null,
      isUpdatingGitHubSSO: false,
      isUpdatingStaticAdmin: false,
      sharedSSO: null,
      staticAdminDisabled: false,
      username: null,
      userGroups: [],
      rbacRoles: [],
    });
  });
});
