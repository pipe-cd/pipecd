import { Drawer } from "@material-ui/core";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import {
  ApplicationForm,
  ApplicationFormValue,
  emptyFormValues,
  validationSchema,
} from "~/components/application-form";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  Application,
  selectById as selectAppById,
} from "~/modules/applications";
import {
  clearUpdateTarget,
  updateApplication,
} from "~/modules/update-application";

export interface EditApplicationDrawerProps {
  onUpdated: () => void;
}

export const EditApplicationDrawer: FC<EditApplicationDrawerProps> = memo(
  function EditApplicationDrawer({ onUpdated }) {
    const dispatch = useAppDispatch();

    const app = useAppSelector<Application.AsObject | undefined>((state) =>
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
      validationSchema,
      enableReinitialize: true,
      async onSubmit(values) {
        if (!app) {
          return;
        }
        await dispatch(updateApplication({ ...values, applicationId: app.id }));
        onUpdated();
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
  }
);
