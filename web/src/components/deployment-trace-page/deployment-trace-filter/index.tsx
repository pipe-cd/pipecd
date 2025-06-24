import { Box, TextField } from "@mui/material";
import { FC, useMemo, useState } from "react";
import { FilterView } from "~/components/filter-view";
import debounce from "~/utils/debounce";

type Props = {
  filterValues: { commitHash?: string };
  onClear: () => void;
  onChange: (options: { commitHash?: string }) => void;
};

const DEBOUNCE_INPUT_WAIT = 1000;

const DeploymentTraceFilter: FC<Props> = ({
  filterValues,
  onClear,
  onChange,
}) => {
  const [commitHash, setCommitHash] = useState<string | null>(
    filterValues.commitHash ?? ""
  );

  const debounceChangeCommitHash = useMemo(
    () => debounce(onChange, DEBOUNCE_INPUT_WAIT),
    [onChange]
  );

  const onChangeCommitHash = (commitHash: string): void => {
    debounceChangeCommitHash({ commitHash: commitHash });
  };

  return (
    <FilterView
      onClear={() => {
        onClear();
        setCommitHash("");
      }}
    >
      <Box
        sx={{
          width: "100%",
          marginTop: 4,
        }}
      >
        <TextField
          id="commit-hash"
          label="Commit hash"
          variant="outlined"
          fullWidth
          value={commitHash || ""}
          onChange={(e) => {
            const text = e.target.value;
            setCommitHash(text);
            onChangeCommitHash(text);
          }}
        />
      </Box>
    </FilterView>
  );
};

export default DeploymentTraceFilter;
