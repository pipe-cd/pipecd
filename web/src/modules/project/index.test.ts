import { projectSlice, fetchProject, updateStaticAdmin } from "./";
import { parseRBACPolicies } from "./index";
import {
  ProjectRBACPolicy,
  ProjectRBACResource,
} from "pipecd/web/model/project_pb";

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

describe("parseRBACPolicies", () => {
  it("should parse RBAC policies all", () => {
    const policies = parseRBACPolicies({ policies: "resources=*;actions=*" });
    const expected = [
      {
        resourcesList: [
          {
            type: ProjectRBACResource.ResourceType.ALL,
            labelsMap: [],
          },
        ],
        actionsList: [ProjectRBACPolicy.Action.ALL],
      },
    ];
    expect(policies[0].toObject()).toEqual(expected[0]);
  });

  it("should parse RBAC policies with resources and actions specified", () => {
    const policies = parseRBACPolicies({
      policies: "resources=application;actions=get,create",
    });
    const expected = [
      {
        resourcesList: [
          {
            type: ProjectRBACResource.ResourceType.APPLICATION,
            labelsMap: [],
          },
        ],
        actionsList: [
          ProjectRBACPolicy.Action.GET,
          ProjectRBACPolicy.Action.CREATE,
        ],
      },
    ];
    expect(policies[0].toObject()).toEqual(expected[0]);
  });

  it("should parse RBAC policies with multiple resources and actions specified", () => {
    const policies = parseRBACPolicies({
      policies: "resources=application,deployment;actions=get,create",
    });
    const expected = [
      {
        resourcesList: [
          {
            type: ProjectRBACResource.ResourceType.APPLICATION,
            labelsMap: [],
          },
          {
            type: ProjectRBACResource.ResourceType.DEPLOYMENT,
            labelsMap: [],
          },
        ],
        actionsList: [
          ProjectRBACPolicy.Action.GET,
          ProjectRBACPolicy.Action.CREATE,
        ],
      },
    ];
    expect(policies[0].toObject()).toEqual(expected[0]);
  });

  it("should parse RBAC policies with multiple policies", () => {
    const policies = parseRBACPolicies({
      policies: `resources=application;actions=get
resources=deployment;actions=get,create`,
    });
    const expected = [
      {
        resourcesList: [
          {
            type: ProjectRBACResource.ResourceType.APPLICATION,
            labelsMap: [],
          },
        ],
        actionsList: [ProjectRBACPolicy.Action.GET],
      },
      {
        resourcesList: [
          {
            type: ProjectRBACResource.ResourceType.DEPLOYMENT,
            labelsMap: [],
          },
        ],
        actionsList: [
          ProjectRBACPolicy.Action.GET,
          ProjectRBACPolicy.Action.CREATE,
        ],
      },
    ];
    expect(policies[0].toObject()).toEqual(expected[0]);
    expect(policies[1].toObject()).toEqual(expected[1]);
  });
});
