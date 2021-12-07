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
  Accordion,
  AccordionSummary,
  AccordionDetails, FormControl, InputLabel, Select,
} from "@material-ui/core";
import ExpandMore from '@material-ui/icons/ExpandMore';
import { FormikProps } from "formik";
import { FC, memo, ReactElement, useCallback, useState } from "react";
import * as yup from "yup";
import { APPLICATION_KIND_TEXT, APPLICATION_KIND_BY_NAME } from "~/constants/application-kind";
import { UI_TEXT_CANCEL, UI_TEXT_SAVE, UI_TEXT_ADD } from "~/constants/ui-text";
import { useAppSelector } from "~/hooks/redux";
import { ApplicationKind } from "~/modules/applications";
import { selectAllEnvs } from "~/modules/environments";
import { Piped, selectAllPipeds, selectPipedById, selectPipedsByEnv } from "~/modules/pipeds";
import { ApplicationInfo } from "~/modules/unregistered-applications";

const createCloudProviderListFromPiped = ({
  kind,
  piped,
}: {
  piped?: Piped.AsObject;
  kind: ApplicationKind;
}): Array<{ name: string; value: string }> => {
  if (!piped) {
    return [{ name: "None", value: "" }];
  }

  return piped.cloudProvidersList
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
    return [{
      name: "None",
      value: "",
      branch: "",
      remote: "",
    }];
  }

  return piped.repositoriesList.map((repo) => ({
    name: repo.id,
    value: repo.id,
    branch: repo.branch,
    remote: repo.remote,
  }));
}

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

function a11yProps(index: number) {
  return {
    id: `simple-tab-${index}`,
    'aria-controls': `simple-tabpanel-${index}`,
  };
}

export const ApplicationFormTabs: React.FC<ApplicationFormProps> = (props) => {
  const [selectedTabIndex, setSelectedTabIndex] = useState(0);

  const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
    setSelectedTabIndex(newValue);
  };

  return (
    <Box width={600}>
      <Box>
        <Tabs value={selectedTabIndex} onChange={handleChange} aria-label="basic tabs example">
          <Tab label="Add manually" {...a11yProps(0)} />
          <Tab label="Add from Git (Alpha)" {...a11yProps(1)} />
        </Tabs>
      </Box>
      <TabPanel selected={selectedTabIndex === 0} index={0}>
        <ApplicationForm {...props} />
      </TabPanel>
      <TabPanel selected={selectedTabIndex === 1} index={1}>
        <UnregisteredApplicationList {...props} />
      </TabPanel>
    </Box>
  );
}

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

