import {
  Box,
  Button,
  Checkbox,
  Divider,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import CheckBoxIcon from "@material-ui/icons/CheckBox";
import CheckBoxOutlineBlankIcon from "@material-ui/icons/CheckBoxOutlineBlank";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { FormikProps } from "formik";
import { FC, memo } from "react";
import * as yup from "yup";
import { useAppSelector } from "~/hooks/redux";
import {
  Environment,
  selectAllEnvs,
  selectEnvEntities,
} from "~/modules/environments";

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
  setFieldValue,
  values,
  isValid,
  isSubmitting,
}) {
  const classes = useStyles();
  const envs = useAppSelector(selectAllEnvs);
  const entities = useAppSelector(selectEnvEntities);

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
        <Autocomplete
          multiple
          id="environments"
          options={envs}
          disableCloseOnSelect
          value={
            values.envIds.map((id) => entities[id]) as Environment.AsObject[]
          }
          onChange={(_, newValue) => {
            setFieldValue(
              "envIds",
              newValue.map((env) => env.id)
            );
          }}
          getOptionLabel={(option) => option.name}
          disabled={isSubmitting}
          renderOption={(option, { selected }) => (
            <>
              <Checkbox
                icon={<CheckBoxOutlineBlankIcon fontSize="small" />}
                checkedIcon={<CheckBoxIcon fontSize="small" />}
                style={{ marginRight: 8 }}
                checked={selected}
                color="primary"
              />
              {option.name}
            </>
          )}
          renderInput={(params) => (
            <TextField
              {...params}
              variant="outlined"
              label="Environments"
              margin="dense"
              placeholder="Environments"
              fullWidth
            />
          )}
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
