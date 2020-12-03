import React, { FC, useRef } from "react";
import {
  makeStyles,
  Dialog,
  DialogActions,
  DialogTitle,
  DialogContent,
  Typography,
  IconButton,
  Button,
} from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import copy from "copy-to-clipboard";
import { addToast } from "../modules/toasts";
import { useDispatch } from "react-redux";

const useStyles = makeStyles(() => ({
  key: {
    wordBreak: "break-all",
  },
}));

interface Props {
  open: boolean;
  generatedKey: string | null;
  onClose: () => void;
}

const DIALOG_TITLE = "Generated API Key";
const VALUE_CAPTION = "API Key";

export const GeneratedApiKeyDialog: FC<Props> = ({
  open,
  onClose,
  generatedKey,
}) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const keyRef = useRef<HTMLDivElement>(null);

  const handleOnClickCopy = (): void => {
    if (generatedKey) {
      copy(generatedKey);
      dispatch(addToast({ message: "API Key copied to clipboard" }));
    }
  };

  return (
    <Dialog open={open}>
      <DialogTitle>{DIALOG_TITLE}</DialogTitle>
      <DialogContent>
        <Typography variant="caption">{VALUE_CAPTION}</Typography>
        <Typography variant="body2" className={classes.key} ref={keyRef}>
          {generatedKey}
          <IconButton
            size="small"
            style={{ marginLeft: 8 }}
            aria-label="Copy API Key"
            onClick={handleOnClickCopy}
          >
            <CopyIcon style={{ fontSize: 20 }} />
          </IconButton>
        </Typography>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
};
