import {
  Box,
  Button,
  Divider,
  FormControl,
  TextField,
  Typography,
} from "@mui/material";
import { FC, useEffect, useMemo } from "react";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { sortFunc } from "~/utils/common";
import { ApplicationFormProps } from "..";
import { useFormik } from "formik";
import * as yup from "yup";

import FormSelectInput from "../../form-select-input";
import { Autocomplete } from "@mui/material";
import { GroupTwoCol, StyledForm } from "../styles";
import { SpinnerIcon } from "~/styles/button";
import { useAddApplication } from "~/queries/applications/use-add-application";
import { useUpdateApplication } from "~/queries/applications/use-update-application";
import { useGetPipeds } from "~/queries/pipeds/use-get-pipeds";
import { Piped } from "~~/model/piped_pb";

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
  const { mutate: addApplication } = useAddApplication();
  const { mutate: updateApplication } = useUpdateApplication();

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
        <GroupTwoCol>
          <FormSelectInput
            id="piped"
            label="Piped"
            value={isLoadingPiped ? "" : values.pipedId}
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

export default ApplicationFormManualV1;
