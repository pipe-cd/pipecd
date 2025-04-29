import makeStyles from "@mui/styles/makeStyles";
import { FC, memo } from "react";
import { CopyIconButton } from "../copy-icon-button";

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
  copyButton: {
    marginLeft: theme.spacing(2),
  },
}));

export interface TextWithCopyButtonProps {
  name: string;
  value: string;
}

export const TextWithCopyButton: FC<TextWithCopyButtonProps> = memo(
  function TextWithCopyButton({ name, value }) {
    const classes = useStyles();
    return (
      <fieldset className={classes.root}>
        <input readOnly value={value} className={classes.input} />
        <legend>{name}</legend>
        <div>
          <CopyIconButton
            name={name}
            value={value}
            size="small"
            className={classes.copyButton}
          />
        </div>
      </fieldset>
    );
  }
);