export const validationSchema = yup.object().shape({
  name: yup.string().required(),
  kind: yup.number().required(),
  // TODO: Make all environment fields in the form in optional
  env: yup.string().required(),
  pipedId: yup.string().required(),
  repo: yup
    .object({
      id: yup.string().required(),
      remote: yup.string().required(),
      branch: yup.string().required(),
    })
    .required(),
  repoPath: yup.string().required(),
  cloudProvider: yup.string().required(),
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

export type ApplicationFormProps = FormikProps<ApplicationFormValue> & {
  title: string;
  onClose: () => void;
  onAddFromGit: () => void;
  disableGitPath?: boolean;
};

export const emptyFormValues: ApplicationFormValue = {
  name: "",
  env: "",
  kind: ApplicationKind.KUBERNETES,
  pipedId: "",
  repoPath: "",
  configFilename: ".pipe.yaml",
  cloudProvider: "",
  repo: {
    id: "",
    remote: "",
    branch: "",
  },
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
    disableGitPath = false,
  }) {
    const classes = useStyles();

    const environments = useAppSelector(selectAllEnvs);

    const pipeds = useAppSelector<Piped.AsObject[]>((state) =>
      values.env !== "" ? selectPipedsByEnv(state.pipeds, values.env) : []
    );

    const selectedPiped = useAppSelector(selectPipedById(values.pipedId));

    const cloudProviders = createCloudProviderListFromPiped({
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
              items={repositories}
              disabled={
                selectedPiped === undefined ||
                repositories.length === 0 ||
                isSubmitting ||
                disableGitPath
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
                selectedPiped === undefined || isSubmitting || disableGitPath
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
            disabled={
              selectedPiped === undefined || isSubmitting || disableGitPath
            }
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


interface UnregisteredApplicationsFilterOptions {
  pipedId: string;
  cloudProvider: string;
  kind: string;
}

interface UnregisteredApplicationFilterProps {
  onChange: (options: UnregisteredApplicationsFilterOptions) => void;
}

const UnregisteredApplicationFilter: FC<UnregisteredApplicationFilterProps> = memo(
  function UnregisteredApplicationFilter({ onChange }) {
    const classes = useStyles();
    const ps = useAppSelector((state) =>
      selectAllPipeds(state)
    );
    const pipeds = ps.filter((piped) => !piped.disabled)

    const [selectedPipedId, setSelectedPipedId] = useState("")
    const selectedPiped = useAppSelector(selectPipedById(selectedPipedId));
    const cloudProviders = selectedPiped ? selectedPiped.cloudProvidersList : [];

    var options: UnregisteredApplicationsFilterOptions;
    const handleUpdateFilterValue = (
      optionPart: Partial<UnregisteredApplicationsFilterOptions>
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
            label="PipedId"
            className={classes.select}
            onChange={(e) => {
              setSelectedPipedId(e.target.value as string);
              handleUpdateFilterValue({
                pipedId: (e.target.value as string)
              });
            }}
          >
            {pipeds.map((e) => (
              <MenuItem value={e.id} key={`piped-${e.id}`}>
                {e.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <div className={classes.inputGroupSpace} />
        <FormControl className={classes.formItem} variant="outlined">
          <InputLabel id="filter-cloud-provider">Cloud Provider</InputLabel>
          <Select
            labelId="filter-cloud-provider"
            id="filter-cloud-provider"
            label="CloudProvider"
            className={classes.select}
            disabled={selectedPipedId === ""}
            onChange={(e) => {
              const values = (e.target.value as ReadonlyArray<string>);
              handleUpdateFilterValue({
                cloudProvider: values[0],
                kind: values[1],
                pipedId: selectedPipedId,
              });
            }}
          >
            {cloudProviders.map((e) => (
              <MenuItem value={([e.name, e.type] as ReadonlyArray<string>)} key={`cloud-provider-${e.name}`}>
                {e.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </div>
    );
  }
);

const UnregisteredApplicationList: FC<ApplicationFormProps> = memo(
  function ApplicationForm({
    onClose,
    onAddFromGit,
  }) {
    const classes = useStyles();
    // TODO: Call api to fetch all unregisteredApps
    /*const apps = useAppSelector<ApplicationInfo.AsObject[]>((state) =>
        fetchUnregisteredApplications()
    );*/
    const apps: ApplicationInfo.AsObject[] = [];

    const [selectedPipedId, setSelectedPipedId] = useState("");
    const [selectedKind, setSelectedKind] = useState("");

    const handleFilterChange = useCallback(
      (options: UnregisteredApplicationsFilterOptions) => {
        setSelectedPipedId(options.pipedId);
        setSelectedKind(options.kind);
        // TODO: Set cloud provider as well
      },
      [],
    );
    return (
      <Box width="100%">
        <Typography className={classes.title} variant="h6">
          Select the application to register
        </Typography>
        <Divider />
        <UnregisteredApplicationFilter
          onChange={handleFilterChange}
        />
        <Accordion disabled>
          <AccordionSummary
            aria-controls="table-header-content"
            id="table-header"
          >
            <Typography>Name (Repository)</Typography>
          </AccordionSummary>
        </Accordion>
        <div>
          {apps
            .filter(app =>
              app.pipedId === selectedPipedId &&
              app.kind === APPLICATION_KIND_BY_NAME[selectedKind]
            )
            .map((app, i) => (
              <Accordion>
                <AccordionSummary
                  expandIcon={<ExpandMore />}
                  aria-controls={"panel-" + i + "-content"}
                  id={"panel-" + i + "-header"}
                >
                  <Typography>{app.name} ({app.repoId})</Typography>
                </AccordionSummary>
                <AccordionDetails>
                  <Typography>
                    <div className={classes.inputGroup}>
                      <TextField
                        id={"kind-" + i}
                        label="Kind"
                        margin="dense"
                        fullWidth
                        variant="outlined"
                        disabled
                        value={APPLICATION_KIND_TEXT[app.kind]}
                        className={classes.textInput}
                      />
                    </div>
                    <div className={classes.inputGroup}>
                      <TextField
                        id={"path-" + i}
                        label="Path"
                        margin="dense"
                        variant="outlined"
                        disabled
                        value={app.path}
                        className={classes.textInput}
                      />
                      <div className={classes.inputGroupSpace} />
                      <TextField
                        id={"configFilename-" + i}
                        label="Config Filename"
                        margin="dense"
                        variant="outlined"
                        disabled
                        value={app.configFilename}
                        className={classes.textInput}
                      />
                    </div>
                    {app.labelsMap.map((label, j) => (
                      <div className={classes.inputGroup}>
                        <TextField
                          id={"label-" + i + "-" + j}
                          label={"Label " + j}
                          margin="dense"
                          variant="outlined"
                          disabled
                          value={label[0] + ": " + label[1]}
                          className={classes.textInput}
                        />
                      </div>
                    ))}
                    <Button
                      color="primary"
                      type="submit"
                      onClick={onAddFromGit}
                    >
                      {UI_TEXT_ADD}
                    </Button>
                  </Typography>
                </AccordionDetails>
              </Accordion>
            ))}
        </div>
        <Button onClick={onClose}>
          {UI_TEXT_CANCEL}
        </Button>
      </Box>
    );
  }
);
