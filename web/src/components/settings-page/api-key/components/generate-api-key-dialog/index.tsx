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
} from "@mui/material";
import { useFormik } from "formik";
import { FC } from "react";
import * as yup from "yup";
import { API_KEY_ROLE_TEXT } from "~/constants/api-key-role-text";
import { APIKey } from "pipecd/web/model/apikey_pb";

export interface GenerateAPIKeyDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; role: APIKey.Role }) => void;
}

const validationSchema = yup.object({
  name: yup.string().min(1).required(),
  role: yup.number().required(),
});

export const GenerateAPIKeyDialog: FC<GenerateAPIKeyDialogProps> = ({
  onClose,
  onSubmit,
  open,
}) => {
  const formik = useFormik({
    initialValues: {
      name: "",
      role: APIKey.Role.READ_ONLY,
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
              <MenuItem value={APIKey.Role.READ_ONLY}>
                {API_KEY_ROLE_TEXT[APIKey.Role.READ_ONLY]}
              </MenuItem>
              <MenuItem value={APIKey.Role.READ_WRITE}>
                {API_KEY_ROLE_TEXT[APIKey.Role.READ_WRITE]}
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
