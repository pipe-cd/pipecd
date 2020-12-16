import React, { FC, memo, useCallback } from "react";
import { Drawer } from "@material-ui/core";
import { PipedForm, PipedFormValues, validationSchema } from "./piped-form";
import { useFormik } from "formik";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { editPiped, fetchPipeds, Piped, selectById } from "../modules/pipeds";
import { AppDispatch } from "../store";
import { addToast } from "../modules/toasts";
import { UPDATE_PIPED_SUCCESS } from "../constants/toast-text";

interface Props {
  pipedId: string | null;
  onClose: () => void;
}

export const EditPipedDrawer: FC<Props> = memo(function EditPipedDrawer({
  pipedId,
  onClose,
}) {
  const dispatch = useDispatch<AppDispatch>();
  const piped = useSelector<AppState, Piped | undefined>((state) =>
    pipedId ? selectById(state.pipeds, pipedId) : undefined
  );

  const formik = useFormik<PipedFormValues>({
    initialValues: {
      name: piped?.name || "",
      desc: piped?.desc || "",
      envIds: piped?.envIdsList || [],
    },
    enableReinitialize: true,
    validationSchema,
    async onSubmit({ desc, envIds, name }) {
      if (!pipedId) {
        return;
      }

      await dispatch(editPiped({ pipedId, name, desc, envIds })).then(() => {
        dispatch(fetchPipeds(true));
        dispatch(
          addToast({ message: UPDATE_PIPED_SUCCESS, severity: "success" })
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
    <Drawer anchor="right" open={Boolean(piped)} onClose={handleClose}>
      <PipedForm
        title={`Edit piped "${piped?.name}"`}
        {...formik}
        onClose={handleClose}
      />
    </Drawer>
  );
});
