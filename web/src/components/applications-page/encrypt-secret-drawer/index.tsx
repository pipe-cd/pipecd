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
import { FC, memo, useCallback, useEffect, useRef, useMemo } from "react";
import * as yup from "yup";
import { TextWithCopyButton } from "~/components/text-with-copy-button";
import { UI_TEXT_CANCEL } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import { Piped, fetchPipeds, selectAllPipeds } from "~/modules/pipeds";
import { sortFunc } from "~/utils/common";
import {
  clearSealedSecret,
  generateSealedSecret,
} from "~/modules/sealed-secret";

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
    const dispatch = useAppDispatch();

    const pipeds = useAppSelector<Piped.AsObject[]>(selectAllPipeds);
    const sealedSecret = useAppSelector<string | null>(
      (state) => state.sealedSecret.data
    );

    const pipedOptions = useMemo(() => {
      return pipeds
        .filter((piped) => !piped.disabled)
        .sort((a, b) => sortFunc(a.name, b.name));
    }, [pipeds]);

    useEffect(() => {
      if (open) {
        dispatch(fetchPipeds(false));
      }
    }, [dispatch, open]);

    const formik = useFormik({
      initialValues: {
        secretData: "",
        base64: false,
        pipedId: "",
      },
      validationSchema,
      async onSubmit(values) {
        await dispatch(
          generateSealedSecret({
            data: values.secretData,
            pipedId: values.pipedId,
            base64Encoding: values.base64,
          })
        );
      },
    });

    const handleClose = useCallback(() => {
      onClose();
      dispatch(clearSealedSecret());
    }, [dispatch, onClose]);

    const prevOpen = useRef(false);
    useEffect(() => {
      if (!prevOpen.current && open) {
        formik.resetForm();
      }
      prevOpen.current = open;
    }, [open, formik]);

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
                <Button onClick={handleClose}>Close</Button>
              </Box>
            </Box>
          ) : (
            <form onSubmit={formik.handleSubmit}>
              <FormControl
                fullWidth
                margin="dense"
                error={formik.touched.pipedId && Boolean(formik.errors.pipedId)}
              >
                <InputLabel id="piped-select">Piped</InputLabel>
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
                {formik.touched.pipedId && formik.errors.pipedId && (
                  <Typography variant="caption" color="error">
                    {formik.errors.pipedId}
                  </Typography>
                )}
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
                autoFocus
                onChange={formik.handleChange}
                error={
                  formik.touched.secretData && Boolean(formik.errors.secretData)
                }
                helperText={
                  formik.touched.secretData && formik.errors.secretData
                }
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
                    formik.isSubmitting ||
                    formik.isValid === false ||
                    formik.dirty === false
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
