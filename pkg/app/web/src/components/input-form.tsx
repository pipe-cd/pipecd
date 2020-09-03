import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import EditIcon from "@material-ui/icons/Edit";
import React, { FC, useState } from "react";

const useStyles = makeStyles((theme) => ({
  content: {
    display: "flex",
    alignItems: "center",
  },
  name: {
    color: theme.palette.text.secondary,
    marginRight: theme.spacing(2),
    minWidth: 120,
  },
  dialog: {
    minWidth: 400,
  },
}));

interface Props {
  name: string;
  currentValue: string | undefined;
  onSave: (value: string | undefined) => void;
  // If true, edit will start with empty value.
  isSecret?: boolean;
}

export const InputForm: FC<Props> = ({
  name,
  currentValue,
  onSave,
  isSecret = false,
}) => {
  const classes = useStyles();
  const [edit, setEdit] = useState(false);
  const initialValue = isSecret ? "" : currentValue;
  const [value, setValue] = useState(initialValue);

  const handleSave = (): void => {
    onSave(value);
    setEdit(!edit);
  };

  return (
    <div>
      <div className={classes.content}>
        <Typography variant="subtitle1" className={classes.name}>
          {name}
        </Typography>
        <Typography variant="body1">{currentValue}</Typography>
        <IconButton onClick={() => setEdit(!edit)}>
          <EditIcon fontSize="small" />
        </IconButton>
      </div>
      <Dialog
        open={edit}
        onEnter={() => setValue(initialValue)}
        PaperProps={{ className: classes.dialog }}
        onClose={() => {
          setEdit(false);
        }}
      >
        <DialogTitle>Edit {name}</DialogTitle>
        <DialogContent>
          <TextField
            value={value}
            variant="outlined"
            margin="dense"
            label={name}
            fullWidth
            onChange={(e) => setValue(e.currentTarget.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => {
              setEdit(false);
            }}
          >
            CANCEL
          </Button>
          <Button
            onClick={handleSave}
            type="submit"
            color="primary"
            disabled={Boolean(value) === false || value === currentValue}
          >
            SAVE
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
