import { FC, memo } from "react";
import {
  makeStyles,
  Divider,
  Typography,
  TextField,
  Button,
  Box,
} from "@material-ui/core";
import { useFormik } from "formik";
import * as yup from "yup";

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
  },
}));

const validationSchema = yup.object({
  name: yup.string().required("Required"),
  desc: yup.string().required("Required"),
});

export interface AddEnvFormProps {
  projectName: string;
  onSubmit: (props: { name: string; desc: string }) => void;
  onCancel: () => void;
}

export const AddEnvForm: FC<AddEnvFormProps> = memo(function AddEnvForm({
  projectName,
  onSubmit,
  onCancel,
}) {
  const classes = useStyles();
  const formik = useFormik({
    initialValues: {
      name: "",
      desc: "",
    },
    validationSchema,
    onSubmit: (values, actions) => {
      onSubmit(values);
      actions.resetForm();
    },
  });

  const handleCancel = (): void => {
    onCancel();
    formik.resetForm();
  };

  return (
    <Box width={600}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add a new environment to "${projectName}" project`}</Typography>
      <Divider />
      <form className={classes.form} onSubmit={formik.handleSubmit}>
        <TextField
          id="name"
          name="name"
          label="Name"
          variant="outlined"
          margin="dense"
          onChange={formik.handleChange}
          value={formik.values.name}
          fullWidth
          required
        />
        <TextField
          id="desc"
          name="desc"
          label="Description"
          variant="outlined"
          margin="dense"
          onChange={formik.handleChange}
          value={formik.values.desc}
          fullWidth
          required
        />
        <Button
          color="primary"
          type="submit"
          disabled={formik.isValid === false}
        >
          SAVE
        </Button>
        <Button onClick={handleCancel}>CANCEL</Button>
      </form>
    </Box>
  );
});
