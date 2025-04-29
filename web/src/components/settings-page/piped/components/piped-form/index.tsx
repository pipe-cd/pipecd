import {
  Box,
  Button,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from "@mui/material";
import { FormikProps } from "formik";
import { FC, memo } from "react";
import * as yup from "yup";

export const validationSchema = yup.object().shape({
  name: yup.string().required(),
  desc: yup.string().required(),
});

export interface PipedFormValues {
  name: string;
  desc: string;
}

export type PipedFormProps = FormikProps<PipedFormValues> & {
  title: string;
  onClose: () => void;
};

export const PipedForm: FC<PipedFormProps> = memo(function PipedForm({
  title,
  onClose,
  handleSubmit,
  handleChange,
  values,
  isValid,
  isSubmitting,
}) {
  return (
    <Box width={600}>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{title}</DialogTitle>
        <DialogContent>
          <TextField
            id="name"
            name="name"
            label="Name"
            variant="outlined"
            margin="dense"
            onChange={handleChange}
            value={values.name}
            fullWidth
            required
            disabled={isSubmitting}
          />
          <TextField
            id="desc"
            name="desc"
            label="Description"
            variant="outlined"
            margin="dense"
            fullWidth
            required
            onChange={handleChange}
            value={values.desc}
            disabled={isSubmitting}
          />
        </DialogContent>
        <DialogActions>
          <Button
            color="primary"
            type="submit"
            disabled={isValid === false || isSubmitting}
          >
            SAVE
          </Button>
          <Button onClick={onClose} disabled={isSubmitting}>
            CANCEL
          </Button>
        </DialogActions>
      </form>
    </Box>
  );
});
