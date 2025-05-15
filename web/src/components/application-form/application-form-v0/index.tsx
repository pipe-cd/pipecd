import {
  Box,
  Button,
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Step,
  StepContent,
  StepLabel,
  Stepper,
  TextField,
  Typography,
} from "@mui/material";
import { FC, memo, useEffect, useMemo, useState } from "react";
import {
  APPLICATION_KIND_BY_NAME,
  APPLICATION_KIND_TEXT,
} from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppDispatch, useAppSelector } from "~/hooks/redux";
import {
  addApplication,
  ApplicationGitRepository,
  ApplicationKind,
} from "~/modules/applications";
import {
  ApplicationInfo,
  fetchUnregisteredApplications,
  selectAllUnregisteredApplications,
} from "~/modules/unregistered-applications";
import { sortFunc } from "~/utils/common";
import { ApplicationFormProps } from "..";

import DialogConfirm from "~/components/dialog-confirm";
import { selectAllPipeds } from "~/modules/pipeds";
import { Autocomplete } from "@mui/material";
import { GroupTwoCol } from "../styles";

const ADD_FROM_GIT_CONFIRM_DIALOG_TITLE = "Add Application";
const ADD_FROM_GIT_CONFIRM_DIALOG_DESCRIPTION =
  "Are you sure you want to add the application?";

enum STEP {
  SELECT_PIPED_AND_PLATFORM,
  SELECT_APPLICATION,
  CONFIRM_INFORMATION,
}

