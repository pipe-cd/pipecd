import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router";
import { ApplicationDetail } from "../../components/application-detail";
import { ApplicationStateView } from "../../components/application-state-view";
import { fetchApplication } from "../../modules/applications";
import { fetchApplicationStateById } from "../../modules/applications-live-state";
import { useInterval } from "../../hooks/use-interval";
import { addToast } from "../../modules/toasts";
import { AppDispatch } from "../../store";
import { AppState } from "../../modules";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const dispatch = useDispatch<AppDispatch>();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const hasLiveStateError = useSelector<AppState, boolean>(
    (state) => state.applicationLiveState.hasError
  );

  const fetchData = (): void => {
    if (applicationId) {
      if (hasLiveStateError === false) {
        dispatch(fetchApplicationStateById(applicationId)).then((result) => {
          if (fetchApplicationStateById.rejected.match(result)) {
            dispatch(
              addToast({
                message: "Failed to get application live state",
                severity: "error",
              })
            );
          }
        });
      }
      dispatch(fetchApplication(applicationId));
    }
  };

  useEffect(fetchData, [applicationId, dispatch, hasLiveStateError]);
  useInterval(fetchData, applicationId ? FETCH_INTERVAL : null);

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </>
  );
});
