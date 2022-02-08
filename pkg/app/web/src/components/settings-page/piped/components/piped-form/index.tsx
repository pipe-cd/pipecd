import {
  Box,
  Button,
  Divider,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import { FormikProps } from "formik";
import { FC, memo } from "react";
import * as yup from "yup";

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
  },
}));

export const validationSchema = yup.object().shape({
  name: yup.string().required(),
  desc: yup.string().required(),
});

export interface PipedFormValues {
  name: string;
  desc: string;
  envIds: string[];
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
  const classes = useStyles();

  return (
    <Box width={600}>
      <Typography className={classes.title} variant="h6">
        {title}
      </Typography>
      <Divider />
      <form className={classes.form} onSubmit={handleSubmit}>
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
          onChange={handleChange}
          value={values.desc}
          fullWidth
          required
          disabled={isSubmitting}
        />
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
      </form>
    </Box>
  );
});
