import { Box, Typography } from "@mui/material";
import { FC } from "react";

type Props = {
  data: {
    key: string;
    title: string;
    description?: string;
    color: string;
  }[];
};

const LegendRow: FC<Props> = ({ data }) => {
  return (
    <Box
      display={"flex"}
      justifyContent={"center"}
      alignItems={"center"}
      columnGap={"10px"}
      mt={2}
      width={"100%"}
    >
      {data.map((item) => (
        <Box
          display={"flex"}
          alignItems={"center"}
          columnGap={"10px"}
          width={"fit-content"}
          key={item.key}
        >
          <Box
            bgcolor={item.color}
            component={"span"}
            sx={{
              width: "30px",
              height: "30px",
              borderRadius: "50%",
            }}
          />
          <Box>
            <Typography variant="body2">{item.title}</Typography>
            {item.description ? (
              <Typography variant="caption">{item.description}</Typography>
            ) : null}
          </Box>
        </Box>
      ))}
    </Box>
  );
};

export default LegendRow;
