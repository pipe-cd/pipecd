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
  onSubmit: (description: string) => void;
  onClose: () => void;
}

export const AddPipedForm: FC<Props> = ({ projectName, onSubmit, onClose }) => {
  const classes = useStyles();
  const [description, setDescription] = useState<string>("");

  function handleSave() {
    onSubmit(description);
  }

  return (
    <div className={classes.root}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add new piped to "${projectName}"`}</Typography>
      <Divider />
      <form className={classes.form}>
        <TextField
          label="description"
          variant="outlined"
          margin="dense"
          onChange={(e) => setDescription(e.currentTarget.value)}
          value={description}
          fullWidth
        />
        <Button
          color="primary"
          type="button"
          onClick={handleSave}
          disabled={description === ""}
        >
          SAVE
        </Button>
        <Button onClick={onClose}>CANCEL</Button>
      </form>
    </div>
  );
};
