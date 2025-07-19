import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Typography,
} from "@mui/material";
import { FC, memo } from "react";
import { TextWithCopyButton } from "~/components/text-with-copy-button";

const DIALOG_TITLE = "Generated API Key";
const VALUE_CAPTION = "API Key";

type Props = {
  generatedKey: string | null;
  onClose: () => void;
};

export const GeneratedAPIKeyDialog: FC<Props> = memo(
  function GeneratedAPIKeyDialog({ generatedKey, onClose }) {
    return (
      <Dialog open={Boolean(generatedKey)} fullWidth>
        <DialogTitle>{DIALOG_TITLE}</DialogTitle>
        <DialogContent>
          <Typography variant="caption">{VALUE_CAPTION}</Typography>
          {generatedKey ? (
            <TextWithCopyButton name="API Key" value={generatedKey} />
          ) : null}
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Close</Button>
        </DialogActions>
      </Dialog>
    );
  }
);
