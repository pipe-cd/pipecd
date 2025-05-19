import {
  Box,
  Button,
  CircularProgress,
  Divider,
  FormControl,
  TextField,
  Typography,
} from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC, useEffect, useMemo } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { Piped, selectAllPipeds, selectPipedById } from "~/modules/pipeds";
import { sortFunc } from "~/utils/common";
import { ApplicationFormProps } from "..";
import { useFormik } from "formik";
import * as yup from "yup";

import { unwrapResult, useAppDispatch, useAppSelector } from "~/hooks/redux";
import { addApplication } from "~/modules/applications";
import FormSelectInput from "../../form-select-input";
import { updateApplication } from "~/modules/update-application";
import { Autocomplete } from "@mui/material";

type FormValues = {
  name: string;
  pipedId: string;
  repoPath: string;
  configFilename: string;
  deployTargets: { pluginName: string; deployTarget: string }[];
  repo: {
    id: string;
    remote: string;
    branch: string;
  };
  labels: Array<[string, string]>;
};

type DeployTargetOption = {
  pluginName: string;
  deployTarget: string;
  value: string;
};

export const emptyFormValues: FormValues = {
  name: "",
  pipedId: "",
  repoPath: "",
  deployTargets: [],
  configFilename: "app.pipecd.yaml",
  repo: {
    id: "",
    remote: "",
    branch: "",
  },
  labels: new Array<[string, string]>(),
};

const createRepoListFromPiped = (
  piped?: Piped.AsObject
): Array<{ name: string; value: string; branch: string; remote: string }> => {
  if (!piped) {
    return [
      {
        name: "None",
        value: "",
        branch: "",
        remote: "",
      },
    ];
  }

  return piped.repositoriesList.map((repo) => ({
    name: repo.id,
    value: repo.id,
    branch: repo.branch,
    remote: repo.remote,
  }));
};

