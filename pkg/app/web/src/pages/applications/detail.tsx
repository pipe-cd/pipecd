import { FC, memo, useEffect } from "react";
import { useParams } from "react-router-dom";
import { ApplicationDetail } from "~/components/application-detail";
import { ApplicationStateView } from "~/components/application-state-view";
import { useAppSelector, useAppDispatch } from "~/hooks/redux";
import { useInterval } from "~/hooks/use-interval";
import { fetchApplication } from "~/modules/applications";
import {
  fetchApplicationStateById,
  selectHasError,
} from "~/modules/applications-live-state";

const FETCH_INTERVAL = 4000;

export const ApplicationDetailPage: FC = memo(function ApplicationDetailPage() {
  const dispatch = useAppDispatch();
  const params = useParams<{ applicationId: string }>();
  const applicationId = decodeURIComponent(params.applicationId);
  const [hasFetchApplicationError, hasLiveStateError] = useAppSelector<
    [boolean, boolean]
  >((state) => [
    state.applications.fetchApplicationError !== null,
    selectHasError(state.applicationLiveState, params.applicationId),
  ]);

  useEffect(() => {
    if (applicationId) {
      dispatch(fetchApplicationStateById(applicationId));
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

  useInterval(
    () => {
      if (applicationId) {
        dispatch(fetchApplicationStateById(applicationId));
      }
    },
    applicationId && hasLiveStateError === false ? FETCH_INTERVAL : null
  );

  return (
    <>
      <ApplicationDetail applicationId={applicationId} />
      <ApplicationStateView applicationId={applicationId} />
    </>
  );
});
