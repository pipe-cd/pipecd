import { Drawer } from "@material-ui/core";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import { UPDATE_PIPED_SUCCESS } from "~/constants/toast-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { editPiped, fetchPipeds, selectPipedById } from "~/modules/pipeds";
import { addToast } from "~/modules/toasts";
import { PipedForm, PipedFormValues, validationSchema } from "../piped-form";

export interface EditPipedDrawerProps {
  pipedId: string | null;
  onClose: () => void;
}

export const EditPipedDrawer: FC<EditPipedDrawerProps> = memo(
  function EditPipedDrawer({ pipedId, onClose }) {
    const dispatch = useAppDispatch();
    const piped = useAppSelector(selectPipedById(pipedId));

    const formik = useFormik<PipedFormValues>({
      initialValues: {
        name: piped?.name || "",
        desc: piped?.desc || "",
        envIds: [],
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
  }
);
