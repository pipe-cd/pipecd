import { Autocomplete, TextField } from "@mui/material";
import { FC, useEffect, useMemo, useState } from "react";

export const LabelAutoComplete: FC<{
  value?: string[];
  onChange?: (value: string[]) => void;
  options: string[];
}> = ({ value, onChange, options }) => {
  const [selectedLabels, setSelectedLabels] = useState<string[]>(value ?? []);
  const [userOptions, setUserOptions] = useState<string[]>([]);
  const [inputValue, setInputValue] = useState<string>("");

  const fullOptions = useMemo(() => {
    const newOptions = new Set([...options, ...userOptions]);
    return Array.from(newOptions);
  }, [options, userOptions]);

  useEffect(() => {
    if (!value) return;
    setSelectedLabels(value);
  }, [value]);

  const handleChange = (newLabels: string[]): void => {
    const labels = new Set<string>();
    newLabels.forEach((label) => {
      const labelParts = label.split(":");
      if (labelParts.length !== 2) return;
      if (labelParts[0].length === 0) return;
      if (labelParts[1].length === 0) return;
      labels.add(label);
    });
    const newValue = Array.from(labels);

    setSelectedLabels(newValue);
    onChange?.(newValue);
    setUserOptions((prev) => {
      const newOptions = new Set([...prev, ...newValue]);
      return Array.from(newOptions);
    });
  };

  return (
    <Autocomplete
      multiple
      autoHighlight
      id="labels"
      noOptionsText="No selectable labels"
      options={fullOptions}
      value={selectedLabels}
      onChange={(_, newValue) => {
        handleChange(newValue);
      }}
      freeSolo
      inputValue={inputValue}
      onInputChange={(_, value) => setInputValue(value)}
      onKeyDown={(e) => {
        if (e.key === "Enter") {
          e.preventDefault();
          setInputValue("");
        }
      }}
      renderInput={(params) => (
        <TextField
          {...params}
          variant="outlined"
          label="Labels"
          margin="dense"
          placeholder="key:value"
          fullWidth
        />
      )}
    />
  );
};

export default LabelAutoComplete;
