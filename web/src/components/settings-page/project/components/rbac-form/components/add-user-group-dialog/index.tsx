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
  makeStyles,
} from "@material-ui/core";
import { useFormik } from "formik";
import { FC } from "react";
import * as yup from "yup";
import { useAppSelector } from "~/hooks/redux";

export interface AddUserGroupDialogProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (values: { ssoGroup: string; role: string }) => void;
}

const useStyles = makeStyles(() => ({
  formItem: {
    width: "50%",
  },
}));

const validationSchema = yup.object({
  ssoGroup: yup.string().min(1).required(),
  role: yup.string().required(),
});

export const AddUserGroupDialog: FC<AddUserGroupDialogProps> = ({
  onClose,
  onSubmit,
  open,
}) => {
  const classes = useStyles();
  const formik = useFormik({
    initialValues: {
      ssoGroup: "",
      role: "",
    },
    validationSchema,
    validateOnMount: true,
    onSubmit: (values, actions) => {
      onSubmit({
        ssoGroup: values.ssoGroup,
        role: values.role,
      });
      actions.resetForm();
    },
    onReset: () => {
      onClose();
    },
  });
  const roles = useAppSelector((state) => state.project.rbacRoles);

  return (
    <Dialog open={open} onClose={onClose}>
      <form onSubmit={formik.handleSubmit} onReset={formik.handleReset}>
        <DialogTitle>Add User Group</DialogTitle>
        <DialogContent>
          <TextField
            id="ssoGroup"
            name="ssoGroup"
            label="Team/Group"
            variant="outlined"
            margin="dense"
            autoFocus
            value={formik.values.ssoGroup}
            onChange={formik.handleChange}
            required
            fullWidth
          />
          <FormControl
            className={classes.formItem}
            variant="outlined"
            margin="dense"
          >
            <InputLabel id="role">Role</InputLabel>
            <Select
              id="role"
              name="role"
              value={formik.values.role}
              label="Role"
              onChange={formik.handleChange}
            >
              {roles.map((role, i) => (
                <MenuItem key={i} value={role.name}>
                  {role.name}
                </MenuItem>
              ))}
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
            Add
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
