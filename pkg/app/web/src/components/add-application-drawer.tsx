import {
  Box,
  Button,
  CircularProgress,
  Divider,
  Drawer,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
} from "@material-ui/core";
import { useFormik } from "formik";
import { ApplicationKind } from "pipe/pkg/app/web/model/common_pb";
import React, { FC, ReactElement, memo } from "react";
import { useSelector } from "react-redux";
import * as Yup from "yup";
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

const emptyItems = [{ name: "None", value: "" }];
const createCloudProviderListFromPiped = ({
  kind,
  piped,
}: {
  piped?: Piped;
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

const validationSchema = Yup.object().shape({
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

interface Props {
  open: boolean;
  isAdding: boolean;
  projectName: string;
  onSubmit: (state: {
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
  }) => void;
  onClose: () => void;
}

const initialFormValues = {
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

export const AddApplicationDrawer: FC<Props> = memo(
  function AddApplicationName({
    open,
    isAdding,
    projectName,
    onSubmit,
    onClose,
  }) {
    const classes = useStyles();
    const formik = useFormik({
      initialValues: initialFormValues,
      validateOnMount: true,
      validationSchema,
      onSubmit: (values) => {
        onSubmit(values);
      },
    });

    const environments = useSelector<AppState, Environment[]>((state) =>
      selectEnvironments(state.environments)
    );

    const pipeds = useSelector<AppState, Piped[]>((state) =>
      formik.values.env !== ""
        ? selectPipedsByEnv(state.pipeds, formik.values.env)
        : []
    );

    const selectedPiped = useSelector<AppState, Piped | undefined>((state) =>
      selectPipedById(state.pipeds, formik.values.pipedId)
    );

    const handleClose = (): void => {
      onClose();
      formik.resetForm({ values: initialFormValues });
    };

    const cloudProviders = createCloudProviderListFromPiped({
      piped: selectedPiped,
      kind: formik.values.kind,
    });

    return (
      <Drawer
        anchor="right"
        open={open}
        onClose={handleClose}
        ModalProps={{ disableBackdropClick: isAdding }}
      >
        <Box width={600}>
          <Typography
            className={classes.title}
            variant="h6"
          >{`Add a new application to "${projectName}" project`}</Typography>
          <Divider />
          <form className={classes.form} onSubmit={formik.handleSubmit}>
            <TextField
              id="name"
              name="name"
              label="Name"
              variant="outlined"
              margin="dense"
              onChange={formik.handleChange}
              value={formik.values.name}
              fullWidth
              required
              disabled={isAdding}
              className={classes.textInput}
            />

            <FormSelectInput
              id="kind"
              label="Kind"
              value={`${formik.values.kind}`}
              items={Object.keys(APPLICATION_KIND_TEXT).map((key) => ({
                name:
                  APPLICATION_KIND_TEXT[(key as unknown) as ApplicationKind],
                value: key,
              }))}
              onChange={({ value }) =>
                formik.setFieldValue("kind", parseInt(value, 10))
              }
              disabled={isAdding}
            />

            <div className={classes.inputGroup}>
              <FormSelectInput
                id="env"
                label="Environment"
                value={formik.values.env}
                items={environments.map((v) => ({ name: v.name, value: v.id }))}
                onChange={(item) => {
                  formik.setValues({
                    ...initialFormValues,
                    name: formik.values.name,
                    kind: formik.values.kind,
                    env: item.value,
                  });
                }}
                disabled={isAdding}
              />
              <div className={classes.inputGroupSpace} />
              <FormSelectInput
                id="piped"
                label="Piped"
                value={formik.values.pipedId}
                onChange={({ value }) => {
                  formik.setValues({
                    ...initialFormValues,
                    name: formik.values.name,
                    kind: formik.values.kind,
                    env: formik.values.env,
                    pipedId: value,
                  });
                }}
                items={pipeds.map((piped) => ({
                  name: `${piped.name} (${piped.id})`,
                  value: piped.id,
                }))}
                disabled={isAdding || !formik.values.env || pipeds.length === 0}
              />
            </div>

            <div className={classes.inputGroup}>
              <FormSelectInput
                id="git-repo"
                label="Repository"
                value={formik.values.repo.id || ""}
                onChange={(value) =>
                  formik.setFieldValue("repo", {
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
                disabled={selectedPiped === undefined || isAdding}
              />

              <div className={classes.inputGroupSpace} />
              {/** TODO: Check path is accessible */}
              <TextField
                id="repoPath"
                label="Path"
                variant="outlined"
                margin="dense"
                disabled={selectedPiped === undefined || isAdding}
                onChange={formik.handleChange}
                value={formik.values.repoPath}
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
              disabled={selectedPiped === undefined || isAdding}
              onChange={formik.handleChange}
              value={formik.values.configFilename}
              fullWidth
              className={classes.textInput}
            />

            <FormSelectInput
              id="cloudProvider"
              label="Cloud Provider"
              value={formik.values.cloudProvider}
              onChange={({ value }) =>
                formik.setFieldValue("cloudProvider", value)
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
              type="submit"
              disabled={formik.isValid === false || isAdding}
            >
              {UI_TEXT_SAVE}
              {isAdding && (
                <CircularProgress
                  size={24}
                  className={classes.buttonProgress}
                />
              )}
            </Button>
            <Button onClick={handleClose} disabled={isAdding}>
              {UI_TEXT_CANCEL}
            </Button>
          </form>
        </Box>
      </Drawer>
    );
  }
);
