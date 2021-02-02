import React from "react";
import { SyncStateReason } from "./";

export default {
  title: "APPLICATION/SyncStateReason",
  component: SyncStateReason,
};

export const overview: React.FC = () => (
  <SyncStateReason
    summary="There are 1 missing manifests and 2 redundant manifests."
    detail={`The following 1 manifests are defined in Git, but NOT appearing in the cluster:
    - apiVersion=v1, kind=Service, namespace=default, name=wait-approvalThe following 2 manifests are NOT defined in Git, but appearing in the cluster:
    - apiVersion=apps/v1, kind=Deployment, namespace=default, name=wait-approval-canary- apiVersion=v1, kind=Service, namespace=default, name=wait-approval`}
  />
);

export const diff: React.FC = () => (
  <SyncStateReason
    summary="Summary message"
    detail={`message\n\n+ added-line\n- deleted-line`}
  />
);