const ApplicationFormSuggestionV0: FC<ApplicationFormProps> = ({
  title,
  onClose,
  onFinished: onAdded,
}) => {
  const [activeStep, setActiveStep] = useState(STEP.SELECT_PIPED_AND_PLATFORM);
  const [showConfirm, setShowConfirm] = useState(false);
  const [selectedPipedId, setSelectedPipedId] = useState("");
  const [selectedKind, setSelectedKind] = useState("");
  const [selectedPlatformProvider, setSelectedPlatformProvider] = useState("");
  const [
    selectedApp,
    setSelectedApp,
  ] = useState<ApplicationInfo.AsObject | null>(null);
  const [appToAdd, setAppToAdd] = useState({
    name: "",
    pipedId: "",
    repo: {} as ApplicationGitRepository.AsObject,
    repoPath: "",
    configFilename: "",
    kind: ApplicationKind.KUBERNETES,
    platformProvider: "",
    labels: new Array<[string, string]>(),
  });
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchUnregisteredApplications());
  }, [dispatch]);

  const apps = useAppSelector(selectAllUnregisteredApplications);
  const ps = useAppSelector(selectAllPipeds);

  const appOptions = useMemo(
    () =>
      apps
        .filter(
          (app) =>
            app.pipedId === selectedPipedId &&
            app.kind === APPLICATION_KIND_BY_NAME[selectedKind]
        )
        .sort((a, b) => sortFunc(a.name, b.name)),
    [apps, selectedKind, selectedPipedId]
  );

  const pipedOptions = useMemo(() => {
    return ps
      .filter((piped) => !piped.disabled)
      .sort((a, b) => sortFunc(a.name, b.name));
  }, [ps]);

  const platformProviderOptions = useMemo(() => {
    const selectedPiped = ps.find((piped) => piped.id === selectedPipedId);

    if (!selectedPiped) return [];
    return [
      ...selectedPiped.platformProvidersList,
      ...selectedPiped.cloudProvidersList,
    ];
  }, [ps, selectedPipedId]);

  /**
   * Auto change step based on selectedApp and selectedPipedId
   */
  useEffect(() => {
    if (selectedApp) {
      setActiveStep(STEP.CONFIRM_INFORMATION);
      return;
    }

    if (selectedPlatformProvider) {
      setActiveStep(STEP.SELECT_APPLICATION);
      return;
    }
    setActiveStep(STEP.SELECT_PIPED_AND_PLATFORM);
  }, [selectedApp, selectedPipedId, selectedPlatformProvider]);

  const onSubmitForm = (): void => {
    if (!selectedApp) return;

    setAppToAdd({
      name: selectedApp.name,
      pipedId: selectedApp.pipedId,
      repo: { id: selectedApp.repoId } as ApplicationGitRepository.AsObject,
      repoPath: selectedApp.path,
      configFilename: selectedApp.configFilename,
      kind: selectedApp.kind,
      platformProvider: selectedPlatformProvider,
      labels: selectedApp.labelsMap,
    });
    setShowConfirm(true);
  };

  const onCreateApplication = async (): Promise<void> => {
    await dispatch(addApplication(appToAdd));
    setShowConfirm(false);
    onAdded();
  };

  const onSelectPiped = (value: string): void => {
    setSelectedApp(null);
    setSelectedPipedId(value);
    setSelectedPlatformProvider("");
  };

  const onSelectPlatformProvider = (platformName: string): void => {
    const platformProvider = platformProviderOptions.find(
      (item) => item.name === platformName
    );
    if (!platformProvider) return;

    setSelectedApp(null);
    const kind = platformProvider.type;
    if (kind) setSelectedKind(kind);
    if (platformProvider) setSelectedPlatformProvider(platformName);
    if (platformProvider) setActiveStep(STEP.SELECT_APPLICATION);
  };

  return (
    <>
      <Box width="100%">
        <Typography variant="h6" p={2}>
          {title}
        </Typography>
        <Divider />
        <Stepper activeStep={activeStep} orientation="vertical" sx={{ p: 3 }}>
          <Step key="Select piped and platform provider" active>
            <StepLabel>Select piped and platform provider</StepLabel>
            <StepContent>
              <GroupTwoCol mt={2}>
                <FormControl variant="outlined">
                  <InputLabel id="filter-piped">Piped</InputLabel>
                  <Select
                    labelId="filter-piped"
                    id="filter-piped"
                    label="Piped"
                    value={selectedPipedId}
                    fullWidth
                    onChange={(e) => onSelectPiped(e.target.value as string)}
                  >
                    {pipedOptions.map((e) => (
                      <MenuItem value={e.id} key={e.id}>
                        {e.name} ({e.id})
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>

                <FormControl variant="outlined">
                  <InputLabel id="filter-platform-provider">
                    Platform Provider
                  </InputLabel>
                  <Select
                    labelId="filter-platform-provider"
                    id="filter-platform-provider"
                    label="PlatformProvider"
                    fullWidth
                    disabled={selectedPipedId === ""}
                    value={selectedPlatformProvider}
                    onChange={(e) =>
                      onSelectPlatformProvider(e.target.value as string)
                    }
                  >
                    {platformProviderOptions.map((e) => (
                      <MenuItem value={e.name} key={e.name}>
                        {e.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </GroupTwoCol>
            </StepContent>
          </Step>
          <Step key="Select application to add" expanded={activeStep !== 0}>
            <StepLabel>Select application to add</StepLabel>
            <StepContent>
              <FormControl sx={{ mt: 2 }} variant="outlined" fullWidth>
                <Autocomplete
                  id="filter-app"
                  options={appOptions}
                  getOptionLabel={(app) =>
                    `name: ${app.name}, repo: ${app.repoId}`
                  }
                  value={selectedApp}
                  onChange={(_e, value) => {
                    setSelectedApp(value || null);
                  }}
                  openOnFocus
                  renderInput={(params) => (
                    <TextField
                      {...params}
                      label="Application"
                      variant="outlined"
                    />
                  )}
                />
              </FormControl>
            </StepContent>
          </Step>
          <Step key="Confirm information before adding">
            <StepLabel>Confirm information before adding</StepLabel>
            <StepContent>
              {selectedApp && (
                <Box>
                  <TextField
                    id={"kind"}
                    label="Kind"
                    margin="dense"
                    fullWidth
                    variant="outlined"
                    value={APPLICATION_KIND_TEXT[selectedApp.kind]}
                    inputProps={{ readOnly: true }}
                  />
                  <GroupTwoCol>
                    <TextField
                      id={"path"}
                      label="Path"
                      margin="dense"
                      variant="outlined"
                      value={selectedApp.path}
                      fullWidth
                      inputProps={{ readOnly: true }}
                    />
                    <TextField
                      id={"configFilename-"}
                      label="Config Filename"
                      margin="dense"
                      variant="outlined"
                      value={selectedApp.configFilename}
                      fullWidth
                      inputProps={{ readOnly: true }}
                    />
                  </GroupTwoCol>
                  {selectedApp.labelsMap.map((label, index) => (
                    <Box key={label[0]}>
                      <TextField
                        id={"label-" + "-" + index}
                        label={"Label " + index}
                        margin="dense"
                        variant="outlined"
                        value={label[0] + ": " + label[1]}
                        fullWidth
                        inputProps={{ readOnly: true }}
                      />
                    </Box>
                  ))}
                </Box>
              )}
            </StepContent>
          </Step>
        </Stepper>

        <Box pl={2}>
          <Button
            color="primary"
            type="submit"
            onClick={onSubmitForm}
            disabled={!selectedApp}
          >
            {UI_TEXT_SAVE}
          </Button>

          <Button onClick={onClose}>{UI_TEXT_CANCEL}</Button>
        </Box>
      </Box>

      <DialogConfirm
        open={showConfirm}
        onClose={() => setShowConfirm(false)}
        onCancel={() => setShowConfirm(false)}
        title={ADD_FROM_GIT_CONFIRM_DIALOG_TITLE}
        description={ADD_FROM_GIT_CONFIRM_DIALOG_DESCRIPTION}
        onConfirm={onCreateApplication}
      />
    </>
  );
};

export default memo(ApplicationFormSuggestionV0);
