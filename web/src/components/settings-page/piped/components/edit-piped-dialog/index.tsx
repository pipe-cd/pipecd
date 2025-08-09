import { Dialog } from "@mui/material";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import { UPDATE_PIPED_SUCCESS } from "~/constants/toast-text";
import { PipedForm, PipedFormValues, validationSchema } from "../piped-form";
import { Piped } from "pipecd/web/model/piped_pb";
import { useToast } from "~/contexts/toast-context";
import { useEditPiped } from "~/queries/pipeds/use-edit-piped";
export interface EditPipedDrawerProps {
  piped: Piped.AsObject | null;
  onClose: () => void;
}

export const EditPipedDialog: FC<EditPipedDrawerProps> = memo(
  function EditPipedDialog({ piped, onClose }) {
    const { addToast } = useToast();
    const { mutateAsync: editPiped } = useEditPiped();
    const formik = useFormik<PipedFormValues>({
      initialValues: {
        name: piped?.name || "",
        desc: piped?.desc || "",
      },
      enableReinitialize: true,
      validationSchema,
      onSubmit({ desc, name }) {
        if (!piped) {
          return;
        }

        editPiped({ pipedId: piped.id, name, desc }).then(() => {
          addToast({ message: UPDATE_PIPED_SUCCESS, severity: "success" });
          onClose();
        });
      },
    });

    const handleClose = useCallback(() => {
      onClose();
      formik.resetForm();
    }, [formik, onClose]);

    return (
      <Dialog open={Boolean(piped)} onClose={handleClose}>
        <PipedForm
          title={`Edit piped "${piped?.name}"`}
          {...formik}
          onClose={handleClose}
        />
      </Dialog>
    );
  }
);
