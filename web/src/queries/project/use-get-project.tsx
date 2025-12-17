import {
  useQuery,
  UseQueryOptions,
  UseQueryResult,
} from "@tanstack/react-query";
import * as projectAPI from "~/api/project";
import {
  ProjectUserGroup,
  ProjectRBACRole,
  ProjectSSOConfig,
  ProjectRBACConfig,
} from "~~/model/project_pb";

export type GitHubSSO = ProjectSSOConfig.GitHub.AsObject;
export type Teams = ProjectRBACConfig.AsObject;

type ProjectDetail = {
  id: string | null;
  desc: string | null;
  username: string | null;
  teams: Teams | null;
  sharedSSO: string | null;
  staticAdminDisabled: boolean;
  github: GitHubSSO | null;
  userGroups: ProjectUserGroup.AsObject[] | [];
  rbacRoles: ProjectRBACRole.AsObject[] | [];
  disabled: boolean;
};

export const useGetProject = (
  queryOption: UseQueryOptions<ProjectDetail> = {}
): UseQueryResult<ProjectDetail> => {
  return useQuery({
    queryKey: ["project", "detail"],
    queryFn: async () => {
      const { project } = await projectAPI.getProject();

      if (!project) {
        return {
          id: null,
          desc: null,
          staticAdminDisabled: false,
          username: null,
          teams: null,
          github: null,
          sharedSSO: null,
          userGroups: [],
          rbacRoles: [],
          disabled: false,
        };
      }

      return {
        id: project.id,
        desc: project.desc,
        staticAdminDisabled: project.staticAdminDisabled,
        username: project.staticAdmin?.username || "",
        teams: project.rbac ?? null,
        github: project.sso?.github ?? null,
        sharedSSO: project.sharedSsoName,
        userGroups: project.userGroupsList,
        rbacRoles: project.rbacRolesList,
        disabled: project.disabled,
      };
    },
    placeholderData: {
      id: null,
      desc: null,
      staticAdminDisabled: false,
      username: null,
      teams: null,
      github: null,
      sharedSSO: null,
      userGroups: [],
      rbacRoles: [],
      disabled: false,
    },
    retry: false,
    refetchOnMount: false,
    refetchOnReconnect: false,
    refetchOnWindowFocus: false,
    staleTime: 120000, // 2 minutes
    cacheTime: 300000, // 5 minutes
    ...queryOption,
  });
};
