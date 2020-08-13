import React, { FC, useState } from "react";
import {
  makeStyles,
  Divider,
  Typography,
  TextField,
  Button,
  Checkbox,
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

const useStyles = makeStyles((theme) => ({
  root: {
    width: 600,
  },
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
  },
}));

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
  const [name, setName] = useState<string>("");
  const [desc, setDesc] = useState<string>("");
  const [envs, setEnvs] = useState<Environment[]>([]);

  const handleSave = (): void => {
    onSubmit({ name, desc, envIds: envs.map((env) => env.id) });
  };

  return (
    <div className={classes.root}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add new piped to "${projectName}"`}</Typography>
      <Divider />
      <form className={classes.form}>
        <TextField
          label="Name"
          variant="outlined"
          margin="dense"
          onChange={(e) => setName(e.currentTarget.value)}
          value={name}
          fullWidth
          required
        />
        <TextField
          label="Description"
          variant="outlined"
          margin="dense"
          onChange={(e) => setDesc(e.currentTarget.value)}
          value={desc}
          fullWidth
          required
        />
        <Autocomplete
          multiple
          id="environments"
          options={environments}
          disableCloseOnSelect
          value={envs}
          onChange={(_, newValue) => {
            setEnvs(newValue);
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
              required
            />
          )}
        />
        <Button
          color="primary"
          type="button"
          onClick={handleSave}
          disabled={name === "" || desc === "" || envs.length === 0}
        >
          SAVE
        </Button>
        <Button onClick={onClose}>CANCEL</Button>
      </form>
    </div>
  );
};
