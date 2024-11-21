import { ApplicationKind } from "~/modules/applications";

export const APPLICATION_KIND_TEXT: Record<ApplicationKind, string> = {
  [ApplicationKind.KUBERNETES]: "KUBERNETES",
  [ApplicationKind.TERRAFORM]: "TERRAFORM",
  [ApplicationKind.LAMBDA]: "LAMBDA",
  [ApplicationKind.CLOUDRUN]: "CLOUDRUN",
  [ApplicationKind.ECS]: "ECS",
  [ApplicationKind.APPLICATION]: "APPLICATION",
};

export const APPLICATION_KIND_BY_NAME: Record<string, ApplicationKind> = {
  [APPLICATION_KIND_TEXT[ApplicationKind.KUBERNETES]]:
    ApplicationKind.KUBERNETES,
  [APPLICATION_KIND_TEXT[ApplicationKind.TERRAFORM]]: ApplicationKind.TERRAFORM,
  [APPLICATION_KIND_TEXT[ApplicationKind.LAMBDA]]: ApplicationKind.LAMBDA,
  [APPLICATION_KIND_TEXT[ApplicationKind.CLOUDRUN]]: ApplicationKind.CLOUDRUN,
  [APPLICATION_KIND_TEXT[ApplicationKind.ECS]]: ApplicationKind.ECS,
  [APPLICATION_KIND_TEXT[ApplicationKind.APPLICATION]]:
    ApplicationKind.APPLICATION,
};
