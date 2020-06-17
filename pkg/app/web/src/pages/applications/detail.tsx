import React, { FC, memo, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router";
import { ApplicationDetail } from "../../components/application-detail";
import { AppState } from "../../modules";
import {
  ApplicationLiveState,
  fetchApplicationById,
  selectById,
} from "../../modules/applications-live-state";
import { fetchApplications } from "../../modules/applications";
import { ApplicationStateView } from "../../components/application-state-view";

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
    }

    // TODO: Fetch only current active application data
    dispatch(fetchApplications());
  }, [applicationId]);

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
