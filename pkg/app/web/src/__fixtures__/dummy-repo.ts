import { ApplicationGitRepository } from "pipecd/pkg/app/web/model/common_pb";

export function createApplicationGitRepository(
  o: ApplicationGitRepository.AsObject
): ApplicationGitRepository {
  const repo = new ApplicationGitRepository();
  repo.setId(o.id);
  repo.setBranch(o.branch);
  repo.setRemote(o.remote);
  return repo;
}

export const dummyRepo: ApplicationGitRepository.AsObject = {
  id: "debug-repo",
  remote: "git@github.com:pipe-cd/debug.git",
  branch: "master",
};
