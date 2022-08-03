import { ApplicationGitPath } from "pipecd/web/model/common_pb";
import { createApplicationGitRepository } from "./dummy-repo";

export function createGitPathFromObject(
  o: ApplicationGitPath.AsObject
): ApplicationGitPath {
  const path = new ApplicationGitPath();
  path.setPath(o.path);
  path.setUrl(o.url);
  path.setConfigFilename(o.configFilename);
  if (o.repo) {
    path.setRepo(createApplicationGitRepository(o.repo));
  }
  return path;
}
