import React, { FC, memo } from "react";
import { useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  ApplicationLiveState,
  selectById,
} from "../modules/applications-live-state";
import { KubernetesStateView } from "./kubernetes-state-view";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import { CircularProgress, makeStyles } from "@material-ui/core";

const useStyles = makeStyles(() => ({
  loading: {
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    flex: 1,
  },
}));

interface Props {
  applicationId: string;
}

export const ApplicationStateView: FC<Props> = memo(
  function ApplicationStateView({ applicationId }) {
    const classes = useStyles();
    const liveState = useSelector<AppState, ApplicationLiveState | undefined>(
      (state) => selectById(state.applicationLiveState, applicationId)
    );

    if (!liveState) {
      return (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      );
    }

    switch (liveState.kind) {
      case ApplicationKind.KUBERNETES:
        return (
          <KubernetesStateView
            resources={liveState.kubernetes?.resourcesList || []}
          />
        );
      default:
    }

    // NOTE: other resource types are not implemented.
    return null;
  }
);
