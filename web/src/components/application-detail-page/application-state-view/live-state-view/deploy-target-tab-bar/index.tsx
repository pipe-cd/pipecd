import { Tab, Tabs } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { grey } from "@mui/material/colors";
import { FC } from "react";

type Props = {
  tabs: string[];
  selectedTab: string;
  onSelectTab: (tab: string) => void;
};

const useStyles = makeStyles((theme) => ({
  rootTabs: {
    minHeight: 10,
  },
  tab: {
    minHeight: 10,
    padding: 0,
    cursor: "pointer",
    color: grey[500],
    backgroundColor: grey[200],
    border: `1px solid ${grey[200]}`,
    minWidth: 100,
  },
  activeTab: {
    backgroundColor: theme.palette.background.paper,
    color: theme.palette.text.primary + " !important",
  },
  tabsIndicator: {
    backgroundColor: grey[900] + " !important",
  },
}));

const DeployTargetTabBar: FC<Props> = ({ tabs, selectedTab, onSelectTab }) => {
  const classes = useStyles();

  return (
    <Tabs
      value={tabs.indexOf(selectedTab) >= 0 ? tabs.indexOf(selectedTab) : 0}
      onChange={(e, activeIndex) => onSelectTab(tabs[activeIndex])}
      variant="scrollable"
      aria-label="icon label tabs"
      classes={{
        root: classes.rootTabs,
        indicator: classes.tabsIndicator,
      }}
    >
      {tabs.map((tab) => (
        <Tab
          key={tab}
          label={tab}
          classes={{
            root: classes.tab,
            selected: classes.activeTab,
          }}
        />
      ))}
    </Tabs>
  );
};

export default DeployTargetTabBar;
