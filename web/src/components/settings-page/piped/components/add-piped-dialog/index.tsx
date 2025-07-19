import { Dialog } from "@mui/material";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import { ADD_PIPED_SUCCESS } from "~/constants/toast-text";
import { unwrapResult, useAppDispatch } from "~/hooks/redux";
import { addPiped } from "~/modules/pipeds";
import { addToast } from "~/modules/toasts";
import { PipedForm, PipedFormValues, validationSchema } from "../piped-form";
import useProjectName from "~/contexts/auth-context/use-project-name";

export interface AddPipedDrawerProps {
  open: boolean;
  onClose: () => void;
}

export const AddPipedDialog: FC<AddPipedDrawerProps> = memo(
  function AddPipedDialog({ open, onClose }) {
    const dispatch = useAppDispatch();
    const projectName = useProjectName();

    const formik = useFormik<PipedFormValues>({
      initialValues: { name: "", desc: "" },
      validationSchema,
      validateOnMount: true,
      async onSubmit(values) {
        await dispatch(addPiped(values))
          .then(unwrapResult)
          .then(() => {
            dispatch(
              addToast({ message: ADD_PIPED_SUCCESS, severity: "success" })
            );
            onClose();
          })
          .catch(() => undefined);
      },
    });

    const handleClose = useCallback(() => {
      onClose();
      formik.resetForm();
    }, [formik, onClose]);

    return (
      <Dialog open={open} onClose={handleClose}>
        <PipedForm
          title={`Add a new piped to "${projectName}" project`}
          {...formik}
          onClose={handleClose}
        />
      </Dialog>
    );
  }
);
