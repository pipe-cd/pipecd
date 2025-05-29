import { FC, memo } from "react";
import ApplicationCount from "./application-count";
import ApplicationByPiped from "./application-by-piped";
import PipedCount from "./piped-count";
import Deployment24h from "./deployment-24h";
import { Box } from "@mui/material";

export const StatisticInformation: FC = memo(function StatisticInformation() {
  return (
    <Box
      sx={(theme) => ({
        width: "fit-content",
        margin: "0 auto",
        display: "grid",
        gap: 2,
        gridTemplateColumns: "repeat(2, 1fr)",
        [theme.breakpoints.up("xl")]: {
          gridTemplateColumns: "repeat(4, 1fr)",
        },
      })}
    >
      <ApplicationCount />
      <ApplicationByPiped />
      <PipedCount />
      <Deployment24h />
    </Box>
  );
});
