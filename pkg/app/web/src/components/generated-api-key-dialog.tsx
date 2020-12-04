import React, { FC, useCallback, memo, useRef } from "react";
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
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { clearGeneratedKey } from "../modules/api-keys";

const useStyles = makeStyles(() => ({
  key: {
    wordBreak: "break-all",
  },
}));

const DIALOG_TITLE = "Generated API Key";
const VALUE_CAPTION = "API Key";
const TOAST_MESSAGE = "API Key copied to clipboard";

export const GeneratedAPIKeyDialog: FC = memo(function GeneratedAPIKeyDialog() {
  const classes = useStyles();
  const dispatch = useDispatch();
  const generatedKey = useSelector<AppState, string | null>(
    (state) => state.apiKeys.generatedKey
  );
  const open = Boolean(generatedKey);
  const keyRef = useRef<HTMLDivElement>(null);

  const handleOnClickCopy = useCallback((): void => {
    if (generatedKey) {
      copy(generatedKey);
      dispatch(addToast({ message: TOAST_MESSAGE }));
    }
  }, [generatedKey, dispatch]);

  const handleClose = useCallback(() => {
    dispatch(clearGeneratedKey());
  }, [dispatch]);

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
        <Button onClick={handleClose}>Close</Button>
      </DialogActions>
    </Dialog>
  );
});
