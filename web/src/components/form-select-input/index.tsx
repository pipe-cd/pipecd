import { FormControl, InputLabel, MenuItem, Select } from "@mui/material";
import React, { useEffect } from "react";

type BaseOption = {
  label?: string;
  value: string;
  [key: string]: unknown;
};

type Props<T> = {
  id: string;
  label?: string;
  value?: string;
  options?: T[];
  required?: boolean;
  onChange?: (value: string, option: T) => void;
  getOptionLabel?: (option: T) => React.ReactNode;
  disabled?: boolean;
  defaultValue?: string;
};

const FormSelectInput = <T extends BaseOption>({
  id,
  label = "",
  value,
  options = [],
  required = false,
  onChange,
  disabled = false,
  getOptionLabel = (option: T) => option.label ?? option.value,
  defaultValue = "",
}: Props<T>): JSX.Element => {
  const [internalValue, setInternalValue] = React.useState<string>(
    defaultValue
  );

  useEffect(() => {
    if (value !== undefined) {
      setInternalValue(value);
    }
  }, [value]);

  return (
    <FormControl fullWidth variant="outlined" required={required}>
      {label && <InputLabel id={id}>{label}</InputLabel>}
      <Select
        labelId={id}
        id={id}
        label={label}
        value={internalValue}
        fullWidth
        onChange={(e) => {
          const inputValue = e.target.value as string;
          const item = options.find((e) => e.value === inputValue);
          if (!item) return;

          if (onChange) onChange?.(inputValue, item);
          setInternalValue(inputValue);
        }}
        disabled={disabled}
      >
        {options.map((op) => (
          <MenuItem
            value={String(op.value)}
            key={String(op.value)}
            disabled={!!op?.disabled}
          >
            {getOptionLabel(op)}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};
export default FormSelectInput;
