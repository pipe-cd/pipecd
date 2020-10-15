import {
  Button,
  Divider,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
  CircularProgress,
} from "@material-ui/core";
import {
  ApplicationGitRepository,
  ApplicationKind,
} from "pipe/pkg/app/web/model/common_pb";
import React, { FC, ReactElement, useReducer } from "react";
import { useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../constants/ui-text";
import { AppState } from "../modules";
import {
  Environment,
  selectAll as selectEnvironments,
} from "../modules/environments";
import {
  Piped,
  selectById as selectPipedById,
  selectPipedsByEnv,
} from "../modules/pipeds";

const useStyles = makeStyles((theme) => ({
  root: {
    width: 600,
  },
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
  },
  textInput: {
    flex: 1,
  },
  inputGroup: {
    display: "flex",
  },
  inputGroupSpace: {
    width: theme.spacing(3),
  },
  buttonProgress: {
    color: theme.palette.primary.main,
    position: "absolute",
    top: "50%",
    left: "50%",
    marginTop: -12,
    marginLeft: -12,
  },
}));

function FormSelectInput<T extends { name: string; value: string }>({
  id,
  name,
  value,
  items,
  onChange,
  disabled = false,
}: {
  id: string;
  name: string;
  value: string;
  items: T[];
  onChange: (value: T) => void;
  disabled?: boolean;
}): ReactElement {
  return (
    <TextField
      id={id}
      fullWidth
      required
      select
      disabled={disabled}
      label={name}
      variant="outlined"
      margin="dense"
      onChange={(e) => {
        const nextItem = items.find((item) => item.value === e.target.value);
        if (nextItem) {
          onChange(nextItem);
        }
      }}
      value={value}
      style={{ flex: 1 }}
    >
      {items.map((item) => (
        <MenuItem key={item.name} value={item.value}>
          {item.name}
        </MenuItem>
      ))}
    </TextField>
  );
}

type FormKey =
  | "name"
  | "env"
  | "pipedId"
  | "repoPath"
  | "configFilename"
  | "cloudProvider";
type FormAction =
  | { type: "update-piped"; value: string }
  | { type: "update-kind"; value: ApplicationKind }
  | { type: "update-repo"; value: ApplicationGitRepository.AsObject }
  | { type: "update-form-value"; key: FormKey; value: string };
export type AddApplicationFormState = Record<FormKey, string> & {
  kind: ApplicationKind;
  repo: ApplicationGitRepository.AsObject;
};
function reducer(
  state: AddApplicationFormState,
  action: FormAction
): AddApplicationFormState {
  switch (action.type) {
    case "update-piped":
      // NOTE: clear values that references to piped.
      return {
        ...state,
        pipedId: action.value,
        repoPath: "",
        configFilename: "",
        cloudProvider: "",
        repo: { id: "", remote: "", branch: "" },
      };
    case "update-form-value":
      return { ...state, [action.key]: action.value };
    case "update-repo":
      return { ...state, repo: action.value };
    case "update-kind":
      return { ...state, kind: action.value };
    default:
      return state;
  }
}

const emptyItems = [{ name: "None", value: "" }];

interface Props {
  isAdding: boolean;
  projectName: string;
  onSubmit: (state: AddApplicationFormState) => void;
  onClose: () => void;
}

