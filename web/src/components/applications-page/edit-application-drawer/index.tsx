import { Drawer } from "@material-ui/core";
import { FC, useCallback, useState } from "react";
import ApplicationFormManual from "~/components/application-form/application-form-manual";
import DialogConfirm from "~/components/dialog-confirm";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { UI_TEXT_CANCEL, UI_TEXT_DISCARD } from "~/constants/ui-text";
import {
  Application,
  selectById as selectAppById,
} from "~/modules/applications";
import { clearUpdateTarget } from "~/modules/update-application";

type Props = {
  onUpdated: () => void;
};

const CONFIRM_DIALOG_TITLE = "Quit editing application?";
const CONFIRM_DIALOG_DESCRIPTION =
  "Form values inputs so far will not be saved.";

const EditApplicationDrawer: FC<Props> = ({ onUpdated }) => {
  const [showConfirm, setShowConfirm] = useState(false);
  const [isFormDirty, setIsFormDirty] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const dispatch = useAppDispatch();

  const app = useAppSelector<Application.AsObject | undefined>((state) =>
    state.updateApplication.targetId
      ? selectAppById(state.applications, state.updateApplication.targetId)
      : undefined
  );

  const handleClose = useCallback(() => {
    if (isFormDirty) {
      setShowConfirm(true);
    } else {
      dispatch(clearUpdateTarget());
    }
  }, [dispatch, isFormDirty]);

  return (
    <Drawer
      anchor="right"
      open={Boolean(app)}
      onClose={() => {
        if (isSubmitting) return;
        handleClose();
      }}
    >
      <ApplicationFormManual
        title={`Edit "${app?.name}"`}
        onClose={handleClose}
        disableApplicationInfo
        onFinished={onUpdated}
        setIsFormDirty={setIsFormDirty}
        setIsSubmitting={setIsSubmitting}
        detailApp={app}
      />
      <DialogConfirm
        open={showConfirm}
        title={CONFIRM_DIALOG_TITLE}
        description={CONFIRM_DIALOG_DESCRIPTION}
        onCancel={() => setShowConfirm(false)}
        cancelText={UI_TEXT_CANCEL}
        confirmText={UI_TEXT_DISCARD}
        onConfirm={() => {
          setShowConfirm(false);
          dispatch(clearUpdateTarget());
        }}
      />
    </Drawer>
  );
};

export default EditApplicationDrawer;
