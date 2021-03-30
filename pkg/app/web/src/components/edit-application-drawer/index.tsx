import React, { FC, memo, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../../modules";
import {
  clearUpdateTarget,
  updateApplication,
} from "../../modules/update-application";
import {
  Application,
  fetchApplications,
  selectById as selectAppById,
} from "../../modules/applications";
import {
  ApplicationForm,
  validationSchema,
  ApplicationFormValue,
  emptyFormValues,
} from "../application-form";
import { AppDispatch } from "../../store";
import { useFormik } from "formik";
import { Drawer } from "@material-ui/core";

export const EditApplicationDrawer: FC = memo(function EditApplicationDrawer() {
  const dispatch = useDispatch<AppDispatch>();

  const app = useSelector<AppState, Application.AsObject | undefined>((state) =>
    state.updateApplication.targetId
      ? selectAppById(state.applications, state.updateApplication.targetId)
      : undefined
  );

  const formik = useFormik<ApplicationFormValue>({
    initialValues: app
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
      : emptyFormValues,
    validateOnMount: true,
    validationSchema,
    enableReinitialize: true,
    async onSubmit(values) {
      if (!app) {
        return;
      }
      await dispatch(updateApplication({ ...values, applicationId: app.id }));
      dispatch(fetchApplications());
    },
  });

  const handleClose = useCallback(() => {
    dispatch(clearUpdateTarget());
  }, [dispatch]);

  return (
    <Drawer
      anchor="right"
      open={Boolean(app)}
      onClose={handleClose}
      ModalProps={{ disableBackdropClick: formik.isSubmitting }}
    >
      <ApplicationForm
        {...formik}
        title={`Edit "${app?.name}"`}
        onClose={handleClose}
        disableGitPath
      />
    </Drawer>
  );
});
