import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  IconButton,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import React, { FC, memo, useCallback, useRef } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AppState } from "../../modules";
import { Application, selectById } from "../../modules/applications";
import {
  generateSealedSecret,
  clearSealedSecret,
} from "../../modules/sealed-secret";
import copy from "copy-to-clipboard";
import { addToast } from "../../modules/toasts";
import { useFormik } from "formik";
import * as Yup from "yup";
import { AppDispatch } from "../../store";
import { UI_TEXT_CANCEL, UI_TEXT_CLOSE } from "../../constants/ui-text";

const useStyles = makeStyles((theme) => ({
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
const BASE64_CHECKBOX_LABEL =
  "Use base64 encoding before encrypting the secret";

const validationSchema = Yup.object({
  secretData: Yup.string().required(),
  base64: Yup.bool(),
});

export const SealedSecretDialog: FC<Props> = memo(function SealedSecretDialog({
  open,
  applicationId,
  onClose,
}) {
  const classes = useStyles();
  const dispatch = useDispatch<AppDispatch>();
  const secretRef = useRef<HTMLDivElement>(null);

  const [application, isLoading, sealedSecret] = useSelector<
    AppState,
    [Application.AsObject | undefined, boolean, string | null]
  >((state) => [
    applicationId ? selectById(state.applications, applicationId) : undefined,
    state.sealedSecret.isLoading,
    state.sealedSecret.data,
  ]);

  const formik = useFormik({
    initialValues: {
      secretData: "",
      base64: false,
    },
    validationSchema,
    validateOnMount: true,
    async onSubmit(values) {
      if (!application) {
        return;
      }
      await dispatch(
        generateSealedSecret({
          data: values.secretData,
          pipedId: application.pipedId,
          base64Encoding: values.base64,
        })
      );
    },
  });

  const handleOnEnter = useCallback(() => {
    formik.resetForm();
  }, [formik]);

  const handleOnExited = (): void => {
    // Clear state after closed dialog
    setTimeout(() => {
      dispatch(clearSealedSecret());
    }, 200);
  };

  const handleOnClickCopy = useCallback(() => {
    if (sealedSecret) {
      copy(sealedSecret);
      dispatch(addToast({ message: "Secret copied to clipboard" }));
    }
  }, [dispatch, sealedSecret]);

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
      {sealedSecret ? (
        <>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption" color="textSecondary">
              Encrypted secret data
            </Typography>
            <Typography
              variant="body2"
              className={classes.encryptedSecret}
              ref={secretRef}
            >
              {sealedSecret}
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
            <Button onClick={onClose}>{UI_TEXT_CLOSE}</Button>
          </DialogActions>
        </>
      ) : (
        <form onSubmit={formik.handleSubmit}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption" color="textSecondary">
              Application
            </Typography>
            <Typography variant="body1" className={classes.targetApp}>
              {application.name}
            </Typography>
            <TextField
              id="secretData"
              name="secretData"
              value={formik.values.secretData}
              variant="outlined"
              margin="dense"
              label="Secret Data"
              multiline
              fullWidth
              rows={4}
              required
              autoFocus
              onChange={formik.handleChange}
            />
            <FormControlLabel
              control={
                <Checkbox
                  color="primary"
                  checked={formik.values.base64}
                  onChange={formik.handleChange}
                  name="base64"
                />
              }
              disabled={formik.isSubmitting}
              label={BASE64_CHECKBOX_LABEL}
            />
          </DialogContent>
          <DialogActions>
            <Button onClick={onClose} disabled={isLoading}>
              {UI_TEXT_CANCEL}
            </Button>
            <Button
              color="primary"
              type="submit"
              disabled={isLoading || formik.isValid === false}
            >
              Encrypt
            </Button>
          </DialogActions>
        </form>
      )}
    </Dialog>
  );
});
