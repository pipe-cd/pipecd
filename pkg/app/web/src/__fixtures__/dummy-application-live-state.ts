import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { ApplicationLiveState } from "../modules/applications-live-state";
import { dummyApplication } from "./dummy-application";
import { dummyEnv } from "./dummy-environment";

export const dummyApplicationLiveState: ApplicationLiveState = {
  applicationId: dummyApplication.id,
  envId: dummyEnv.id,
  kind: ApplicationKind.KUBERNETES,
  pipedId: "piped-1",
  version: { index: 1, timestamp: 0 },
  projectId: "project-1",
  cloudrun: {},
  lambda: {},
  terraform: {},
  kubernetes: { resourcesList: [] },
};
