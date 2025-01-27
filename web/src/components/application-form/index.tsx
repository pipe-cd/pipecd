import { Box, makeStyles, Tabs, Tab, IconButton } from "@material-ui/core";
import { Help } from "@material-ui/icons";
import { useState } from "react";
import { Application, ApplicationKind } from "~/modules/applications";
import ApplicationFormV1 from "./application-form-v1";
import ApplicationFormV0 from "./application-form-v0";
import ApplicationFormManual from "./application-form-manual";
import TabPanel from "./tab-panel";

const useStyles = makeStyles(() => ({
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

const tabProps = (tabKey: number): { id: string; "aria-controls": string } => {
  return {
    id: `tab-${tabKey}`,
    "aria-controls": `tabpanel-${tabKey}`,
  };
};

const tabPanelProps = (
  tabKey: number
): { id: string; "aria-labelledby": string } => {
  return {
    id: `tabpanel-${tabKey}`,
    "aria-labelledby": `tab-${tabKey}`,
  };
};

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

const FORM_SUGGESTION_DOC_URL =
  "https://pipecd.dev/docs/user-guide/managing-application/adding-an-application/#picking-from-a-list-of-unused-apps-suggested-by-pipeds";

export type ApplicationFormProps = {
  title: string;
  onClose: () => void;
  onFinished: () => void;
  disableApplicationInfo?: boolean;
  setIsFormDirty?: (state: boolean) => void;
  setIsSubmitting?: (state: boolean) => void;
  detailApp?: Application.AsObject;
};

enum TabKeys {
  V0 = 0,
  V1 = 1,
  MANUAL = 2,
}

export const ApplicationFormTabs: React.FC<ApplicationFormProps> = (props) => {
  const classes = useStyles();

  const [selectedTabIndex, setSelectedTabIndex] = useState(TabKeys.V0);

  const handleChange = (
    _event: React.ChangeEvent<Record<string, unknown>>,
    newValue: number
  ): void => {
    setSelectedTabIndex(newValue);
    props.setIsFormDirty?.(false);
    props.setIsSubmitting?.(false);
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
                href={FORM_SUGGESTION_DOC_URL}
                target="_blank"
                rel="noopener noreferrer"
              >
                <Help fontSize="small" />
              </IconButton>
            }
            {...tabProps(TabKeys.V0)}
          />
          <Tab
            className={classes.tabLabel}
            label="PIPED V1 ADD FROM SUGGESTIONS"
            {...tabProps(TabKeys.V1)}
          />
          <Tab
            className={classes.tabLabel}
            label="ADD MANUALLY"
            icon=" "
            {...tabProps(TabKeys.MANUAL)}
          />
        </Tabs>
      </Box>
      <TabPanel
        selected={selectedTabIndex === TabKeys.V0}
        {...tabPanelProps(TabKeys.V0)}
      >
        <ApplicationFormV0 {...props} />
      </TabPanel>
      <TabPanel
        selected={selectedTabIndex === TabKeys.V1}
        {...tabPanelProps(TabKeys.V1)}
      >
        <ApplicationFormV1 {...props} />
      </TabPanel>
      <TabPanel
        selected={selectedTabIndex === TabKeys.MANUAL}
        {...tabPanelProps(TabKeys.MANUAL)}
      >
        <ApplicationFormManual {...props} />
      </TabPanel>
    </Box>
  );
};
