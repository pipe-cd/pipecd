import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import React, { FC, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../modules";
import { Application, selectById } from "../modules/applications";
import {
  generateSealedSecret,
  SealedSecret,
  clearSealedSecret,
} from "../modules/sealed-secret";

const useStyles = makeStyles((theme) => ({
  description: {
    marginBottom: theme.spacing(2),
  },
  targetApp: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
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
      dispatch(generateSealedSecret({ data, pipedId: application.pipedId }));
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
            <TextField
              value={sealedSecret.data}
              variant="outlined"
              margin="dense"
              label="Secret Data"
              multiline
              fullWidth
              rows={6}
              required
              autoFocus
              onFocus={(e) => {
                e.target.select();
              }}
            />
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
