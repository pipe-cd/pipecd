import React, { FC, memo, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import {
  clearUpdateTarget,
  updateApplication,
} from "../modules/update-application";
import {
  Application,
  fetchApplications,
  selectById as selectAppById,
} from "../modules/applications";
import {
  ApplicationFormDrawer,
  ApplicationFormValue,
} from "./application-form-drawer";
import { AppDispatch } from "../store";

export const EditApplicationDrawer: FC = memo(function EditApplicationDrawer() {
  const dispatch = useDispatch<AppDispatch>();
  const [applicationId, isUpdating] = useSelector<
    AppState,
    [string | null, boolean]
  >((state) => [
    state.updateApplication.targetId,
    state.updateApplication.updating,
  ]);

  const app = useSelector<AppState, Application | undefined>((state) =>
    applicationId ? selectAppById(state.applications, applicationId) : undefined
  );

  const handleClose = useCallback(() => {
    dispatch(clearUpdateTarget());
  }, [dispatch]);

  const handleSubmit = useCallback(
    (values: ApplicationFormValue) => {
      if (app) {
        dispatch(updateApplication({ ...values, applicationId: app.id })).then(
          () => {
            dispatch(fetchApplications());
          }
        );
      }
    },
    [dispatch, app]
  );

  return (
    <ApplicationFormDrawer
      open={Boolean(app)}
      title={`Edit "${app?.name}"`}
      onSubmit={handleSubmit}
      isProcessing={isUpdating}
      onClose={handleClose}
      initialFormValues={
        app
          ? {
              name: app.name,
              env: app.envId,
              kind: app.kind,
              pipedId: app.pipedId,
              repoPath: app.gitPath?.path || "",
              repo: app.gitPath?.repo || { id: "", remote: "", branch: "" },
              configFilename: app.gitPath?.configFilename || "",
              cloudProvider: app.cloudProvider,
            }
          : undefined
      }
    />
  );
});
