import React from "react";
import { ApplicationKind } from "../modules/applications";
import { ApplicationKindIcon } from "./application-kind-icon";

export default {
  title: "ApplicationKindIcon",
  component: ApplicationKindIcon,
};

export const overview: React.FC = () => (
  <ApplicationKindIcon kind={ApplicationKind.KUBERNETES} />
);
