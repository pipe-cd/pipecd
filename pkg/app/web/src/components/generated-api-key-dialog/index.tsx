import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@material-ui/core";
import { FC, memo, useCallback } from "react";
import { useAppSelector, useAppDispatch } from "../../hooks/redux";
import { clearGeneratedKey } from "../../modules/api-keys";
import { TextWithCopyButton } from "../text-with-copy-button";

const DIALOG_TITLE = "Generated API Key";
const VALUE_CAPTION = "API Key";

export const GeneratedAPIKeyDialog: FC = memo(function GeneratedAPIKeyDialog() {
  const dispatch = useAppDispatch();
  const generatedKey = useAppSelector<string | null>(
    (state) => state.apiKeys.generatedKey
  );
  const open = Boolean(generatedKey);

  const handleClose = useCallback(() => {
    dispatch(clearGeneratedKey());
  }, [dispatch]);

  return (
    <Dialog open={open} fullWidth>
      <DialogTitle>{DIALOG_TITLE}</DialogTitle>
      <DialogContent>
        <Typography variant="caption">{VALUE_CAPTION}</Typography>
        {generatedKey ? (
          <TextWithCopyButton name="API Key" value={generatedKey} />
        ) : null}
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
});
