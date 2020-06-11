import React, { memo, FC, useEffect } from "react";
import { useParams } from "react-router";
import { useDispatch, useSelector } from "react-redux";
import {
  fetchApplicationById,
  selectById,
  ApplicationLiveState,
} from "../../modules/applications";
import { ApplicationDetail } from "../../components/application-detail";
import { AppState } from "../../modules";
import { ApplicationSyncStatus } from "pipe/pkg/app/web/model/application_pb";

export const ApplicationDetailPage: FC = memo(() => {
  const dispatch = useDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const application = useSelector<AppState, ApplicationLiveState | undefined>(
    (state) => selectById(state.applications, applicationId)
  );

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplicationById(applicationId));
    }
  }, [applicationId]);

  if (!application) {
    return <div>loading</div>;
  }

  return (
    <div>
      <ApplicationDetail
        name={application.applicationId}
        env={application.envId}
        version={`${application.version.index}`}
        piped={application.pipedId}
        deployedAt={application.version.timestamp}
        deploymentId="aaa"
        status={ApplicationSyncStatus.DEPLOYING}
        description="hello"
      />
    </div>
  );
});
