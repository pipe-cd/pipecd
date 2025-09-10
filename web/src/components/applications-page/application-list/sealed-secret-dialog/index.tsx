import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  TextField,
  Typography,
} from "@mui/material";
import { useFormik } from "formik";
import { FC, memo, useCallback, useState } from "react";
import * as yup from "yup";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import { UI_TEXT_CANCEL, UI_TEXT_CLOSE } from "~/constants/ui-text";
import { Application } from "~/types/applications";
import { useGenerateSealedSecret } from "~/queries/sealed-secret/use-generate-sealed-secret";

export interface SealedSecretDialogProps {
  application?: Application.AsObject | null;
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
  function SealedSecretDialog({ open, application, onClose }) {
    const [sealedSecret, setSealedSecret] = useState<string | null>(null);

    const { mutateAsync: generateSealedSecret } = useGenerateSealedSecret();

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
        return generateSealedSecret(
          {
            data: values.secretData,
            pipedId: application.pipedId,
            base64Encoding: values.base64,
          },
          {
            onSuccess: (data) => {
              setSealedSecret(data);
            },
          }
        );
      },
    });

    const handleOnEnter = useCallback(() => {
      formik.resetForm();
    }, [formik]);

    const handleClose = useCallback(() => {
      onClose();
    }, [onClose]);

    if (!application) {
      return null;
    }

    return (
      <Dialog
        open={open}
        TransitionProps={{
          onEnter: handleOnEnter,
        }}
        onClose={handleClose}
      >
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
              <Typography
                variant="body1"
                sx={(theme) => ({
                  color: theme.palette.text.primary,
                  fontWeight: theme.typography.fontWeightMedium,
                })}
              >
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
                minRows={4}
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
