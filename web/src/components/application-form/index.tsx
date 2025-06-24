import { Box, Tabs, Tab, IconButton } from "@mui/material";
import { Help } from "@mui/icons-material";
import { useState } from "react";
import { Application } from "~/modules/applications";
import ApplicationFormV1 from "./application-form-v1";
import ApplicationFormV0 from "./application-form-v0";
import ApplicationFormManualV0 from "./application-form-manual-v0";
import TabPanel from "./tab-panel";

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

const FORM_SUGGESTION_DOC_URL =
  "https://pipecd.dev/docs/user-guide/managing-application/adding-an-application/#picking-from-a-list-of-unused-apps-suggested-by-pipeds";

export type ApplicationFormProps = {
  title: string;
  onClose: () => void;
  onFinished: () => void;
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
  const [selectedTabIndex, setSelectedTabIndex] = useState(TabKeys.V0);

  const handleChange = (_event: unknown, newValue: number): void => {
    setSelectedTabIndex(newValue);
    props.setIsFormDirty?.(false);
    props.setIsSubmitting?.(false);
  };

  return (
    <Box
      sx={{
        width: 600,
      }}
    >
      <Box>
        <Tabs
          value={selectedTabIndex}
          onChange={handleChange}
          aria-label="basic tabs example"
        >
          <Tab
            label="PIPED V0 ADD FROM SUGGESTIONS"
            iconPosition="end"
            sx={{ maxWidth: "210px" }}
            icon={
              <IconButton
                size="small"
                href={FORM_SUGGESTION_DOC_URL}
                target="_blank"
                rel="noopener noreferrer"
                sx={{ marginLeft: "0px !important" }}
              >
                <Help fontSize="small" />
              </IconButton>
            }
            {...tabProps(TabKeys.V0)}
          />
          <Tab
            sx={{ maxWidth: "180px" }}
            label="PIPED V1 ADD FROM SUGGESTIONS"
            {...tabProps(TabKeys.V1)}
          />
          <Tab
            sx={{ maxWidth: "210px" }}
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
        <ApplicationFormManualV0 {...props} />
      </TabPanel>
    </Box>
  );
};
