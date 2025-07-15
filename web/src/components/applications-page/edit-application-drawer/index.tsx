import { Drawer } from "@mui/material";
import { FC, useCallback, useMemo, useState } from "react";
import ApplicationFormManualV0 from "~/components/application-form/application-form-manual-v0";
import ApplicationFormManualV1 from "~/components/application-form/application-form-manual-v1";
import DialogConfirm from "~/components/dialog-confirm";
import { UI_TEXT_CANCEL, UI_TEXT_DISCARD } from "~/constants/ui-text";
import { Application } from "~/types/applications";

type Props = {
  open: boolean;
  application?: Application.AsObject;
  onUpdated: () => void;
  onClose: () => void;
};

enum PipedVersion {
  V0 = "v0",
  V1 = "v1",
}

const CONFIRM_DIALOG_TITLE = "Quit editing application?";
const CONFIRM_DIALOG_DESCRIPTION =
  "Form values inputs so far will not be saved.";

const EditApplicationDrawer: FC<Props> = ({
  onUpdated,
  application: app,
  open,
  onClose,
}) => {
  const [showConfirm, setShowConfirm] = useState(false);
  const [isFormDirty, setIsFormDirty] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleClose = useCallback(() => {
    if (isFormDirty) {
      setShowConfirm(true);
    } else {
      onClose();
    }
  }, [isFormDirty, onClose]);

  const pipedVersion = useMemo(() => {
    if (!app) return PipedVersion.V0;

    if (!app.platformProvider) return PipedVersion.V1;

    return PipedVersion.V0;
  }, [app]);

  const editProps = useMemo(
    () => ({
      title: `Edit "${app?.name}"`,
      onClose: handleClose,
      onFinished: onUpdated,
      setIsFormDirty: setIsFormDirty,
      setIsSubmitting: setIsSubmitting,
      detailApp: app,
    }),
    [app, handleClose, onUpdated]
  );

  return (
    <Drawer
      anchor="right"
      open={Boolean(app) && open}
      onClose={() => {
        if (isSubmitting) return;
        handleClose();
      }}
    >
      {pipedVersion === PipedVersion.V0 && (
        <ApplicationFormManualV0 {...editProps} />
      )}
      {pipedVersion === PipedVersion.V1 && (
        <ApplicationFormManualV1 {...editProps} />
      )}
      <DialogConfirm
        open={showConfirm}
        title={CONFIRM_DIALOG_TITLE}
        description={CONFIRM_DIALOG_DESCRIPTION}
        onCancel={() => setShowConfirm(false)}
        cancelText={UI_TEXT_CANCEL}
        confirmText={UI_TEXT_DISCARD}
        onConfirm={() => {
          setShowConfirm(false);
          onClose();
        }}
      />
    </Drawer>
  );
};

export default EditApplicationDrawer;
