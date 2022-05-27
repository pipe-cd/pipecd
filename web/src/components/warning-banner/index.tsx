import { FC, memo } from "react";
import { AppBar, makeStyles } from "@material-ui/core";
import CloseIcon from "@material-ui/icons/Close";

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
          Version {process.env.PIPECD_VERSION} is available! Find all{" "}
          <a
            href={releaseNoteURL}
            target="_blank"
            rel="noreferrer"
            className={classes.highlight}
          >
            details on this release note here.
          </a>
        </div>
        <CloseIcon className={classes.close} onClick={onClose} />
      </AppBar>
    );
  }
);
