import { Dialog } from "@mui/material";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import { ADD_PIPED_SUCCESS } from "~/constants/toast-text";
import { PipedForm, PipedFormValues, validationSchema } from "../piped-form";
import useProjectName from "~/contexts/auth-context/use-project-name";
import { useAddPiped } from "~/queries/pipeds/use-add-piped";
import { useToast } from "~/contexts/toast-context";

export interface AddPipedDrawerProps {
  open: boolean;
  onClose: () => void;
  onSuccess: (data: { id: string; key: string }) => void;
}

export const AddPipedDialog: FC<AddPipedDrawerProps> = memo(
  function AddPipedDialog({ open, onClose, onSuccess }) {
    const { mutateAsync: addPiped } = useAddPiped();
    const { addToast } = useToast();
    const projectName = useProjectName();

    const formik = useFormik<PipedFormValues>({
      initialValues: { name: "", desc: "" },
      validationSchema,
      validateOnMount: true,
      onSubmit(values) {
        addPiped(values).then((res) => {
          addToast({ message: ADD_PIPED_SUCCESS, severity: "success" });
          formik.resetForm();
          onSuccess({ id: res.id, key: res.key });
        });
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
