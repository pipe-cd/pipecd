import { ApplicationGitRepository } from "pipe/pkg/app/web/model/common_pb";

export const dummyRepo: ApplicationGitRepository.AsObject = {
  id: "debug-repo",
  remote: "git@github.com:pipe-cd/debug.git",
  branch: "master",
};
