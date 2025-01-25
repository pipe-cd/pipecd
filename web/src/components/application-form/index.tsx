import {
  Box,
  Button,
  CircularProgress,
  Divider,
  makeStyles,
  MenuItem,
  TextField,
  Typography,
  Tabs,
  Tab,
  FormControl,
  InputLabel,
  Select,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stepper,
  Step,
  StepLabel,
  StepContent,
  IconButton,
} from "@material-ui/core";
import { Help } from "@material-ui/icons";
import { FormikProps } from "formik";
import {
  FC,
  memo,
  ReactElement,
  useCallback,
  useState,
  useEffect,
} from "react";
import * as yup from "yup";
import {
  APPLICATION_KIND_TEXT,
  APPLICATION_KIND_BY_NAME,
} from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE } from "~/constants/ui-text";
import { useAppSelector, useAppDispatch } from "~/hooks/redux";
import { ApplicationKind } from "~/modules/applications";
import { Piped, selectAllPipeds, selectPipedById } from "~/modules/pipeds";
import {
  ApplicationInfo,
  selectAllUnregisteredApplications,
  fetchUnregisteredApplications,
} from "~/modules/unregistered-applications";
import {
  addApplication,
  ApplicationGitRepository,
} from "~/modules/applications";
import ApplicationFormV1 from "./application-form-v1";

const ADD_FROM_GIT_CONFIRM_DIALOG_TITLE = "Add Application";
const ADD_FROM_GIT_CONFIRM_DIALOG_DESCRIPTION =
  "Are you sure you want to add the application?";

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
  formItem: {
    width: "100%",
    marginTop: theme.spacing(4),
  },
  select: {
    width: "100%",
  },
  accordionDetail: {
    width: "100%",
  },
  button: {
    margin: theme.spacing(2),
  },
  actionsContainer: {
    marginBottom: theme.spacing(2),
  },
  tabLabel: {
    minHeight: 0,
    "& .MuiTab-wrapper": {
      flexDirection: "row-reverse",
      maxWidth: 200,
    },
    "& .MuiTab-wrapper > *:first-child": {
      marginBottom: 0,
    },
    "& .MuiIconButton-sizeSmall": {
      padding: "0 3px 3px 3px",
    },
  },
}));

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  selected: boolean;
}

function TabPanel(props: TabPanelProps): ReactElement {
  return (
    <div
      role="tabpanel"
      hidden={!props.selected}
      id={`simple-tabpanel-${props.index}`}
      aria-labelledby={`simple-tab-${props.index}`}
    >
      {props.selected && (
        <Box>
          <Typography>{props.children}</Typography>
        </Box>
      )}
    </div>
  );
}

