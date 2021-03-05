import { Drawer } from "@material-ui/core";
import { useFormik } from "formik";
import React, { FC, memo, useCallback } from "react";
import { useDispatch, useSelector } from "react-redux";
import { ADD_PIPED_SUCCESS } from "../../constants/toast-text";
import { selectProjectName } from "../../modules/me";
import { addPiped } from "../../modules/pipeds";
import { addToast } from "../../modules/toasts";
import { AppDispatch } from "../../store";
import { PipedForm, PipedFormValues, validationSchema } from "../piped-form";

export interface AddPipedDrawerProps {
  open: boolean;
  onClose: () => void;
}

export const AddPipedDrawer: FC<AddPipedDrawerProps> = memo(
  function AddPipedDrawer({ open, onClose }) {
    const dispatch = useDispatch<AppDispatch>();
    const projectName = useSelector(selectProjectName);

    const formik = useFormik<PipedFormValues>({
      initialValues: { name: "", desc: "", envIds: [] },
      validationSchema,
      validateOnMount: true,
      async onSubmit(values) {
        await dispatch(addPiped(values)).then(() => {
          dispatch(
            addToast({ message: ADD_PIPED_SUCCESS, severity: "success" })
          );
          onClose();
        });
      },
    });

    const handleClose = useCallback(() => {
      onClose();
      formik.resetForm();
    }, [formik, onClose]);

    return (
      <Drawer anchor="right" open={open} onClose={handleClose}>
        <PipedForm
          title={`Add a new piped to "${projectName}" project`}
          {...formik}
          onClose={handleClose}
        />
      </Drawer>
    );
  }
);
