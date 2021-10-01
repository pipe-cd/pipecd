import { FC, memo, useEffect } from "react";
import { useParams } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import { fetchApplication } from "~/modules/applications";
import { ApplicationDetail } from "./application-detail";
import { ApplicationStateView } from "./application-state-view";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const dispatch = useAppDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const [hasFetchApplicationError] = useAppSelector<
    [boolean]
  >((state) => [
    state.applications.fetchApplicationError !== null,
  ]);

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplication(applicationId));
    }
  }, [applicationId, dispatch]);

  useInterval(
    () => {
      if (applicationId) {
        dispatch(fetchApplication(applicationId));
      }
    },
    applicationId && hasFetchApplicationError === false ? FETCH_INTERVAL : null
  );

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </>
  );
});
