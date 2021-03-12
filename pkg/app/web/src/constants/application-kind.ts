import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";

export const APPLICATION_KIND_TEXT: Record<ApplicationKind, string> = {
  [ApplicationKind.KUBERNETES]: "KUBERNETES",
  [ApplicationKind.TERRAFORM]: "TERRAFORM",
  [ApplicationKind.CROSSPLANE]: "CROSSPLANE",
  [ApplicationKind.LAMBDA]: "LAMBDA",
  [ApplicationKind.CLOUDRUN]: "CLOUDRUN",
  [ApplicationKind.ECS]: "ECS",
};
