import { Box, Button, Divider, TextField, Typography } from "@mui/material";
import { FC, useEffect, useMemo } from "react";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { sortFunc } from "~/utils/common";
import { ApplicationFormProps } from "..";
import { useFormik } from "formik";
import * as yup from "yup";

import FormSelectInput from "../../form-select-input";
import { GroupTwoCol, StyledForm } from "../styles";
import { SpinnerIcon } from "~/styles/button";
import { useAddApplication } from "~/queries/applications/use-add-application";
import { useUpdateApplication } from "~/queries/applications/use-update-application";
import { useGetPipeds } from "~/queries/pipeds/use-get-pipeds";
import { Piped } from "~~/model/piped_pb";
import { ApplicationKind } from "~/types/applications";

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
  detailApp,
}) => {
  const { mutate: addApplication } = useAddApplication();
  const { mutate: updateApplication } = useUpdateApplication();

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
        updateApplication(
          {
            ...values,
            applicationId: detailApp.id,
          },
          {
            onSuccess: () => {
              formik.resetForm();
              onFinished();
            },
          }
        );
      }
      if (!detailApp) {
        addApplication(values, {
          onSuccess: () => {
            formik.resetForm();
            onFinished();
          },
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

  const { data: ps = [], isLoading: isLoadingPiped } = useGetPipeds({
    withStatus: true,
  });
  const pipedOptions = ps
    .filter((piped) => !piped.disabled || piped.id === detailApp?.pipedId)
    .sort((a, b) => sortFunc(a.name, b.name));

  const selectedPiped = useMemo(() => {
    return pipedOptions.find((piped) => piped.id === values.pipedId);
  }, [pipedOptions, values.pipedId]);

  const platformProviders = createPlatformProviderListFromPiped({
    piped: selectedPiped,
    kind: values.kind,
  });

  const repositories = createRepoListFromPiped(selectedPiped);

  const disableApplicationInfo = !!detailApp;

  return (
    <Box
      sx={{
        width: "100%",
      }}
    >
      <Typography
        variant="h6"
        sx={{
          p: 2,
        }}
      >
        {title}
      </Typography>
      <Divider />
      <StyledForm onSubmit={handleSubmit}>
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

        <GroupTwoCol>
          <FormSelectInput
            id="piped"
            label="Piped"
            value={isLoadingPiped ? "" : values.pipedId}
            onChange={(value) => {
              setValues({
                ...emptyFormValues,
                name: values.name,
                kind: values.kind,
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
          <FormSelectInput
            id="platformProvider"
            label="Platform Provider"
            value={isLoadingPiped ? "" : values.platformProvider}
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
        </GroupTwoCol>

        <GroupTwoCol>
          <FormSelectInput
            id="git-repo"
            label="Repository"
            value={isLoadingPiped ? "" : values.repo.id || ""}
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
        </GroupTwoCol>

        <TextField
          id="configFilename"
          label="Config Filename"
          variant="outlined"
          disabled={selectedPiped === undefined || isSubmitting}
          onChange={handleChange}
          value={values.configFilename}
          fullWidth
          required
        />

        <Box
          sx={{
            my: 2,
          }}
        >
          <Button
            color="primary"
            type="submit"
            disabled={isValid === false || isSubmitting || dirty === false}
          >
            {UI_TEXT_SAVE}
            {isSubmitting && <SpinnerIcon />}
          </Button>
          <Button onClick={onClose} disabled={isSubmitting}>
            {UI_TEXT_CANCEL}
          </Button>
        </Box>
      </StyledForm>
    </Box>
  );
};

export default ApplicationFormManualV0;
