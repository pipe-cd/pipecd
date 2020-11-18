import React, { FC } from "react";
import {
  makeStyles,
  Divider,
  Typography,
  TextField,
  Button,
  Checkbox,
  Box,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import { AppState } from "../modules";
import {
  Environment,
  selectAll as selectAllEnvs,
} from "../modules/environments";
import { useSelector } from "react-redux";
import CheckBoxOutlineBlankIcon from "@material-ui/icons/CheckBoxOutlineBlank";
import CheckBoxIcon from "@material-ui/icons/CheckBox";
import { useFormik } from "formik";
import * as Yup from "yup";

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
  },
}));

const validationSchema = Yup.object().shape({
  name: Yup.string().required(),
  desc: Yup.string().required(),
  envs: Yup.array().required(),
});

interface Props {
  projectName: string;
  onSubmit: (props: { name: string; desc: string; envIds: string[] }) => void;
  onClose: () => void;
}

export const AddPipedForm: FC<Props> = ({ projectName, onSubmit, onClose }) => {
  const classes = useStyles();
  const environments = useSelector<AppState, Environment[]>((state) =>
    selectAllEnvs(state.environments)
  );

  const formik = useFormik<{ name: string; desc: string; envs: Environment[] }>(
    {
      initialValues: {
        name: "",
        desc: "",
        envs: [],
      },
      validationSchema,
      validateOnMount: true,
      onSubmit: (values) => {
        onSubmit({
          name: values.name,
          desc: values.desc,
          envIds: values.envs.map((env) => env.id),
        });
      },
    }
  );

  return (
    <Box width={600}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add a new piped to "${projectName}" project`}</Typography>
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
        <Autocomplete
          multiple
          id="environments"
          options={environments}
          disableCloseOnSelect
          value={formik.values.envs}
          onChange={(_, newValue) => {
            formik.setFieldValue("envs", newValue);
          }}
          getOptionLabel={(option) => option.name}
          renderOption={(option, { selected }) => (
            <React.Fragment>
              <Checkbox
                icon={<CheckBoxOutlineBlankIcon fontSize="small" />}
                checkedIcon={<CheckBoxIcon fontSize="small" />}
                style={{ marginRight: 8 }}
                checked={selected}
                color="primary"
              />
              {option.name}
            </React.Fragment>
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
          disabled={formik.isValid === false}
        >
          SAVE
        </Button>
        <Button onClick={onClose}>CANCEL</Button>
      </form>
    </Box>
  );
};
