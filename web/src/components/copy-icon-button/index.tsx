import { IconButton } from "@mui/material";
import { FileCopyOutlined as CopyIcon } from "@mui/icons-material";
import { FC, useCallback, memo } from "react";
import { useToast } from "~/contexts/toast-context";

export interface CopyIconButtonProps {
  name: string;
  value: string;
  size?: "small" | "medium";
  className?: string;
}

export const CopyIconButton: FC<CopyIconButtonProps> = memo(
  function CopyIconButton({ name, value, className, size }) {
    const { addToast } = useToast();
    const handleCopy = useCallback(() => {
      navigator.clipboard.writeText(value).then(() => {
        addToast({ message: `${name} copied to clipboard.` });
      });
    }, [value, addToast, name]);

    return (
      <IconButton
        className={className}
        aria-label={`Copy ${name}`}
        onClick={handleCopy}
        size={size}
      >
        <CopyIcon fontSize={size === "small" ? "small" : "medium"} />
      </IconButton>
    );
  }
);
