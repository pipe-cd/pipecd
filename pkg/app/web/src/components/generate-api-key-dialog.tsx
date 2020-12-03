import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  TextField,
} from "@material-ui/core";
import { useFormik } from "formik";
import React, { FC } from "react";
import * as Yup from "yup";
import { API_KEY_ROLE_TEXT } from "../constants/api-key-role-text";
import { APIKeyModel } from "../modules/api-keys";
interface Props {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; role: APIKeyModel.Role }) => void;
}

const validationSchema = Yup.object({
  name: Yup.string().min(1).required(),
  role: Yup.number().required(),
});

export const GenerateAPIKeyDialog: FC<Props> = ({
  onClose,
  onSubmit,
  open,
}) => {
  const formik = useFormik({
    initialValues: {
      name: "",
      role: APIKeyModel.Role.READ_ONLY,
    },
    validationSchema,
    validateOnMount: true,
    onSubmit: (values, actions) => {
      onSubmit({
        name: values.name,
        role: values.role,
      });
      actions.resetForm();
    },
    onReset: () => {
      onClose();
    },
  });

  return (
    <Dialog open={open} onClose={onClose}>
      <form onSubmit={formik.handleSubmit} onReset={formik.handleReset}>
        <DialogTitle>Generate API Key</DialogTitle>
        <DialogContent>
          <TextField
            id="name"
            name="name"
            label="Name"
            variant="outlined"
            margin="dense"
            autoFocus
            value={formik.values.name}
            onChange={formik.handleChange}
            required
            fullWidth
          />
          <FormControl variant="outlined" margin="dense">
            <InputLabel id="role">Role</InputLabel>
            <Select
              id="role"
              name="role"
              value={formik.values.role}
              label="Role"
              onChange={formik.handleChange}
            >
              <MenuItem value={APIKeyModel.Role.READ_ONLY}>
                {API_KEY_ROLE_TEXT[APIKeyModel.Role.READ_ONLY]}
              </MenuItem>
              <MenuItem value={APIKeyModel.Role.READ_WRITE}>
                {API_KEY_ROLE_TEXT[APIKeyModel.Role.READ_WRITE]}
              </MenuItem>
            </Select>
          </FormControl>
        </DialogContent>
        <DialogActions>
          <Button type="reset">Close</Button>
          <Button
            type="submit"
            color="primary"
            disabled={formik.isValid === false || formik.dirty === false}
          >
            Generate
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
