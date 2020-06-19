import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router";
import { ApplicationDetail } from "../../components/application-detail";
import { ApplicationStateView } from "../../components/application-state-view";
import { AppState } from "../../modules";
import { fetchApplication } from "../../modules/applications";
import {
  ApplicationLiveState,
  fetchApplicationById,
  selectById,
} from "../../modules/applications-live-state";
import { useInterval } from "../../utils/use-interval";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(() => {
  const dispatch = useDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const application = useSelector<AppState, ApplicationLiveState | undefined>(
    (state) => selectById(state.applicationLiveState, applicationId)
  );

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplicationById(applicationId));
      dispatch(fetchApplication(applicationId));
    }
  }, [applicationId]);

  useInterval(
    () => {
      dispatch(fetchApplicationById(applicationId));
      dispatch(fetchApplication(applicationId));
    },
    applicationId ? FETCH_INTERVAL : null
  );

  if (!application) {
    return <div>loading</div>;
  }

  return (
    <div>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </div>
  );
});
