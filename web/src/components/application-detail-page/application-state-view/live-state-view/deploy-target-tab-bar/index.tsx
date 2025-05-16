import { Tab, Tabs } from "@mui/material";
import { grey } from "@mui/material/colors";
import { FC } from "react";

type Props = {
  tabs: string[];
  selectedTab: string;
  onSelectTab: (tab: string) => void;
};

const DeployTargetTabBar: FC<Props> = ({ tabs, selectedTab, onSelectTab }) => {
  return (
    <Tabs
      value={tabs.indexOf(selectedTab) >= 0 ? tabs.indexOf(selectedTab) : 0}
      onChange={(e, activeIndex) => onSelectTab(tabs[activeIndex])}
      variant="scrollable"
      aria-label="icon label tabs"
      sx={{ minHeight: 10 }}
      TabIndicatorProps={{
        sx: { backgroundColor: grey[900] },
      }}
    >
      {tabs.map((tab) => (
        <Tab
          key={tab}
          label={tab}
          sx={(theme) => ({
            minHeight: "10px",
            padding: "2px 10px",
            cursor: "pointer",
            color: grey[500],
            backgroundColor: grey[200],
            border: `1px solid ${grey[200]}`,
            minWidth: 100,
            "&.Mui-selected": {
              backgroundColor: theme.palette.background.paper,
              color: theme.palette.text.primary,
            },
          })}
        />
      ))}
    </Tabs>
  );
};

export default DeployTargetTabBar;