export const AddApplicationForm: FC<Props> = ({
  isAdding,
  projectName,
  onSubmit,
  onClose,
}) => {
  const classes = useStyles();
  const [formState, dispatch] = useReducer(reducer, {
    name: "",
    env: "",
    pipedId: "",
    repo: { id: "", remote: "", branch: "" },
    repoPath: "",
    configFilename: "",
    kind: ApplicationKind.KUBERNETES, // default value
    cloudProvider: "",
  });

  const environments = useSelector<AppState, Environment[]>((state) =>
    selectEnvironments(state.environments)
  );

  const pipeds = useSelector<AppState, Piped[]>((state) =>
    formState.env ? selectPipedsByEnv(state.pipeds, formState.env) : []
  );

  const selectedPiped = useSelector<AppState, Piped | undefined>((state) =>
    selectPipedById(state.pipeds, formState.pipedId)
  );

  const handleSave = (): void => {
    onSubmit(formState);
  };

  const isSomeEmptyFormValue = (): boolean => {
    return (Object.keys(formState) as FormKey[]).some((key) => {
      // NOTE: configFilename is optional
      if (key === "configFilename") {
        return false;
      }
      return formState[key] === "";
    });
  };

  const cloudProviders =
    selectedPiped?.cloudProvidersList
      ?.filter(
        (provider) =>
          provider.type ===
          APPLICATION_KIND_TEXT[(formState.kind as unknown) as ApplicationKind]
      )
      .map((provider) => ({
        name: provider.name,
        value: provider.name,
      })) || emptyItems;

  return (
    <div className={classes.root}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add a new application to "${projectName}" project`}</Typography>
      <Divider />
      <form className={classes.form}>
        <TextField
          id="application-name"
          label="Name"
          variant="outlined"
          margin="dense"
          onChange={(e) =>
            dispatch({
              type: "update-form-value",
              key: "name",
              value: e.target.value,
            })
          }
          value={formState.name}
          fullWidth
          required
          disabled={isAdding}
          className={classes.textInput}
        />

        <FormSelectInput
          id="application-kind"
          name="Kind"
          value={`${formState.kind}`}
          items={Object.keys(APPLICATION_KIND_TEXT).map((key) => ({
            name: APPLICATION_KIND_TEXT[(key as unknown) as ApplicationKind],
            value: key,
          }))}
          onChange={({ value }) =>
            dispatch({
              type: "update-kind",
              value: parseInt(value, 10) as ApplicationKind,
            })
          }
          disabled={isAdding}
        />

        <div className={classes.inputGroup}>
          <FormSelectInput
            id="application-env"
            name="Environment"
            value={formState.env}
            items={environments.map((v) => ({ name: v.name, value: v.id }))}
            onChange={(item) =>
              dispatch({
                type: "update-form-value",
                key: "env",
                value: item.value,
              })
            }
            disabled={isAdding}
          />
          <div className={classes.inputGroupSpace} />
          <FormSelectInput
            id="application-piped"
            name="Piped"
            value={formState.pipedId}
            onChange={({ value }) => {
              dispatch({ type: "update-piped", value });
            }}
            items={pipeds.map((piped) => ({
              name: `${piped.name} (${piped.id})`,
              value: piped.id,
            }))}
            disabled={isAdding || !formState.env}
          />
        </div>

        <div className={classes.inputGroup}>
          <FormSelectInput
            id="application-git-repo"
            name="Repository"
            value={formState.repo?.id || ""}
            onChange={(value) =>
              dispatch({
                type: "update-repo",
                value: {
                  id: value.value,
                  branch: value.branch,
                  remote: value.remote,
                },
              })
            }
            items={
              selectedPiped?.repositoriesList?.map((repo) => ({
                name: repo.id,
                value: repo.id,
                branch: repo.branch,
                remote: repo.remote,
              })) || []
            }
            disabled={selectedPiped === undefined || isAdding}
          />

          <div className={classes.inputGroupSpace} />
          {/** TODO: Check path is accessible */}
          <TextField
            id="application-repo-path"
            label="Path"
            variant="outlined"
            margin="dense"
            disabled={selectedPiped === undefined || isAdding}
            onChange={(e) =>
              dispatch({
                type: "update-form-value",
                key: "repoPath",
                value: e.target.value,
              })
            }
            value={formState.repoPath}
            fullWidth
            required
            className={classes.textInput}
          />
        </div>

        <TextField
          id="application-config-filename"
          label="Config Filename"
          variant="outlined"
          margin="dense"
          disabled={selectedPiped === undefined || isAdding}
          onChange={(e) =>
            dispatch({
              type: "update-form-value",
              key: "configFilename",
              value: e.target.value,
            })
          }
          value={formState.configFilename}
          fullWidth
          className={classes.textInput}
        />

        <FormSelectInput
          id="application-cloud-provider"
          name="Cloud Provider"
          value={formState.cloudProvider}
          onChange={({ value }) =>
            dispatch({
              type: "update-form-value",
              key: "cloudProvider",
              value,
            })
          }
          items={cloudProviders}
          disabled={
            selectedPiped === undefined ||
            cloudProviders.length === 0 ||
            isAdding
          }
        />

        <Button
          color="primary"
          type="button"
          onClick={handleSave}
          disabled={isSomeEmptyFormValue() || isAdding}
        >
          {UI_TEXT_SAVE}
          {isAdding && (
            <CircularProgress size={24} className={classes.buttonProgress} />
          )}
        </Button>
        <Button onClick={onClose} disabled={isAdding}>
          {UI_TEXT_CANCEL}
        </Button>
      </form>
    </div>
  );
};
