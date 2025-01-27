import {
  FormControl,
  InputLabel,
  makeStyles,
  MenuItem,
  Select,
} from "@material-ui/core";
import React, { useEffect } from "react";

const useStyles = makeStyles(() => ({
  formItem: {
    width: "100%",
  },
  select: {
    width: "100%",
  },
}));

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

  const classes = useStyles();

  return (
    <FormControl
      className={classes.formItem}
      variant="outlined"
      required={required}
    >
      {label && <InputLabel id={id}>{label}</InputLabel>}
      <Select
        labelId={id}
        id={id}
        label={label}
        value={internalValue}
        className={classes.select}
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
          <MenuItem value={String(op.value)} key={String(op.value)}>
            {getOptionLabel(op)}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};
export default FormSelectInput;
