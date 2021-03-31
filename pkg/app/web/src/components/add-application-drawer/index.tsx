import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Drawer,
} from "@material-ui/core";
import { useFormik } from "formik";
import React, { FC, memo, useCallback, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { UI_TEXT_CANCEL, UI_TEXT_DISCARD } from "../../constants/ui-text";
import { addApplication, fetchApplications } from "../../modules/applications";
import { selectProjectName } from "../../modules/me";
import { AppDispatch } from "../../store";
import {
  ApplicationForm,
  ApplicationFormValue,
  emptyFormValues,
  validationSchema,
} from "../application-form";

export interface AddApplicationDrawerProps {
  open: boolean;
  onClose: () => void;
}

const CONFIRM_DIALOG_TITLE = "Quit adding application?";
const CONFIRM_DIALOG_DESCRIPTION =
  "Form values inputs so far will not be saved.";

export const AddApplicationDrawer: FC<AddApplicationDrawerProps> = memo(
  function AddApplicationDrawer({ open, onClose }) {
    const dispatch = useDispatch<AppDispatch>();
    const [showConfirm, setShowConfirm] = useState(false);
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

    const projectName = useSelector(selectProjectName);

    const handleClose = useCallback(() => {
      if (formik.dirty) {
        setShowConfirm(true);
      } else {
        onClose();
        formik.resetForm();
      }
    }, [onClose, formik]);

    return (
      <>
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
        <Dialog open={showConfirm}>
          <DialogTitle>{CONFIRM_DIALOG_TITLE}</DialogTitle>
          <DialogContent>{CONFIRM_DIALOG_DESCRIPTION}</DialogContent>
          <DialogActions>
            <Button onClick={() => setShowConfirm(false)}>
              {UI_TEXT_CANCEL}
            </Button>
            <Button
              color="primary"
              onClick={() => {
                setShowConfirm(false);
                onClose();
                formik.resetForm();
              }}
            >
              {UI_TEXT_DISCARD}
            </Button>
          </DialogActions>
        </Dialog>
      </>
    );
  }
);
