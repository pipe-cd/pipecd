import {
  Button,
  Divider,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
  CircularProgress,
} from "@material-ui/core";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import React, { FC, useReducer } from "react";
import { useSelector } from "react-redux";
import { APPLICATION_KIND_TEXT } from "../constants/application-kind";
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

const FormSelectInput: FC<{
  name: string;
  value: string;
  items: { name: string; value: string }[];
  onChange: (value: string) => void;
  disabled?: boolean;
}> = ({ name, value, items, onChange, disabled = false }) => (
  <TextField
    fullWidth
    required
    select
    disabled={disabled}
    label={name}
    variant="outlined"
    margin="dense"
    onChange={(e) => onChange(e.target.value)}
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

type FormKey =
  | "name"
  | "env"
  | "pipedId"
  | "repoId"
  | "repoPath"
  | "configFilename"
  | "cloudProvider";
type FormAction =
  | { type: "update-piped"; value: string }
  | { type: "update-kind"; value: ApplicationKind }
  | { type: "update-form-value"; key: FormKey; value: string };
export type AddApplicationFormState = Record<FormKey, string> & {
  kind: ApplicationKind;
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
        repoId: "",
        repoPath: "",
        configFilename: "",
        cloudProvider: "",
      };
    case "update-form-value":
      return { ...state, [action.key]: action.value };
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
    repoId: "",
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

  return (
    <div className={classes.root}>
      <Typography
        className={classes.title}
        variant="h6"
      >{`Add a new application to "${projectName}" project`}</Typography>
      <Divider />
      <form className={classes.form}>
        <TextField
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
          name="Kind"
          value={`${formState.kind}`}
          items={Object.keys(APPLICATION_KIND_TEXT).map((key) => ({
            name: APPLICATION_KIND_TEXT[(key as unknown) as ApplicationKind],
            value: key,
          }))}
          onChange={(value) =>
            dispatch({
              type: "update-kind",
              value: (value as unknown) as ApplicationKind,
            })
          }
          disabled={isAdding}
        />

        <div className={classes.inputGroup}>
          <FormSelectInput
            name="Environment"
            value={formState.env}
            items={environments.map((v) => ({ name: v.name, value: v.id }))}
            onChange={(value) =>
              dispatch({ type: "update-form-value", key: "env", value })
            }
            disabled={isAdding}
          />
          <div className={classes.inputGroupSpace} />
          <FormSelectInput
            name="Piped"
            value={formState.pipedId}
            onChange={(value) => dispatch({ type: "update-piped", value })}
            items={pipeds.map((piped) => ({
              name: `${piped.name} (${piped.id})`,
              value: piped.id,
            }))}
            disabled={isAdding}
          />
        </div>

        <div className={classes.inputGroup}>
          <FormSelectInput
            name="Repository"
            value={formState.repoId}
            onChange={(value) =>
              dispatch({ type: "update-form-value", key: "repoId", value })
            }
            items={
              selectedPiped?.repositoriesList?.map((repo) => ({
                name: repo.id,
                value: repo.id,
              })) || emptyItems
            }
            disabled={selectedPiped === undefined || isAdding}
          />

          <div className={classes.inputGroupSpace} />
          {/** TODO: Check path is accessible */}
          <TextField
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
          name="Cloud Provider"
          value={formState.cloudProvider}
          onChange={(value) =>
            dispatch({
              type: "update-form-value",
              key: "cloudProvider",
              value,
            })
          }
          items={
            selectedPiped?.cloudProvidersList?.map((provider) => ({
              name: provider.name,
              value: provider.name,
            })) || emptyItems
          }
          disabled={selectedPiped === undefined || isAdding}
        />

        <Button
          color="primary"
          type="button"
          onClick={handleSave}
          disabled={isSomeEmptyFormValue() || isAdding}
        >
          SAVE
          {isAdding && (
            <CircularProgress size={24} className={classes.buttonProgress} />
          )}
        </Button>
        <Button onClick={onClose} disabled={isAdding}>
          CANCEL
        </Button>
      </form>
    </div>
  );
};
