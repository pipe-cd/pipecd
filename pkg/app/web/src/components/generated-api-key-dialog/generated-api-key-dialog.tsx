import React, { FC, useCallback, memo } from "react";
import {
  Dialog,
  DialogActions,
  DialogTitle,
  DialogContent,
  Typography,
  Button,
} from "@material-ui/core";
import { addToast } from "../../modules/toasts";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../../modules";
import { clearGeneratedKey } from "../../modules/api-keys";
import { COPY_API_KEY } from "../../constants/toast-text";
import { TextWithCopyButton } from "../text-with-copy-button";

const DIALOG_TITLE = "Generated API Key";
const VALUE_CAPTION = "API Key";

export const GeneratedAPIKeyDialog: FC = memo(function GeneratedAPIKeyDialog() {
  const dispatch = useDispatch();
  const generatedKey = useSelector<AppState, string | null>(
    (state) => state.apiKeys.generatedKey
  );
  const open = Boolean(generatedKey);

  const handleOnClickCopy = useCallback((): void => {
    dispatch(addToast({ message: COPY_API_KEY }));
  }, [dispatch]);

  const handleClose = useCallback(() => {
    dispatch(clearGeneratedKey());
  }, [dispatch]);

  return (
    <Dialog open={open} fullWidth>
      <DialogTitle>{DIALOG_TITLE}</DialogTitle>
      <DialogContent>
        <Typography variant="caption">{VALUE_CAPTION}</Typography>
        {generatedKey ? (
          <TextWithCopyButton
            label="Copy API Key"
            value={generatedKey}
            onCopy={handleOnClickCopy}
          />
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
});
