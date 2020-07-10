import React, { FC, useState } from "react";
import {
  makeStyles,
  Divider,
  Typography,
  TextField,
  Button,
} from "@material-ui/core";

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
  onSubmit: (props: { name: string; desc: string }) => void;
  onClose: () => void;
}

export const AddPipedForm: FC<Props> = ({ projectName, onSubmit, onClose }) => {
  const classes = useStyles();
  const [name, setName] = useState<string>("");
  const [desc, setDesc] = useState<string>("");

  const handleSave = (): void => {
    onSubmit({ name, desc });
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
        <Button
          color="primary"
          type="button"
          onClick={handleSave}
          disabled={name === "" || desc === ""}
        >
          SAVE
        </Button>
        <Button onClick={onClose}>CANCEL</Button>
      </form>
    </div>
  );
};
