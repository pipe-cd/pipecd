import { Box, Typography } from "@mui/material";
import { WarningOutlined } from "@mui/icons-material";
import { FC } from "react";

type Props = {
  visible: boolean;
  noDataText?: string;
};

const NO_DATA_TEXT = "No data is available.";

const ChartEmptyData: FC<Props> = ({ visible, noDataText = NO_DATA_TEXT }) => {
  return (
    <Box
      display={visible ? "flex" : "none"}
      width="100%"
      height="100%"
      alignItems="center"
      justifyContent="center"
      position="absolute"
      top={0}
      left={0}
      bgcolor="#fafafabb"
    >
      <Typography
        variant="body1"
        color="textSecondary"
        // className={classes.noDataMessage}
        display="flex"
      >
        <WarningOutlined
          // className={classes.noDataMessageIcon}
          sx={{ mr: 1 }}
        />
        {noDataText}
      </Typography>
    </Box>
  );
};

export default ChartEmptyData;
