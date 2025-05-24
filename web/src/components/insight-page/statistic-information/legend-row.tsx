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
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        columnGap: "10px",
        mt: 2,
        width: "100%",
      }}
    >
      {data.map((item) => (
        <Box
          key={item.key}
          sx={{
            display: "flex",
            alignItems: "center",
            columnGap: "10px",
            width: "fit-content",
          }}
        >
          <Box
            component={"span"}
            sx={{
              bgcolor: item.color,
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
