import { IconButton } from "@material-ui/core";
import copy from "copy-to-clipboard";
import { FileCopyOutlined as CopyIcon } from "@material-ui/icons";
import { FC, useCallback, memo } from "react";
import { useDispatch } from "react-redux";
import { addToast } from "../../modules/toasts";

export interface CopyIconButtonProps {
  name: string;
  value: string;
  size?: "small" | "medium";
  className?: string;
}

export const CopyIconButton: FC<CopyIconButtonProps> = memo(
  function CopyIconButton({ name, value, className, size }) {
    const dispatch = useDispatch();
    const handleCopy = useCallback(() => {
      copy(value);
      dispatch(addToast({ message: `${name} copied to clipboard.` }));
    }, [dispatch, value, name]);

    return (
      <IconButton
        className={className}
        aria-label={`Copy ${name}`}
        onClick={handleCopy}
        size={size}
      >
        <CopyIcon fontSize={size === "small" ? "small" : "default"} />
      </IconButton>
    );
  }
);
