import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router";
import { ApplicationDetail } from "../../components/application-detail";
import { ApplicationStateView } from "../../components/application-state-view";
import { fetchApplication } from "../../modules/applications";
import {
  clearError,
  fetchApplicationStateById,
  selectHasError,
} from "../../modules/applications-live-state";
import { useInterval } from "../../hooks/use-interval";
import { AppState } from "../../modules";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const dispatch = useDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const hasLiveStateError = useSelector<AppState, boolean>((state) =>
    selectHasError(state.applicationLiveState, params.applicationId)
  );

  const fetchData = (): void => {
    if (applicationId) {
      if (hasLiveStateError === false) {
        dispatch(fetchApplicationStateById(applicationId));
      }
      dispatch(fetchApplication(applicationId));
    }
  };

  useEffect(fetchData, [applicationId, dispatch, hasLiveStateError]);
  useInterval(fetchData, applicationId ? FETCH_INTERVAL : null);

  useEffect(() => {
    return () => {
      dispatch(clearError(applicationId));
    };
  }, [dispatch, applicationId]);

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </>
  );
});
