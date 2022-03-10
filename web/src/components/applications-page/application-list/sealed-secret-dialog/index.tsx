import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import { useFormik } from "formik";
import { FC, memo, useCallback } from "react";
import * as yup from "yup";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import { UI_TEXT_CANCEL, UI_TEXT_CLOSE } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { Application, selectById } from "~/modules/applications";
import {
  clearSealedSecret,
  generateSealedSecret,
} from "~/modules/sealed-secret";

const useStyles = makeStyles((theme) => ({
  targetApp: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  secretInput: {
    border: "none",
    fontSize: 14,
    flex: 1,
    textOverflow: "ellipsis",
  },
  encryptedSecret: {
    wordBreak: "break-all",
  },
}));

export interface SealedSecretDialogProps {
  applicationId: string | null;
  open: boolean;
  onClose: () => void;
}

const DIALOG_TITLE = "Encrypting secret data for application";
const BASE64_CHECKBOX_LABEL =
  "Use base64 encoding before encrypting the secret";

const validationSchema = yup.object({
  secretData: yup.string().required(),
  base64: yup.bool(),
});

export const SealedSecretDialog: FC<SealedSecretDialogProps> = memo(
  function SealedSecretDialog({ open, applicationId, onClose }) {
    const classes = useStyles();
    const dispatch = useAppDispatch();

    const application = useAppSelector<Application.AsObject | undefined>(
      (state) =>
        applicationId
          ? selectById(state.applications, applicationId)
          : undefined
    );
    const sealedSecret = useAppSelector<string | null>(
      (state) => state.sealedSecret.data
    );

    const formik = useFormik({
      initialValues: {
        secretData: "",
        base64: false,
      },
      validationSchema,
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

    const handleClose = useCallback(() => {
      onClose();
      dispatch(clearSealedSecret());
    }, [dispatch, onClose]);

    if (!application) {
      return null;
    }

    return (
      <Dialog open={open} onEnter={handleOnEnter} onClose={handleClose}>
        {sealedSecret ? (
          <>
            <DialogTitle>{DIALOG_TITLE}</DialogTitle>
            <DialogContent>
              <Typography variant="caption" color="textSecondary">
                Encrypted secret data
              </Typography>
              <TextWithCopyButton
                name="Encrypted secret"
                value={sealedSecret}
              />
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose}>{UI_TEXT_CLOSE}</Button>
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
              <Button onClick={onClose} disabled={formik.isSubmitting}>
                {UI_TEXT_CANCEL}
              </Button>
              <Button
                color="primary"
                type="submit"
                disabled={
                  formik.isSubmitting ||
                  formik.isValid === false ||
                  formik.dirty === false
                }
              >
                Encrypt
              </Button>
            </DialogActions>
          </form>
        )}
      </Dialog>
    );
  }
);
