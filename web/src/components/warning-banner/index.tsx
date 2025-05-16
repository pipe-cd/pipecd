import { FC, memo } from "react";
import { AppBar } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import CloseIcon from "@mui/icons-material/Close";

const useStyles = makeStyles((theme) => ({
  container: {
    zIndex: theme.zIndex.drawer + 1,
    height: "30px",
    backgroundColor: "#403f4c",
  },
  content: {
    paddingTop: "2px",
    textAlign: "center",
    fontSize: "1rem",
  },
  highlight: {
    backgroundColor: "yellow",
  },
  close: {
    position: "absolute",
    top: "2px",
    right: "2px",
  },
}));

export interface WarningBannerProps {
  onClose: () => void;
}

export const WarningBanner: FC<WarningBannerProps> = memo(
  function BannerWarning({ onClose }) {
    const classes = useStyles();
    const releaseNoteURL = `https://github.com/pipe-cd/pipecd/releases/tag/${process.env.PIPECD_VERSION}`;

    return (
      <AppBar position="static" className={classes.container}>
        <div className={classes.content}>
          Piped version{" "}
          <a
            href={releaseNoteURL}
            target="_blank"
            rel="noreferrer"
            className={classes.highlight}
          >
            {process.env.PIPECD_VERSION}
          </a>{" "}
          is available!
        </div>
        <CloseIcon className={classes.close} onClick={onClose} />
      </AppBar>
    );
  }
);
