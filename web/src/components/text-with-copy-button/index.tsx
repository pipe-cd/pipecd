import { FC, memo } from "react";
import { CopyIconButton } from "../copy-icon-button";
import { Box } from "@mui/material";

export interface TextWithCopyButtonProps {
  name: string;
  value: string;
}

export const TextWithCopyButton: FC<TextWithCopyButtonProps> = memo(
  function TextWithCopyButton({ name, value }) {
    return (
      <Box
        component={"fieldset"}
        sx={(theme) => ({
          display: "flex",
          alignItems: "center",
          height: 64,
          backgroundColor: theme.palette.background.paper,
          marginTop: theme.spacing(1),
          marginBottom: theme.spacing(0.5),
          borderRadius: theme.shape.borderRadius,
          borderWidth: 1,
          borderStyle: "solid",
        })}
      >
        <Box
          component="input"
          readOnly
          value={value}
          sx={{
            border: "none",
            fontSize: 16,
            flex: 1,
            textOverflow: "ellipsis",
            paddingLeft: 1,
          }}
        />
        <legend>{name}</legend>
        <Box ml={2}>
          <CopyIconButton name={name} value={value} size="small" />
        </Box>
      </Box>
    );
  }
);