function a11yProps(index: number): { id: string; "aria-controls": string } {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

const sortComp = (a: string, b: string): number => {
  return a > b ? 1 : -1;
};

export const ApplicationFormTabs: React.FC<ApplicationFormProps> = (props) => {
  const classes = useStyles();

  const [selectedTabIndex, setSelectedTabIndex] = useState(0);

  const handleChange = (
    event: React.ChangeEvent<Record<string, unknown>>,
    newValue: number
  ): void => {
    setSelectedTabIndex(newValue);
  };

  return (
    <Box width={600}>
      <Box>
        <Tabs
          value={selectedTabIndex}
          onChange={handleChange}
          aria-label="basic tabs example"
        >
          <Tab
            className={classes.tabLabel}
            label="PIPED V0 ADD FROM SUGGESTIONS"
            icon={
              <IconButton
                size="small"
                href="https://pipecd.dev/docs/user-guide/managing-application/adding-an-application/#picking-from-a-list-of-unused-apps-suggested-by-pipeds"
                target="_blank"
                rel="noopener noreferrer"
              >
                <Help fontSize="small" />
              </IconButton>
            }
            {...a11yProps(0)}
          />
          <Tab
            className={classes.tabLabel}
            label="PIPED V1 ADD FROM SUGGESTIONS"
            {...a11yProps(1)}
          />
          <Tab
            className={classes.tabLabel}
            label="ADD MANUALLY"
            icon=" "
            {...a11yProps(2)}
          />
        </Tabs>
      </Box>
      <TabPanel selected={selectedTabIndex === 0} index={0}>
        <SelectFromSuggestionsForm {...props} />
      </TabPanel>
      <TabPanel selected={selectedTabIndex === 1} index={1}>
        <ApplicationFormV1 {...props} />
      </TabPanel>
      <TabPanel selected={selectedTabIndex === 2} index={2}>
        <ApplicationForm {...props} />
      </TabPanel>
    </Box>
  );
};

function FormSelectInput<T extends { name: string; value: string }>({
  id,
  label,
  value,
  items,
  required = true,
  onChange,
  disabled = false,
}: {
  id: string;
  label: string;
  value: string;
  items: T[];
  required?: boolean;
  onChange: (value: T) => void;
  disabled?: boolean;
}): ReactElement {
  return (
    <TextField
      id={id}
      name={id}
      label={label}
      fullWidth
      required={required}
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

export const validationSchema = yup.object().shape({
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

export interface ApplicationFormValue {
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
}

export type ApplicationFormProps = FormikProps<ApplicationFormValue> & {
  title: string;
  onClose: () => void;
  disableApplicationInfo?: boolean;
};

export const emptyFormValues: ApplicationFormValue = {
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

export const ApplicationForm: FC<ApplicationFormProps> = memo(
  function ApplicationForm({
    title,
    values,
    handleSubmit,
    handleChange,
    isSubmitting,
    isValid,
    dirty,
    setFieldValue,
    setValues,
    onClose,
    disableApplicationInfo = false,
  }) {
    const classes = useStyles();
    const ps = useAppSelector((state) => selectAllPipeds(state));
    const pipeds = ps
      .filter((piped) => !piped.disabled)
      .sort((a, b) => sortComp(a.name, b.name));

    const selectedPiped = useAppSelector(selectPipedById(values.pipedId));

    const platformProviders = createPlatformProviderListFromPiped({
      piped: selectedPiped,
      kind: values.kind,
    });

    const repositories = createRepoListFromPiped(selectedPiped);

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
            margin="dense"
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
            items={Object.keys(APPLICATION_KIND_TEXT).map((key) => ({
              name: APPLICATION_KIND_TEXT[(key as unknown) as ApplicationKind],
              value: key,
            }))}
            onChange={({ value }) => setFieldValue("kind", parseInt(value, 10))}
            disabled={isSubmitting || disableApplicationInfo}
          />

          <div className={classes.inputGroup}>
            <FormSelectInput
              id="piped"
              label="Piped"
              value={values.pipedId}
              onChange={({ value }) => {
                setValues({
                  ...emptyFormValues,
                  name: values.name,
                  kind: values.kind,
                  pipedId: value,
                });
              }}
              items={pipeds.map((piped) => ({
                name: `${piped.name} (${piped.id})`,
                value: piped.id,
              }))}
              disabled={isSubmitting || pipeds.length === 0}
            />
            <div className={classes.inputGroupSpace} />
            <FormSelectInput
              id="platformProvider"
              label="Platform Provider"
              value={values.platformProvider}
              onChange={({ value }) => setFieldValue("platformProvider", value)}
              items={platformProviders}
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
              onChange={(value) =>
                setFieldValue("repo", {
                  id: value.value,
                  branch: value.branch,
                  remote: value.remote,
                })
              }
              items={repositories}
              disabled={
                selectedPiped === undefined ||
                repositories.length === 0 ||
                isSubmitting ||
                disableApplicationInfo
              }
            />

            <div className={classes.inputGroupSpace} />
            {/** TODO: Check path is accessible */}
            <TextField
              id="repoPath"
              label="Path"
              placeholder="Relative path to app directory"
              variant="outlined"
              margin="dense"
              disabled={
                selectedPiped === undefined ||
                isSubmitting ||
                disableApplicationInfo
              }
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
            required
            className={classes.textInput}
          />

          <Box m={2} />
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
        </form>
      </Box>
    );
  }
);

interface PlatformProviderFilterOptions {
  pipedId: string;
  platformProvider: string;
  kind: string;
}

interface PlatformProviderFilterProps {
  onChange: (options: PlatformProviderFilterOptions) => void;
}

const PlatformProviderFilter: FC<PlatformProviderFilterProps> = memo(
  function PlatformProviderFilter({ onChange }) {
    const classes = useStyles();
    const ps = useAppSelector((state) => selectAllPipeds(state));
    const pipeds = ps
      .filter((piped) => !piped.disabled)
      .sort((a, b) => sortComp(a.name, b.name));

    const [selectedPipedId, setSelectedPipedId] = useState(
      pipeds.length === 1 ? pipeds[0].id : ""
    );
    const selectedPiped = useAppSelector(selectPipedById(selectedPipedId));
    const platformProviders: Array<{
      name: string;
      type: string;
    }> = selectedPiped
      ? [
          ...selectedPiped.platformProvidersList,
          ...selectedPiped.cloudProvidersList,
        ]
      : [];

    let options: PlatformProviderFilterOptions;
    const handleUpdateFilterValue = (
      optionPart: Partial<PlatformProviderFilterOptions>
    ): void => {
      onChange({ ...options, ...optionPart });
    };

    return (
      <div className={classes.inputGroup}>
        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-piped">Piped</InputLabel>
          <Select
            labelId="filter-piped"
            id="filter-piped"
            label="Piped"
            value={selectedPipedId}
            className={classes.select}
            onChange={(e) => {
              setSelectedPipedId(e.target.value as string);
              handleUpdateFilterValue({
                pipedId: e.target.value as string,
              });
            }}
          >
            {pipeds.map((e) => (
              <MenuItem value={e.id} key={`piped-${e.id}`}>
                {e.name} ({e.id})
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <div className={classes.inputGroupSpace} />
        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-platform-provider">
            Platform Provider
          </InputLabel>
          <Select
            labelId="filter-platform-provider"
            id="filter-platform-provider"
            label="PlatformProvider"
            className={classes.select}
            disabled={selectedPipedId === ""}
            onChange={(e) => {
              const values = e.target.value as ReadonlyArray<string>;
              handleUpdateFilterValue({
                platformProvider: values[0],
                kind: values[1],
                pipedId: selectedPipedId,
              });
            }}
          >
            {platformProviders.map((e) => (
              <MenuItem
                value={[e.name, e.type] as ReadonlyArray<string>}
                key={`platform-provider-${e.name}`}
              >
                {e.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </div>
    );
  }
);

const SelectFromSuggestionsForm: FC<ApplicationFormProps> = memo(
  function ApplicationForm({ title, onClose }) {
    const dispatch = useAppDispatch();
    useEffect(() => {
      dispatch(fetchUnregisteredApplications());
    }, [dispatch]);

    const classes = useStyles();
    const apps = useAppSelector<ApplicationInfo.AsObject[]>((state) =>
      selectAllUnregisteredApplications(state)
    );

    const [selectedPipedId, setSelectedPipedId] = useState("");
    const [selectedKind, setSelectedKind] = useState("");
    const [selectedPlatformProvider, setSelectedPlatformProvider] = useState(
      ""
    );
    const [selectedAppIndex, setSelectedAppIndex] = useState(-1);
    const [
      selectedApp,
      setSelectedApp,
    ] = useState<ApplicationInfo.AsObject | null>();
    const [filteredApps, setFilteredApps] = useState<
      ApplicationInfo.AsObject[]
    >([]);

    useEffect(() => {
      setFilteredApps(
        apps
          .filter(
            (app) =>
              app.pipedId === selectedPipedId &&
              app.kind === APPLICATION_KIND_BY_NAME[selectedKind]
          )
          .sort((a, b) => sortComp(a.name, b.name))
      );
    }, [apps, selectedPipedId, selectedKind]);

    const [showConfirm, setShowConfirm] = useState(false);

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

    const handleFilterChange = useCallback(
      (options: PlatformProviderFilterOptions) => {
        setSelectedApp(null);
        setSelectedAppIndex(-1);
        setSelectedPipedId(options.pipedId);
        setSelectedKind(options.kind);
        setSelectedPlatformProvider(options.platformProvider);
        setActiveStep(options.platformProvider ? 1 : 0);
      },
      []
    );

    const [activeStep, setActiveStep] = useState(0);

    return (
      <>
        <Box width="100%">
          <Typography className={classes.title} variant="h6">
            {title}
          </Typography>
          <Divider />
          <Stepper activeStep={activeStep} orientation="vertical">
            <Step key="Select piped and platform provider" active>
              <StepLabel>Select piped and platform provider</StepLabel>
              <StepContent>
                <div className={classes.actionsContainer}>
                  <div>
                    <PlatformProviderFilter onChange={handleFilterChange} />
                  </div>
                </div>
              </StepContent>
            </Step>
            <Step key="Select application to add" expanded={activeStep !== 0}>
              <StepLabel>Select application to add</StepLabel>
              <StepContent>
                <FormControl className={classes.formItem} variant="outlined">
                  <InputLabel id="filter-app">Application</InputLabel>
                  <Select
                    labelId="filter-app"
                    id="filter-app"
                    label="Application"
                    className={classes.select}
                    value={selectedAppIndex}
                    onChange={(e) => {
                      const appIndex = e.target.value as number;
                      setSelectedApp(filteredApps[appIndex]);
                      setSelectedAppIndex(appIndex);
                      setActiveStep(2);
                    }}
                  >
                    {filteredApps.map((app, i) => (
                      <MenuItem
                        value={i}
                        key={`app-${i}-${app.name}-${app.repoId}`}
                      >
                        name: {app.name}, repo: {app.repoId}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </StepContent>
            </Step>
            <Step key="Confirm information before adding">
              <StepLabel>Confirm information before adding</StepLabel>
              <StepContent>
                {selectedApp && (
                  <Typography className={classes.accordionDetail}>
                    <div className={classes.inputGroup}>
                      <TextField
                        id={"kind"}
                        label="Kind"
                        margin="dense"
                        fullWidth
                        variant="outlined"
                        value={APPLICATION_KIND_TEXT[selectedApp.kind]}
                        className={classes.textInput}
                        inputProps={{ readOnly: true }}
                      />
                    </div>
                    <div className={classes.inputGroup}>
                      <TextField
                        id={"path"}
                        label="Path"
                        margin="dense"
                        variant="outlined"
                        value={selectedApp.path}
                        className={classes.textInput}
                        inputProps={{ readOnly: true }}
                      />
                      <div className={classes.inputGroupSpace} />
                      <TextField
                        id={"configFilename-"}
                        label="Config Filename"
                        margin="dense"
                        variant="outlined"
                        value={selectedApp.configFilename}
                        className={classes.textInput}
                        inputProps={{ readOnly: true }}
                      />
                    </div>
                    {selectedApp.labelsMap.map((label, j) => (
                      <div className={classes.inputGroup} key={label[0]}>
                        <TextField
                          id={"label-" + "-" + j}
                          label={"Label " + j}
                          margin="dense"
                          variant="outlined"
                          value={label[0] + ": " + label[1]}
                          className={classes.textInput}
                          inputProps={{ readOnly: true }}
                        />
                      </div>
                    ))}
                  </Typography>
                )}
              </StepContent>
            </Step>
          </Stepper>

          {selectedApp ? (
            <Button
              color="primary"
              type="submit"
              onClick={() => {
                setAppToAdd({
                  name: selectedApp.name,
                  pipedId: selectedApp.pipedId,
                  repo: {
                    id: selectedApp.repoId,
                  } as ApplicationGitRepository.AsObject,
                  repoPath: selectedApp.path,
                  configFilename: selectedApp.configFilename,
                  kind: selectedApp.kind,
                  platformProvider: selectedPlatformProvider,
                  labels: selectedApp.labelsMap,
                });
                setShowConfirm(true);
              }}
            >
              {UI_TEXT_SAVE}
            </Button>
          ) : (
            <Button color="primary" type="submit" disabled>
              {UI_TEXT_SAVE}
            </Button>
          )}
          <Button onClick={onClose}>{UI_TEXT_CANCEL}</Button>
        </Box>
        <Dialog open={showConfirm}>
          <DialogTitle>{ADD_FROM_GIT_CONFIRM_DIALOG_TITLE}</DialogTitle>
          <DialogContent>
            {ADD_FROM_GIT_CONFIRM_DIALOG_DESCRIPTION}
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setShowConfirm(false)}>
              {UI_TEXT_CANCEL}
            </Button>
            <Button
              color="primary"
              onClick={() => {
                dispatch(addApplication(appToAdd));
                setShowConfirm(false);
                onClose();
              }}
            >
              {UI_TEXT_SAVE}
            </Button>
          </DialogActions>
        </Dialog>
      </>
    );
  }
);
