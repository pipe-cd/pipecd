import { FC, memo, useCallback } from "react";
import { IconButton, makeStyles } from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import copy from "copy-to-clipboard";
import { useDispatch } from "react-redux";
import { addToast } from "../../modules/toasts";

const useStyles = makeStyles((theme) => ({
  root: {
    display: "flex",
    height: 64,
    backgroundColor: theme.palette.background.paper,
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(0.5),
    borderRadius: theme.shape.borderRadius,
    borderWidth: 1,
    borderStyle: "solid",
  },
  input: {
    border: "none",
    fontSize: 16,
    flex: 1,
    textOverflow: "ellipsis",
    paddingLeft: theme.spacing(1),
  },
}));

export interface TextWithCopyButtonProps {
  name: string;
  label: string;
  value: string;
}

export const TextWithCopyButton: FC<TextWithCopyButtonProps> = memo(
  function TextWithCopyButton({ name, label, value }) {
    const classes = useStyles();
    const dispatch = useDispatch();
    const handleCopy = useCallback(() => {
      copy(value);
      dispatch(addToast({ message: `${name} copied to clipboard.` }));
    }, [value, name, dispatch]);
    return (
      <fieldset className={classes.root}>
        <input readOnly value={value} className={classes.input} />
        <legend>{name}</legend>
        <div>
          <IconButton
            size="small"
            style={{ marginLeft: 8 }}
            aria-label={label}
            onClick={handleCopy}
          >
            <CopyIcon style={{ fontSize: 20 }} />
          </IconButton>
        </div>
      </fieldset>
    );
  }
);
