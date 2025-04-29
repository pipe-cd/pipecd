import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
} from "@mui/material";
import { useFormik } from "formik";
import { FC } from "react";
import { rbacResourceTypes, rbacActionTypes } from "~/modules/project";
import * as yup from "yup";

export interface AddRoleDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { name: string; policies: string }) => void;
}

// resources=(\*|application|deployment|event|piped|deploymentChain|project|apiKey|insight|,)+;\s*actions=(\*|get|list|create|update|delete|,)+
const validationRgex = new RegExp(
  "resources=(" +
    rbacResourceTypes()
      .map((v) => v.replace(/\*/, "\\*"))
      .join("|") +
    "|,)+;\\s*actions=(" +
    rbacActionTypes()
      .map((v) => v.replace(/\*/, "\\*"))
      .join("|") +
    "|,)+"
);

const validationSchema = yup.object({
  name: yup.string().min(1).required(),
  policies: yup
    .string()
    .matches(validationRgex, "Invalid policy format")
    .required(),
});

export const AddRoleDialog: FC<AddRoleDialogProps> = ({
  onClose,
  onSubmit,
  open,
}) => {
  const formik = useFormik({
    initialValues: {
      name: "",
      policies: "",
    },
    validationSchema,
    validateOnMount: true,
    onSubmit: (values, actions) => {
      onSubmit({
        name: values.name,
        policies: values.policies,
      });
      actions.resetForm();
    },
    onReset: () => {
      onClose();
    },
  });

  return (
    <Dialog open={open} onClose={onClose} fullWidth>
      <form onSubmit={formik.handleSubmit} onReset={formik.handleReset}>
        <DialogTitle>Add Role</DialogTitle>
        <DialogContent>
          <TextField
            id="name"
            name="name"
            label="Role"
            variant="outlined"
            margin="dense"
            autoFocus
            value={formik.values.name}
            onChange={formik.handleChange}
            required
          />
          <TextField
            id="policies"
            name="policies"
            label="Policies"
            placeholder="resources=RESOURCE_NAMES;actions=ACTION_NAMES"
            variant="outlined"
            margin="dense"
            value={formik.values.policies}
            onChange={formik.handleChange}
            required
            fullWidth
            multiline={true}
            rows={10}
          />
        </DialogContent>
        <DialogActions>
          <Button type="reset">Close</Button>
          <Button
            type="submit"
            color="primary"
            disabled={formik.isValid === false || formik.dirty === false}
          >
            Add
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
