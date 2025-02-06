import { Box, Typography } from "@material-ui/core";
import React from "react";
import { PropsWithChildren } from "react";

type TabPanelProps = PropsWithChildren<{
  selected: boolean;
  id: string;
  ariaLabelledBy?: string;
}>;

const TabPanel: React.FC<TabPanelProps> = (props) => {
  return (
    <div
      role="tabpanel"
      hidden={!props.selected}
      id={props.id}
      aria-labelledby={props.ariaLabelledBy}
    >
      {props.selected && (
        <Box>
          <Typography>{props.children}</Typography>
        </Box>
      )}
    </div>
  );
};

export default TabPanel;
