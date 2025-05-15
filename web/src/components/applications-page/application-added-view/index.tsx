import { Box, Button, Divider, Typography } from "@mui/material";
import { FC, memo } from "react";

const TEXT = {
  TITLE: "Congratulation!",
  MESSAGE: "Your application has been added successfully.",
  NOTE:
    "Please ensure that your application directory in Git is containing an application config file since PipeCD needs it to know how to run the application's deployments.",
};

export interface ApplicationAddedViewProps {
  onClose: () => void;
}

export const ApplicationAddedView: FC<ApplicationAddedViewProps> = memo(
  function ApplicationAddedView({ onClose }) {
    return (
      <Box width={600} flex={1} display="flex" flexDirection="column">
        <Typography variant="h6" p={2}>
          {TEXT.TITLE}
        </Typography>

        <Divider />

        <Box p={2}>
          <Box mt={2} mb={2}>
            <Typography variant="subtitle1">{TEXT.MESSAGE}</Typography>
            <Typography
              variant="body2"
              sx={{ marginTop: 1, color: "text.secondary" }}
            >
              {TEXT.NOTE}
            </Typography>
          </Box>

          <Box mt={1} textAlign="right">
            <Button onClick={onClose} variant="outlined">
              CLOSE
            </Button>
          </Box>
        </Box>
      </Box>
    );
  }
);
