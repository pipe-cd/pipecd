import { makeStyles, TextField } from "@material-ui/core";
import { FC, useMemo, useState } from "react";
import { FilterView } from "~/components/filter-view";
import debounce from "~/utils/debounce";

type Props = {
  filterValues: { commitHash?: string };
  onClear: () => void;
  onChange: (options: { commitHash?: string }) => void;
};

const useStyles = makeStyles((theme) => ({
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
}));

const DEBOUNCE_INPUT_WAIT = 1000;

const DeploymentTraceFilter: FC<Props> = ({
  filterValues,
  onClear,
  onChange,
}) => {
  const classes = useStyles();
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
      <div className={classes.formItem}>
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
      </div>
    </FilterView>
  );
};

export default DeploymentTraceFilter;
