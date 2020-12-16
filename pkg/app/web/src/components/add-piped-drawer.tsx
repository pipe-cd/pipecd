import React, { FC, memo, useCallback } from "react";
import { Drawer } from "@material-ui/core";
import { PipedForm, PipedFormValues, validationSchema } from "./piped-form";
import { useFormik } from "formik";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { selectProjectName } from "../modules/me";
import { addPiped } from "../modules/pipeds";
import { AppDispatch } from "../store";
import { addToast } from "../modules/toasts";
import { ADD_PIPED_SUCCESS } from "../constants/toast-text";

interface Props {
  open: boolean;
  onClose: () => void;
}

export const AddPipedDrawer: FC<Props> = memo(function AddPipedDrawer({
  open,
  onClose,
}) {
  const dispatch = useDispatch<AppDispatch>();
  const projectName = useSelector<AppState, string>((state) =>
    selectProjectName(state.me)
  );

  const formik = useFormik<PipedFormValues>({
    initialValues: { name: "", desc: "", envIds: [] },
    validationSchema,
    validateOnMount: true,
    async onSubmit(values) {
      await dispatch(addPiped(values)).then(() => {
        dispatch(addToast({ message: ADD_PIPED_SUCCESS, severity: "success" }));
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
});
