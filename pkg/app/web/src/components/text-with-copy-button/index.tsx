import { FC, useCallback } from "react";
import { Box, IconButton, makeStyles } from "@material-ui/core";
import CopyIcon from "@material-ui/icons/FileCopyOutlined";
import copy from "copy-to-clipboard";

const useStyles = makeStyles((theme) => ({
  root: { backgroundColor: theme.palette.background.paper },
  input: { border: "none", fontSize: 14, flex: 1, textOverflow: "ellipsis" },
}));

export interface TextWithCopyButtonProps {
  label: string;
  value: string;
  onCopy: () => void;
}

export const TextWithCopyButton: FC<TextWithCopyButtonProps> = ({
  label,
  value,
  onCopy,
}) => {
  const classes = useStyles();
  const handleCopy = useCallback(() => {
    copy(value);
    onCopy();
  }, [value, onCopy]);
  return (
    <Box
      display="flex"
      p={1}
      border={1}
      borderColor="divider"
      className={classes.root}
    >
      <input readOnly value={value} className={classes.input} />
      <IconButton
        size="small"
        style={{ marginLeft: 8 }}
        aria-label={label}
        onClick={handleCopy}
      >
        <CopyIcon style={{ fontSize: 20 }} />
      </IconButton>
    </Box>
  );
};
