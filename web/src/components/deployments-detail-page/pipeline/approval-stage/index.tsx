import { Paper, Typography } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import WaitIcon from "@mui/icons-material/PauseCircleOutline";
import { FC, memo } from "react";

const useStyles = makeStyles((theme) => ({
  root: (props: { active: boolean }) => ({
    flex: 1,
    display: "inline-flex",
    cursor: "pointer",
    padding: theme.spacing(2),
    "&:hover": {
      backgroundColor: theme.palette.action.hover,
    },
    backgroundColor: props.active
      ? // NOTE: 12%
        theme.palette.primary.main + "1e"
      : undefined,
  }),
  icon: {
    color: theme.palette.warning.main,
  },
  name: {
    marginLeft: theme.spacing(1),
  },
  stageName: {
    fontFamily: theme.typography.fontFamilyMono,
  },
  main: {
    display: "flex",
    justifyContent: "flex-start",
    alignItems: "center",
  },
}));

export interface ApprovalStageProps {
  id: string;
  name: string;
  active: boolean;
  onClick: (stageId: string, stageName: string) => void;
}

export const ApprovalStage: FC<ApprovalStageProps> = memo(
  function ApprovalStage({ id, name, onClick, active }) {
    const classes = useStyles({ active });

    function handleOnClick(): void {
      onClick(id, name);
    }

    return (
      <Paper square className={classes.root} onClick={handleOnClick}>
        <div className={classes.main}>
          <WaitIcon className={classes.icon} />
          <Typography variant="subtitle2" className={classes.name}>
            <div className={classes.stageName}>{name}</div>
          </Typography>
        </div>
      </Paper>
    );
  }
);