const useStyles = makeStyles((theme) => ({
  title: {
    padding: theme.spacing(2),
  },
  form: {
    padding: theme.spacing(2),
    display: "grid",
    gap: theme.spacing(2),
  },
  textInput: {
    flex: 1,
  },
  inputGroup: {
    display: "grid",
    gridTemplateColumns: "1fr 1fr",
    gap: theme.spacing(3),
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

const validationSchema = yup.object().shape({
  name: yup.string().required(),
  pipedId: yup.string().required(),
  repo: yup
    .object({
      id: yup.string().required(),
      remote: yup.string().required(),
      branch: yup.string().required(),
    })
    .required(),
  repoPath: yup.string().required(),
  configFilename: yup.string().required(),
});

const ApplicationFormManualV1: FC<ApplicationFormProps> = ({
  title,
  onClose,
  onFinished,
  setIsFormDirty,
  setIsSubmitting,
  detailApp: detailApp,
}) => {
  const dispatch = useAppDispatch();
  const formik = useFormik<FormValues>({
    initialValues: detailApp
      ? {
          name: detailApp.name,
          pipedId: detailApp.pipedId,
          repoPath: detailApp.gitPath?.path || "",
          repo: detailApp.gitPath?.repo || {
            id: "",
            remote: "",
            branch: "",
          },
          configFilename: detailApp.gitPath?.configFilename || "",
          labels: detailApp.labelsMap,
          deployTargets: detailApp.deployTargetsByPluginMap.reduce(
            (all, [pluginName, { deployTargetsList }]) => {
              deployTargetsList.forEach((deployTarget) => {
                all.push({ pluginName, deployTarget });
              });

              return all;
            },
            [] as { pluginName: string; deployTarget: string }[]
          ),
        }
      : emptyFormValues,
    validationSchema,
    enableReinitialize: true,

    async onSubmit(values) {
      if (detailApp) {
        await dispatch(
          updateApplication({
            ...values,
            applicationId: detailApp.id,
          })
        )
          .then(unwrapResult)
          .then(() => {
            formik.resetForm();
            onFinished();
          });
      }
      if (!detailApp) {
        await dispatch(addApplication(values))
          .then(unwrapResult)
          .then(() => {
            formik.resetForm();
            onFinished();
          });
      }
    },
  });

  const {
    values,
    handleSubmit,
    handleChange,
    isSubmitting,
    isValid,
    dirty,
    setFieldValue,
    setValues,
  } = formik;

  useEffect(() => {
    setIsFormDirty?.(dirty);
  }, [dirty, setIsFormDirty]);

  useEffect(() => {
    setIsSubmitting?.(isSubmitting);
  }, [isSubmitting, setIsSubmitting]);

  const classes = useStyles();
  const ps = useAppSelector((state) => selectAllPipeds(state));
  const pipedOptions = ps
    .filter((piped) => !piped.disabled || piped.id === detailApp?.pipedId)
    .sort((a, b) => sortFunc(a.name, b.name));

  const selectedPiped = useAppSelector(selectPipedById(values.pipedId));

  const repositories = createRepoListFromPiped(selectedPiped);

  const disableApplicationInfo = !!detailApp;

  const deployTargetOptions = useMemo(() => {
    if (!selectedPiped) return [];
    if (selectedPiped.pluginsList.length === 0) return [];

    return selectedPiped.pluginsList.reduce((all, plugin) => {
      plugin.deployTargetsList.forEach((deployTarget) => {
        all.push({
          deployTarget,
          pluginName: plugin.name,
          value: `${deployTarget} - ${plugin.name}`,
        });
      });
      return all;
    }, [] as DeployTargetOption[]);
  }, [selectedPiped]);

  return (
    <Box width="100%">
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
          onChange={handleChange}
          value={values.name}
          fullWidth
          required
          disabled={isSubmitting || disableApplicationInfo}
          className={classes.textInput}
        />
        <div className={classes.inputGroup}>
          <FormSelectInput
            id="piped"
            label="Piped"
            value={values.pipedId}
            onChange={(value) => {
              setValues({
                ...emptyFormValues,
                name: values.name,
                pipedId: value,
              });
            }}
            options={pipedOptions.map((piped) => ({
              label: `${piped.name} (${piped.id})`,
              value: piped.id,
              disabled: piped.disabled,
            }))}
            required
            disabled={isSubmitting || pipedOptions.length === 0}
          />

          <FormControl variant="outlined">
            <Autocomplete
              id="deploy-targets"
              options={deployTargetOptions.map(({ value }) => value)}
              multiple={true}
              value={values.deployTargets.map(
                (item) => `${item.deployTarget} - ${item.pluginName}`
              )}
              disabled={isSubmitting || pipedOptions.length === 0}
              onChange={(_e, value) => {
                const selected = deployTargetOptions.filter((item) =>
                  value.includes(item.value)
                );
                setFieldValue("deployTargets", selected);
              }}
              openOnFocus
              autoComplete={false}
              noOptionsText="No deploy targets found"
              renderInput={(params) => (
                <TextField
                  {...params}
                  label="Deploy targets"
                  variant="outlined"
                />
              )}
            />
          </FormControl>
        </div>

        <div className={classes.inputGroup}>
          <FormSelectInput
            id="git-repo"
            label="Repository"
            value={values.repo.id || ""}
            getOptionLabel={(option) => option.name}
            options={repositories}
            onChange={(_value, item) =>
              setFieldValue("repo", {
                id: item.value,
                branch: item.branch,
                remote: item.remote,
              })
            }
            required
            disabled={
              selectedPiped === undefined ||
              repositories.length === 0 ||
              isSubmitting ||
              disableApplicationInfo
            }
          />

          <TextField
            id="repoPath"
            label="Path"
            placeholder="Relative path to app directory"
            variant="outlined"
            disabled={
              selectedPiped === undefined ||
              isSubmitting ||
              disableApplicationInfo
            }
            onChange={handleChange}
            value={values.repoPath}
            fullWidth
            required
            // className={classes.textInput}
          />
        </div>

        <TextField
          id="configFilename"
          label="Config Filename"
          variant="outlined"
          disabled={selectedPiped === undefined || isSubmitting}
          onChange={handleChange}
          value={values.configFilename}
          fullWidth
          required
          className={classes.textInput}
        />

        <Box my={2}>
          <Button
            color="primary"
            type="submit"
            disabled={isValid === false || isSubmitting || dirty === false}
          >
            {UI_TEXT_SAVE}
            {isSubmitting && (
              <CircularProgress size={24} className={classes.buttonProgress} />
            )}
          </Button>
          <Button onClick={onClose} disabled={isSubmitting}>
            {UI_TEXT_CANCEL}
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default ApplicationFormManualV1;
