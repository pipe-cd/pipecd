import { FC, memo } from "react";
import { AppBar, Box, Link } from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";

export interface WarningBannerProps {
  onClose: () => void;
}

export const WarningBanner: FC<WarningBannerProps> = memo(
  function BannerWarning({ onClose }) {
    const releaseNoteURL = `https://github.com/pipe-cd/pipecd/releases/tag/${process.env.PIPECD_VERSION}`;

    return (
      <AppBar
        position="static"
        sx={(theme) => ({
          zIndex: theme.zIndex.drawer - 1,
          height: "30px",
          backgroundColor: "#403f4c",
        })}
      >
        <Box
          sx={{
            paddingTop: "2px",
            textAlign: "center",
            fontSize: "1rem",
          }}
        >
          Piped version{" "}
          <Link
            href={releaseNoteURL}
            target="_blank"
            rel="noreferrer"
            sx={{
              bgcolor: "yellow",
            }}
          >
            {process.env.PIPECD_VERSION}
          </Link>{" "}
          is available!
        </Box>
        <CloseIcon
          sx={{
            position: "absolute",
            top: "2px",
            right: "2px",
          }}
          onClick={onClose}
        />
      </AppBar>
    );
  }
);
