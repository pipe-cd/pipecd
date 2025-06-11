import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  Typography,
  DialogTitle,
  TextField,
} from "@mui/material";
import { useFormik } from "formik";
import { FC, memo } from "react";
import * as yup from "yup";
import { formalizePoliciesList } from "~/utils/formalize-policies-list";
import { POLICIES_STRING_REGEX } from "~/constants/project";
import { ProjectRBACRole } from "pipecd/web/model/project_pb";

export interface EditRoleDialogProps {
  role: ProjectRBACRole.AsObject | null;
  onClose: () => void;
  onUpdate: (values: { name: string; policies: string }) => void;
}

const validationSchema = yup.object({
  policies: yup
    .string()
    .matches(POLICIES_STRING_REGEX, "Invalid policy format")
    .required(),
});

const DIALOG_TITLE = "Edit Role";

export const EditRoleDialog: FC<EditRoleDialogProps> = memo(
  function EditRoleDialog({ role, onUpdate, onClose }) {
    const formik = useFormik({
      initialValues: {
        policies: formalizePoliciesList({
          policiesList: role?.policiesList || [],
        }),
      },
      enableReinitialize: true,
      validationSchema,
      onSubmit: (values, actions) => {
        onUpdate({
          name: role?.name || "",
          policies: values.policies,
        });
        actions.resetForm();
      },
    });

    return (
      <Dialog open={Boolean(role)} onClose={onClose} fullWidth>
        <form onSubmit={formik.handleSubmit}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption">Role</Typography>
            <Typography
              variant="body1"
              sx={(theme) => ({
                color: theme.palette.text.primary,
                fontWeight: theme.typography.fontWeightMedium,
              })}
            >
              {role?.name}
            </Typography>
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
            <Button type="reset" onClick={onClose}>
              Cancel
            </Button>
            <Button
              type="submit"
              color="primary"
              disabled={formik.isValid === false || formik.dirty === false}
            >
              Edit
            </Button>
          </DialogActions>
        </form>
      </Dialog>
    );
  }
);
