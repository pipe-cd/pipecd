import { Drawer } from "@material-ui/core";
import { useFormik } from "formik";
import React, { FC, memo, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { addApplication, fetchApplications } from "../modules/applications";
import { selectProjectName } from "../modules/me";
import { AppDispatch } from "../store";
import {
  ApplicationForm,
  emptyFormValues,
  validationSchema,
  ApplicationFormValue,
} from "./application-form";

interface Props {
  open: boolean;
  onClose: () => void;
}

export const AddApplicationDrawer: FC<Props> = memo(
  function AddApplicationDrawer({ open, onClose }) {
    const dispatch = useDispatch<AppDispatch>();
    const formik = useFormik<ApplicationFormValue>({
      initialValues: emptyFormValues,
      validateOnMount: true,
      validationSchema,
      enableReinitialize: true,
      async onSubmit(values) {
        await dispatch(addApplication(values));
        dispatch(fetchApplications());
        onClose();
        formik.resetForm();
      },
    });

    const projectName = useSelector<AppState, string>((state) =>
      selectProjectName(state.me)
    );

    const handleClose = useCallback(() => {
      onClose();
      formik.resetForm();
    }, [onClose, formik]);

    return (
      <Drawer
        anchor="right"
        open={open}
        onClose={handleClose}
        ModalProps={{ disableBackdropClick: formik.isSubmitting }}
      >
        <ApplicationForm
          {...formik}
          title={`Add a new application to "${projectName}" project`}
          onClose={handleClose}
        />
      </Drawer>
    );
  }
);
