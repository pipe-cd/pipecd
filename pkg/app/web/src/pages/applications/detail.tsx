import React, { FC, memo, useEffect } from "react";
import { useDispatch } from "react-redux";
import { useParams } from "react-router";
import { ApplicationDetail } from "../../components/application-detail";
import { ApplicationStateView } from "../../components/application-state-view";
import { fetchApplication } from "../../modules/applications";
import { fetchApplicationById } from "../../modules/applications-live-state";
import { useInterval } from "../../utils/use-interval";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const dispatch = useDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplicationById(applicationId));
      dispatch(fetchApplication(applicationId));
    }
  }, [dispatch, applicationId]);

  useInterval(
    () => {
      dispatch(fetchApplicationById(applicationId));
      dispatch(fetchApplication(applicationId));
    },
    applicationId ? FETCH_INTERVAL : null
  );

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </>
  );
});
