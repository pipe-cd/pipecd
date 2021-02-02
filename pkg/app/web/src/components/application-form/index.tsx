import {
  Box,
  Button,
  CircularProgress,
  Divider,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
} from "@material-ui/core";
import { FormikProps } from "formik";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import React, { FC, memo, ReactElement } from "react";
import { useSelector } from "react-redux";
import * as Yup from "yup";
import { APPLICATION_KIND_TEXT } from "../../constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "../../constants/ui-text";
import { AppState } from "../../modules";
import {
  Environment,
  selectAll as selectEnvironments,
} from "../../modules/environments";
import {
  Piped,
  selectById as selectPipedById,
  selectPipedsByEnv,
} from "../../modules/pipeds";

const emptyItems = [{ name: "None", value: "" }];
const createCloudProviderListFromPiped = ({
  kind,
  piped,
}: {
  piped?: Piped.AsObject;
  kind: ApplicationKind;
}): Array<{ name: string; value: string }> => {
  if (!piped) {
    return emptyItems;
  }

  return piped.cloudProvidersList
    .filter((provider) => provider.type === APPLICATION_KIND_TEXT[kind])
    .map((provider) => ({
      name: provider.name,
      value: provider.name,
    }));
};

const useStyles = makeStyles((theme) => ({
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
  label,
  value,
  items,
  onChange,
  disabled = false,
}: {
  id: string;
  label: string;
  value: string;
  items: T[];
  onChange: (value: T) => void;
  disabled?: boolean;
}): ReactElement {
  return (
    <TextField
      id={id}
      name={id}
      label={label}
      fullWidth
      required
      select
      disabled={disabled}
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

export const validationSchema = Yup.object().shape({
  name: Yup.string().required(),
  kind: Yup.number().required(),
  env: Yup.string().required(),
  pipedId: Yup.string().required(),
  repo: Yup.object({
    id: Yup.string().required(),
    remote: Yup.string().required(),
    branch: Yup.string().required(),
  }).required(),
  repoPath: Yup.string().required(),
  cloudProvider: Yup.string().required(),
});

export interface ApplicationFormValue {
  name: string;
  env: string;
  kind: ApplicationKind;
  pipedId: string;
  repoPath: string;
  configFilename: string;
  cloudProvider: string;
  repo: {
    id: string;
    remote: string;
    branch: string;
  };
}

type Props = FormikProps<ApplicationFormValue> & {
  title: string;
  onClose: () => void;
};

export const emptyFormValues: ApplicationFormValue = {
  name: "",
  env: "",
  kind: ApplicationKind.KUBERNETES,
  pipedId: "",
  repoPath: "",
  configFilename: "",
  cloudProvider: "",
  repo: {
    id: "",
    remote: "",
    branch: "",
  },
};

export const ApplicationForm: FC<Props> = memo(function ApplicationForm({
  title,
  values,
  handleSubmit,
  handleChange,
  isSubmitting,
  isValid,
  setFieldValue,
  setValues,
  onClose,
}) {
  const classes = useStyles();

  const environments = useSelector<AppState, Environment.AsObject[]>((state) =>
    selectEnvironments(state.environments)
  );

  const pipeds = useSelector<AppState, Piped.AsObject[]>((state) =>
    values.env !== "" ? selectPipedsByEnv(state.pipeds, values.env) : []
  );

  const selectedPiped = useSelector<AppState, Piped.AsObject | undefined>(
    (state) => selectPipedById(state.pipeds, values.pipedId)
  );

  const cloudProviders = createCloudProviderListFromPiped({
    piped: selectedPiped,
    kind: values.kind,
  });

  return (
    <Box width={600}>
      <Typography className={classes.title} variant="h6">
        {title}
      </Typography>
      <Divider />
      <form className={classes.form} onSubmit={handleSubmit}>
        <TextField
          id="name"
          name="name"
          label="Name"
          variant="outlined"
          margin="dense"
          onChange={handleChange}
          value={values.name}
          fullWidth
          required
          disabled={isSubmitting}
          className={classes.textInput}
        />

        <FormSelectInput
          id="kind"
          label="Kind"
          value={`${values.kind}`}
          items={Object.keys(APPLICATION_KIND_TEXT).map((key) => ({
            name: APPLICATION_KIND_TEXT[(key as unknown) as ApplicationKind],
            value: key,
          }))}
          onChange={({ value }) => setFieldValue("kind", parseInt(value, 10))}
          disabled={isSubmitting}
        />

        <div className={classes.inputGroup}>
          <FormSelectInput
            id="env"
            label="Environment"
            value={values.env}
            items={environments.map((v) => ({ name: v.name, value: v.id }))}
            onChange={(item) => {
              setValues({
                ...emptyFormValues,
                name: values.name,
                kind: values.kind,
                env: item.value,
              });
            }}
            disabled={isSubmitting}
          />
          <div className={classes.inputGroupSpace} />
          <FormSelectInput
            id="piped"
            label="Piped"
            value={values.pipedId}
            onChange={({ value }) => {
              setValues({
                ...emptyFormValues,
                name: values.name,
                kind: values.kind,
                env: values.env,
                pipedId: value,
              });
            }}
            items={pipeds.map((piped) => ({
              name: `${piped.name} (${piped.id})`,
              value: piped.id,
            }))}
            disabled={isSubmitting || !values.env || pipeds.length === 0}
          />
        </div>

        <div className={classes.inputGroup}>
          <FormSelectInput
            id="git-repo"
            label="Repository"
            value={values.repo.id || ""}
            onChange={(value) =>
              setFieldValue("repo", {
                id: value.value,
                branch: value.branch,
                remote: value.remote,
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
            disabled={selectedPiped === undefined || isSubmitting}
          />

          <div className={classes.inputGroupSpace} />
          {/** TODO: Check path is accessible */}
          <TextField
            id="repoPath"
            label="Path"
            variant="outlined"
            margin="dense"
            disabled={selectedPiped === undefined || isSubmitting}
            onChange={handleChange}
            value={values.repoPath}
            fullWidth
            required
            className={classes.textInput}
          />
        </div>

        <TextField
          id="configFilename"
          label="Config Filename"
          variant="outlined"
          margin="dense"
          disabled={selectedPiped === undefined || isSubmitting}
          onChange={handleChange}
          value={values.configFilename}
          fullWidth
          className={classes.textInput}
        />

        <FormSelectInput
          id="cloudProvider"
          label="Cloud Provider"
          value={values.cloudProvider}
          onChange={({ value }) => setFieldValue("cloudProvider", value)}
          items={cloudProviders}
          disabled={
            selectedPiped === undefined ||
            cloudProviders.length === 0 ||
            isSubmitting
          }
        />

        <Button
          color="primary"
          type="submit"
          disabled={isValid === false || isSubmitting}
        >
          {UI_TEXT_SAVE}
          {isSubmitting && (
            <CircularProgress size={24} className={classes.buttonProgress} />
          )}
        </Button>
        <Button onClick={onClose} disabled={isSubmitting}>
          {UI_TEXT_CANCEL}
        </Button>
      </form>
    </Box>
  );
});
