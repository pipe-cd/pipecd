import { Drawer } from "@material-ui/core";
import { FC, memo, useCallback, useState } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_DISCARD } from "~/constants/ui-text";
import { useAppSelector } from "~/hooks/redux";
import { selectProjectName } from "~/modules/me";
import { ApplicationFormTabs } from "~/components/application-form";
import DialogConfirm from "~/components/dialog-confirm";

export interface AddApplicationDrawerProps {
  open: boolean;
  onClose?: () => void;
  onAdded?: () => void;
}

const CONFIRM_DIALOG_TITLE = "Quit adding application?";
const CONFIRM_DIALOG_DESCRIPTION =
  "Form values inputs so far will not be saved.";

const AddApplicationDrawer: FC<AddApplicationDrawerProps> = memo(
  function AddApplicationDrawer({
    open,
    onClose = () => null,
    onAdded = () => null,
  }) {
    const [showConfirm, setShowConfirm] = useState(false);
    const [isFormDirty, setIsFormDirty] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const projectName = useAppSelector(selectProjectName);

    const handleClose = useCallback(() => {
      if (isFormDirty) {
        setShowConfirm(true);
      } else {
        onClose();
      }
    }, [isFormDirty, onClose]);

    return (
      <>
        <Drawer
          anchor="right"
          open={open}
          onClose={() => {
            if (isSubmitting) return;
            handleClose();
          }}
        >
          <ApplicationFormTabs
            title={`Add a new application to "${projectName}" project`}
            onClose={handleClose}
            onFinished={onAdded}
            setIsFormDirty={setIsFormDirty}
            setIsSubmitting={setIsSubmitting}
          />
        </Drawer>

        <DialogConfirm
          open={showConfirm}
          title={CONFIRM_DIALOG_TITLE}
          description={CONFIRM_DIALOG_DESCRIPTION}
          onCancel={() => setShowConfirm(false)}
          cancelText={UI_TEXT_CANCEL}
          confirmText={UI_TEXT_DISCARD}
          onConfirm={() => {
            onClose();
            setShowConfirm(false);
          }}
        />
      </>
    );
  }
);

export default AddApplicationDrawer;
