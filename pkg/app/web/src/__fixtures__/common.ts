import { ApplicationGitPath } from "pipecd/pkg/app/web/model/common_pb";
import { createApplicationGitRepository } from "./dummy-repo";

export function createGitPathFromObject(
  o: ApplicationGitPath.AsObject
): ApplicationGitPath {
  const path = new ApplicationGitPath();
  path.setPath(o.path);
  path.setUrl(o.url);
  path.setConfigFilename(o.configFilename);
  path.setConfigPath(o.configPath);
  if (o.repo) {
    path.setRepo(createApplicationGitRepository(o.repo));
  }
  return path;
}
