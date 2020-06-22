import React, { FC, memo } from "react";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  ApplicationLiveState,
  selectById,
} from "../modules/applications-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";

interface Props {
  applicationId: string;
}

export const ApplicationStateView: FC<Props> = memo(
  function ApplicationStateView({ applicationId }) {
    const liveState = useSelector<AppState, ApplicationLiveState | undefined>(
      (state) => selectById(state.applicationLiveState, applicationId)
    );

    if (!liveState) {
      return null;
    }

    switch (liveState.kind) {
      case ApplicationKind.KUBERNETES:
        return (
          <KubernetesStateView resources={liveState.kubernetes.resourcesList} />
        );
      default:
    }

    // NOTE: other resource types are not implemented.
    return null;
  }
);
