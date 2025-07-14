import {
  Box,
  Button,
  Checkbox,
  Drawer,
  FormControl,
  FormControlLabel,
  InputLabel,
  MenuItem,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import { useFormik } from "formik";
import { FC, memo, useCallback, useEffect, useMemo, useState } from "react";
import * as yup from "yup";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import { UI_TEXT_CANCEL, UI_TEXT_CLOSE } from "~/constants/ui-text";
import { useGetPipeds } from "~/queries/pipeds/use-get-pipeds";
import { useGenerateSealedSecret } from "~/queries/sealed-secret/use-generate-sealed-secret";
import { sortFunc } from "~/utils/common";

export interface EncryptSecretDrawerProps {
  open: boolean;
  onClose: () => void;
}

const DRAWER_TITLE = "Encrypting secret data";
const BASE64_CHECKBOX_LABEL =
  "Use base64 encoding before encrypting the secret";

const validationSchema = yup.object({
  secretData: yup.string().required(),
  base64: yup.bool(),
  pipedId: yup.string().required("Please select a piped"),
});

export const EncryptSecretDrawer: FC<EncryptSecretDrawerProps> = memo(
  function EncryptSecretDrawer({ open, onClose }) {
    const [sealedSecret, setSealedSecret] = useState<string | null>(null);
    const { data: pipeds = [] } = useGetPipeds(
      { withStatus: true },
      { enabled: open }
    );

    const { mutateAsync: generateSealedSecret } = useGenerateSealedSecret();

    const pipedOptions = useMemo(() => {
      return pipeds
        .filter((piped) => !piped.disabled)
        .sort((a, b) => sortFunc(a.name, b.name));
    }, [pipeds]);

    const formik = useFormik({
      initialValues: {
        secretData: "",
        base64: false,
        pipedId: "",
      },
      validationSchema,
      async onSubmit(values) {
        return generateSealedSecret(
          {
            data: values.secretData,
            pipedId: values.pipedId,
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

    const handleClose = useCallback(() => {
      onClose();
    }, [onClose]);

    const { resetForm } = formik;
    useEffect(() => {
      if (open) {
        resetForm();
      }
    }, [open, resetForm]);

    return (
      <Drawer
        anchor="right"
        open={open}
        onClose={handleClose}
        sx={{
          "& .MuiDrawer-paper": {
            width: 600,
          },
        }}
      >
        <Box sx={{ p: 3 }}>
          <Typography variant="h6" sx={{ mb: 3 }}>
            {DRAWER_TITLE}
          </Typography>

          {sealedSecret ? (
            <Box>
              <Typography variant="caption" color="textSecondary">
                Encrypted secret data
              </Typography>
              <TextWithCopyButton
                name="Encrypted secret"
                value={sealedSecret}
              />
              <Box sx={{ mt: 3, display: "flex", justifyContent: "flex-end" }}>
                <Button onClick={handleClose}>{UI_TEXT_CLOSE}</Button>
              </Box>
            </Box>
          ) : (
            <form onSubmit={formik.handleSubmit}>
              <FormControl fullWidth margin="dense">
                <InputLabel id="piped-select" required>
                  Piped
                </InputLabel>
                <Select
                  labelId="piped-select"
                  id="pipedId"
                  name="pipedId"
                  label="Piped"
                  value={formik.values.pipedId}
                  onChange={(e) => {
                    formik.setFieldValue("pipedId", e.target.value);
                  }}
                >
                  {pipedOptions.map((e) => (
                    <MenuItem key={e.id} value={e.id}>
                      {e.name} ({e.id})
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

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

              <Box sx={{ mt: 3, display: "flex", justifyContent: "flex-end" }}>
                <Button
                  onClick={onClose}
                  disabled={formik.isSubmitting}
                  sx={{ mr: 1 }}
                >
                  {UI_TEXT_CANCEL}
                </Button>
                <Button
                  color="primary"
                  type="submit"
                  variant="contained"
                  disabled={
                    formik.isSubmitting
                    // formik.isValid === false ||
                    // formik.dirty === false
                  }
                >
                  Encrypt
                </Button>
              </Box>
            </form>
          )}
        </Box>
      </Drawer>
    );
  }
);

export default EncryptSecretDrawer;
