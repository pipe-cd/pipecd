import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import React, { FC, useRef, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { Application, selectById } from "../modules/applications";
import {
  generateSealedSecret,
  SealedSecret,
  clearSealedSecret,
} from "../modules/sealed-secret";
import copy from "copy-to-clipboard";
import { addToast } from "../modules/toasts";

const useStyles = makeStyles((theme) => ({
  description: {
    marginBottom: theme.spacing(2),
  },
  targetApp: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  encryptedSecret: {
    wordBreak: "break-all",
  },
}));

interface Props {
  applicationId: string | null;
  open: boolean;
  onClose: () => void;
}

const DIALOG_TITLE = "Encrypting secret data for application";

export const SealedSecretDialog: FC<Props> = ({
  open,
  applicationId,
  onClose,
}) => {
  const classes = useStyles();
  const dispatch = useDispatch();
  const [data, setData] = useState("");
  const secretRef = useRef<HTMLDivElement>(null);

  const [application, sealedSecret] = useSelector<
    AppState,
    [Application | undefined, SealedSecret]
  >((state) => [
    applicationId ? selectById(state.applications, applicationId) : undefined,
    state.sealedSecret,
  ]);

  const handleGenerate = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    if (application) {
      dispatch(generateSealedSecret({ data, pipedId: application.pipedId, base64Encoding: false }));
    }
  };

  const handleOnEnter = (): void => {
    setData("");
  };

  const handleOnExited = (): void => {
    // Clear state after closed dialog
    setTimeout(() => {
      dispatch(clearSealedSecret());
    }, 200);
  };

  const handleOnClickCopy = (): void => {
    if (sealedSecret.data) {
      copy(sealedSecret.data);
      dispatch(addToast({ message: "Secret copied to clipboard" }));
    }
  };

  if (!application) {
    return null;
  }

  return (
    <Dialog
      open={open}
      onEnter={handleOnEnter}
      onExit={handleOnExited}
      onClose={onClose}
    >
      {sealedSecret.data ? (
        <>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption">Encrypted secret data</Typography>
            <Typography
              variant="body2"
              className={classes.encryptedSecret}
              ref={secretRef}
            >
              {sealedSecret.data}
              <IconButton
                size="small"
                style={{ marginLeft: 8 }}
                aria-label="Copy secret"
                onClick={handleOnClickCopy}
              >
                <CopyIcon style={{ fontSize: 20 }} />
              </IconButton>
            </Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={onClose}>Close</Button>
          </DialogActions>
        </>
      ) : (
        <form onSubmit={handleGenerate}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption">Application</Typography>
            <Typography variant="body1" className={classes.targetApp}>
              {application.name}
            </Typography>
            <TextField
              value={data}
              variant="outlined"
              margin="dense"
              label="Secret Data"
              multiline
              fullWidth
              rows={4}
              required
              autoFocus
              onChange={(e) => setData(e.currentTarget.value)}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={onClose} disabled={sealedSecret.isLoading}>
              Cancel
            </Button>
            <Button
              color="primary"
              type="submit"
              disabled={sealedSecret.isLoading}
            >
              Encrypt
            </Button>
          </DialogActions>
        </form>
      )}
    </Dialog>
  );
};
