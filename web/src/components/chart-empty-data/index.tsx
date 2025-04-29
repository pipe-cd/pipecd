import { Box, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { WarningOutlined } from "@mui/icons-material";
import { FC } from "react";

type Props = {
  visible: boolean;
  noDataText?: string;
};

const useStyles = makeStyles((theme) => ({
  noDataMessage: {
    display: "flex",
  },
  noDataMessageIcon: {
    marginRight: theme.spacing(1),
  },
}));

const NO_DATA_TEXT = "No data is available.";

const ChartEmptyData: FC<Props> = ({ visible, noDataText = NO_DATA_TEXT }) => {
  const classes = useStyles();

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
        className={classes.noDataMessage}
      >
        <WarningOutlined className={classes.noDataMessageIcon} />
        {noDataText}
      </Typography>
    </Box>
  );
};

export default ChartEmptyData;
