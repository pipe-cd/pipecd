import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Drawer,
} from "@material-ui/core";
import { useFormik } from "formik";
import { FC, memo, useCallback, useState } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_DISCARD } from "~/constants/ui-text";
import { unwrapResult, useAppDispatch, useAppSelector } from "~/hooks/redux";
import { addApplication } from "~/modules/applications";
import { selectProjectName } from "~/modules/me";
import {
  ApplicationFormTabs,
  ApplicationFormValue,
  emptyFormValues,
  validationSchema,
} from "~/components/application-form";

export interface AddApplicationDrawerProps {
  open: boolean;
  onClose?: () => void;
  onAdded?: () => void;
}

const CONFIRM_DIALOG_TITLE = "Quit adding application?";
const CONFIRM_DIALOG_DESCRIPTION =
  "Form values inputs so far will not be saved.";

export const AddApplicationDrawer: FC<AddApplicationDrawerProps> = memo(
  function AddApplicationDrawer({
    open,
    onClose = () => null,
    onAdded = () => null,
  }) {
    const dispatch = useAppDispatch();
    const [showConfirm, setShowConfirm] = useState(false);
    const formik = useFormik<ApplicationFormValue>({
      initialValues: emptyFormValues,
      validationSchema,
      enableReinitialize: true,
      async onSubmit(values) {
        await dispatch(addApplication(values))
          .then(unwrapResult)
          .then(() => {
            onAdded();
            onClose();
            formik.resetForm();
          })
          .catch(() => undefined);
      },
    });

    const projectName = useAppSelector(selectProjectName);

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
          onClose={(_, reason) => {
            if (reason === "backdropClick" && formik.isSubmitting) return;
            handleClose();
          }}
        >
          <ApplicationFormTabs
            {...formik}
            title={`Add a new application to "${projectName}" project`}
            onClose={handleClose}
            onAdded={onAdded}
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
