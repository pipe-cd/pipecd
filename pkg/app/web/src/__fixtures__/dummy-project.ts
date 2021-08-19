import {
  Project,
  ProjectRBACConfig,
  ProjectStaticUser,
} from "pipe/pkg/app/web/model/project_pb";
import {
  createRandTimes,
  randomKeyHash,
  randomUUID,
  randomWords,
} from "./utils";

const [createdAt, updatedAt] = createRandTimes(2);

export const dummyProject: Project.AsObject = {
  id: randomUUID(),
  desc: randomWords(8),
  sharedSsoName: "shared-sso",
  createdAt: createdAt.unix(),
  updatedAt: updatedAt.unix(),
  staticAdminDisabled: false,
  rbac: {
    admin: "admin-team",
    editor: "editor-team",
    viewer: "viewer-team",
    viewerRoleAsDefault: false,
  },
  staticAdmin: {
    username: "static-admin-user",
    passwordHash: randomKeyHash(),
  },
};

export function createProjectFromObject(o: Project.AsObject): Project {
  const project = new Project();
  project.setId(o.id);
  project.setDesc(o.desc);
  project.setSharedSsoName(o.sharedSsoName);
  project.setCreatedAt(o.createdAt);
  project.setUpdatedAt(o.updatedAt);
  project.setStaticAdminDisabled(o.staticAdminDisabled);
  if (o.rbac) {
    const rbac = new ProjectRBACConfig();
    rbac.setAdmin(o.rbac.admin);
    rbac.setEditor(o.rbac.editor);
    rbac.setViewer(o.rbac.viewer);
    rbac.setViewerRoleAsDefault(o.rbac.viewerRoleAsDefault);
    project.setRbac(rbac);
  }
  if (o.staticAdmin) {
    const user = new ProjectStaticUser();
    user.setUsername(o.staticAdmin.username);
    user.setPasswordHash(o.staticAdmin.passwordHash);
    project.setStaticAdmin(user);
  }
  return project;
}
