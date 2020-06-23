import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";

export const APPLICATION_KIND_TEXT: Record<ApplicationKind, string> = {
  [ApplicationKind.KUBERNETES]: "Kubernetes",
  [ApplicationKind.TERRAFORM]: "Terraform",
  [ApplicationKind.CROSSPLANE]: "Crossplane",
  [ApplicationKind.LAMBDA]: "Lambda",
  [ApplicationKind.CLOUDRUN]: "Cloud Run",
};
