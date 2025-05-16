import { Box, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import { FC } from "react";

type Props = {
  data: {
    key: string;
    title: string;
    description?: string;
    color: string;
  }[];
};

const useStyles = makeStyles(() => ({
  root: {
    minWidth: 300,
    display: "inline-block",
    overflow: "visible",
    position: "relative",
  },
  pageTitle: {
    fontWeight: "bold",
  },
  labelDot: {
    width: 30,
    height: 30,
    borderRadius: "50%",
  },
}));

const LegendRow: FC<Props> = ({ data }) => {
  const classes = useStyles();

  return (
    <Box
      display={"flex"}
      justifyContent={"center"}
      alignItems={"center"}
      columnGap={10}
      mt={2}
      width={"100%"}
    >
      {data.map((item) => (
        <Box
          display={"flex"}
          alignItems={"center"}
          columnGap={10}
          key={item.key}
        >
          <Box className={classes.labelDot} bgcolor={item.color} />
          <div>
            <Typography variant="body2">{item.title}</Typography>
            {item.description ? (
              <Typography variant="caption">{item.description}</Typography>
            ) : null}
          </div>
        </Box>
      ))}
    </Box>
  );
};

export default LegendRow;
