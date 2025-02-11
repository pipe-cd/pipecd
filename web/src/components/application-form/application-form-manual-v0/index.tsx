import {
  Box,
  Button,
  CircularProgress,
  Divider,
  makeStyles,
  TextField,
  Typography,
} from "@material-ui/core";
import { FC, useEffect } from "react";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { ApplicationKind } from "~/modules/applications";
import { Piped, selectAllPipeds, selectPipedById } from "~/modules/pipeds";
import { sortFunc } from "~/utils/common";
import { ApplicationFormProps } from "..";
import { useFormik } from "formik";
import * as yup from "yup";

import { unwrapResult, useAppDispatch, useAppSelector } from "~/hooks/redux";
import { addApplication } from "~/modules/applications";
import FormSelectInput from "../../form-select-input";
import { updateApplication } from "~/modules/update-application";

type FormValues = {
  name: string;
  kind: ApplicationKind;
  pipedId: string;
  repoPath: string;
  configFilename: string;
  platformProvider: string;
  repo: {
    id: string;
    remote: string;
    branch: string;
  };
  labels: Array<[string, string]>;
};

export const emptyFormValues: FormValues = {
  name: "",
  kind: ApplicationKind.KUBERNETES,
  pipedId: "",
  repoPath: "",
  configFilename: "app.pipecd.yaml",
  platformProvider: "",
  repo: {
    id: "",
    remote: "",
    branch: "",
  },
  labels: new Array<[string, string]>(),
};

const createPlatformProviderListFromPiped = ({
  kind,
  piped,
}: {
  piped?: Piped.AsObject;
  kind: ApplicationKind;
}): Array<{ name: string; value: string }> => {
  if (!piped) {
    return [{ name: "None", value: "" }];
  }

  const providerList: Array<{ name: string; type: string }> = [
    ...piped.cloudProvidersList,
    ...piped.platformProvidersList,
  ];

  return providerList
    .filter((provider) => provider.type === APPLICATION_KIND_TEXT[kind])
    .map((provider) => ({
      name: provider.name,
      value: provider.name,
    }));
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
  kind: yup.number().required(),
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
  platformProvider: yup.string().required(),
});

const ApplicationFormManualV0: FC<ApplicationFormProps> = ({
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
          kind: detailApp.kind,
          pipedId: detailApp.pipedId,
          repoPath: detailApp.gitPath?.path || "",
          repo: detailApp.gitPath?.repo || {
            id: "",
            remote: "",
            branch: "",
          },
          configFilename: detailApp.gitPath?.configFilename || "",
          platformProvider: detailApp.platformProvider,
          labels: detailApp.labelsMap,
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
  const pipeds = ps
    .filter((piped) => !piped.disabled)
    .sort((a, b) => sortFunc(a.name, b.name));

  const selectedPiped = useAppSelector(selectPipedById(values.pipedId));

  const platformProviders = createPlatformProviderListFromPiped({
    piped: selectedPiped,
    kind: values.kind,
  });

  const repositories = createRepoListFromPiped(selectedPiped);

  const disableApplicationInfo = !!detailApp;

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

        <FormSelectInput
          id="kind"
          label="Kind"
          value={`${values.kind}`}
          options={Object.entries(
            APPLICATION_KIND_TEXT
          ).map(([key, label]) => ({ value: String(key), label }))}
          onChange={(value) => {
            setFieldValue("kind", parseInt(value, 10));
          }}
          getOptionLabel={(option) => option.label}
          required
          disabled={isSubmitting || disableApplicationInfo}
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
                kind: values.kind,
                pipedId: value,
              });
            }}
            options={pipeds.map((piped) => ({
              label: `${piped.name} (${piped.id})`,
              value: piped.id,
            }))}
            required
            disabled={isSubmitting || pipeds.length === 0}
          />
          <FormSelectInput
            id="platformProvider"
            label="Platform Provider"
            value={values.platformProvider}
            onChange={(value) => setFieldValue("platformProvider", value)}
            getOptionLabel={(option) => option.name}
            options={platformProviders}
            required
            disabled={
              selectedPiped === undefined ||
              platformProviders.length === 0 ||
              isSubmitting
            }
          />
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

export default ApplicationFormManualV0;
