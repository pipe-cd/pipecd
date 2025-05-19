import {
  Box,
  Button,
  Divider,
  MenuItem,
  TextField,
  Typography,
  FormControl,
  InputLabel,
  Select,
  Stepper,
  Step,
  StepLabel,
  StepContent,
} from "@mui/material";
import { FC, memo, useCallback, useState, useEffect, useMemo } from "react";
import { APPLICATION_KIND_TEXT } from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppSelector, useAppDispatch, unwrapResult } from "~/hooks/redux";
import {
  selectAllUnregisteredApplications,
  fetchUnregisteredApplications,
  ApplicationInfo,
} from "~/modules/unregistered-applications";
import {
  addApplication,
  ApplicationGitRepository,
} from "~/modules/applications";
import { sortFunc } from "~/utils/common";
import { selectAllPipeds } from "~/modules/pipeds";
import { Autocomplete } from "@mui/material";
import DialogConfirm from "~/components/dialog-confirm";
import { GroupTwoCol } from "../styles";

const ADD_FROM_GIT_CONFIRM_DIALOG_TITLE = "Add Application";
const ADD_FROM_GIT_CONFIRM_DIALOG_DESCRIPTION =
  "Are you sure you want to add the application?";

enum STEP {
  SELECT_PIPED,
  SELECT_APPLICATION,
  CONFIRM_INFORMATION,
}

type DeployTargetOption = {
  pluginName: string;
  deployTarget: string;
  value: string;
};

type Props = {
  title: string;
  onClose: () => void;
  onFinished: () => void;
};

const ApplicationFormSuggestionV1: FC<Props> = ({
  title,
  onClose,
  onFinished: onAdded,
}) => {
  const [loading, setLoading] = useState(false);
  const [activeStep, setActiveStep] = useState(STEP.SELECT_PIPED);
  const [showConfirm, setShowConfirm] = useState(false);
  const [selectedPipedId, setSelectedPipedId] = useState("");
  const [selectedDeployTargets, setSelectedDeployTargets] = useState<
    DeployTargetOption[]
  >([]);
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
    labels: new Array<[string, string]>(),
    deployTargets: new Array<DeployTargetOption>(),
  });
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch(fetchUnregisteredApplications());
  }, [dispatch]);

  const apps = useAppSelector(selectAllUnregisteredApplications);
  const ps = useAppSelector(selectAllPipeds);

  const selectedPiped = useMemo(
    () => ps.find((piped) => piped.id === selectedPipedId),
    [ps, selectedPipedId]
  );

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

  const filteredApps = useMemo(
    () =>
      apps
        .filter((app) => app.pipedId === selectedPipedId)
        .sort((a, b) => sortFunc(a.name, b.name)),
    [apps, selectedPipedId]
  );

  const pipedOptions = useMemo(() => {
    return ps
      .filter((piped) => !piped.disabled)
      .sort((a, b) => sortFunc(a.name, b.name));
  }, [ps]);

  /**
   * Auto change step based on selectedApp and selectedPipedId
   */
  useEffect(() => {
    if (selectedApp) {
      setActiveStep(STEP.CONFIRM_INFORMATION);
      return;
    }

    if (selectedPipedId && selectedDeployTargets.length > 0) {
      setActiveStep(STEP.SELECT_APPLICATION);
      return;
    }
    setActiveStep(STEP.SELECT_PIPED);
  }, [selectedApp, selectedDeployTargets.length, selectedPipedId]);

  /**
   * Init selectedPipedId if there is only one piped
   */
  useEffect(() => {
    if (pipedOptions.length === 1 && !selectedApp) {
      setSelectedPipedId(pipedOptions[0].id);
    }
  }, [pipedOptions, selectedApp]);

  /**
   * Init selectedApp if there is only one app
   */
  useEffect(() => {
    if (filteredApps.length === 1 && !selectedApp) {
      setSelectedApp(filteredApps[0]);
    }
  }, [apps.length, filteredApps, filteredApps.length, selectedApp]);

  const onSelectPiped = useCallback((pipedId: string) => {
    setSelectedPipedId(pipedId);
    setSelectedDeployTargets([]);
    setSelectedApp(null);
  }, []);

  const onSelectDeployTargets = useCallback((targets: DeployTargetOption[]) => {
    setSelectedDeployTargets(targets);
    setSelectedApp(null);
  }, []);

  const onSubmitForm = (): void => {
    if (!selectedApp) return;

    setAppToAdd({
      name: selectedApp.name,
      pipedId: selectedApp.pipedId,
      repo: { id: selectedApp.repoId } as ApplicationGitRepository.AsObject,
      repoPath: selectedApp.path,
      configFilename: selectedApp.configFilename,
      labels: selectedApp.labelsMap,
      deployTargets: selectedDeployTargets,
    });
    setShowConfirm(true);
  };

  const onCreateApplication = (): void => {
    setLoading(true);
    dispatch(addApplication(appToAdd))
      .then(unwrapResult)
      .then(() => {
        onAdded();
      })
      .finally(() => {
        setLoading(true);
        setShowConfirm(false);
      });
  };

  return (
    <>
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
        <Stepper activeStep={activeStep} orientation="vertical" sx={{ p: 3 }}>
          <Step key="Select piped and deploy targets" active>
            <StepLabel>Select piped and deploy targets</StepLabel>
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
                    defaultValue={""}
                    onChange={(e) => {
                      onSelectPiped(e.target.value as string);
                    }}
                  >
                    {pipedOptions.map((e) => (
                      <MenuItem value={e.id} key={`piped-${e.id}`}>
                        {e.name} ({e.id})
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
                <FormControl variant="outlined">
                  <Autocomplete
                    id="deploy-targets"
                    options={deployTargetOptions.map(({ value }) => value)}
                    multiple={true}
                    value={selectedDeployTargets.map((item) => item.value)}
                    disabled={!selectedPipedId}
                    onChange={(_e, value) => {
                      const selected = deployTargetOptions.filter((item) =>
                        value.includes(item.value)
                      );
                      onSelectDeployTargets(selected);
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
            </StepContent>
          </Step>
          <Step key="Select application to add" expanded={activeStep !== 0}>
            <StepLabel>Select application to add</StepLabel>
            <StepContent>
              <FormControl sx={{ mt: 2 }} variant="outlined" fullWidth>
                <Autocomplete
                  id="filter-app"
                  options={filteredApps}
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

        <Box
          sx={{
            pl: 2,
          }}
        >
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
        loading={loading}
      />
    </>
  );
};

export default memo(ApplicationFormSuggestionV1);
