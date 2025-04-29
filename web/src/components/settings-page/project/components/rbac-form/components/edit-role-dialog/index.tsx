import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  Typography,
  DialogTitle,
  TextField,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { useFormik } from "formik";
import { FC, memo } from "react";
import { useAppSelector } from "~/hooks/redux";
import {
  formalizePoliciesList,
  rbacResourceTypes,
  rbacActionTypes,
} from "~/modules/project";
import * as yup from "yup";

const useStyles = makeStyles((theme) => ({
  deleteTargetName: {
    color: theme.palette.text.primary,
    fontWeight: theme.typography.fontWeightMedium,
  },
  description: {
    marginBottom: theme.spacing(2),
  },
}));

export interface EditRoleDialogProps {
  role: string | null;
  onClose: () => void;
  onUpdate: (values: { name: string; policies: string }) => void;
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
  policies: yup
    .string()
    .matches(validationRgex, "Invalid policy format")
    .required(),
});

const DIALOG_TITLE = "Edit Role";

export const EditRoleDialog: FC<EditRoleDialogProps> = memo(
  function EditRoleDialog({ role, onUpdate, onClose }) {
    const classes = useStyles();
    const rs = useAppSelector((state) => state.project.rbacRoles);
    const r = rs.filter((r) => r.name == role)[0];

    const formik = useFormik({
      initialValues: {
        policies: formalizePoliciesList({
          policiesList: r?.policiesList || [],
        }),
      },
      enableReinitialize: true,
      validationSchema,
      onSubmit: (values, actions) => {
        onUpdate({
          name: role || "",
          policies: values.policies,
        });
        actions.resetForm();
      },
    });

    return (
      <Dialog open={Boolean(r)} onClose={onClose} fullWidth>
        <form onSubmit={formik.handleSubmit}>
          <DialogTitle>{DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            <Typography variant="caption">Role</Typography>
            <Typography variant="body1" className={classes.deleteTargetName}>
              {role}
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
