import React, { FC, memo } from "react";
import { makeStyles, Paper, Typography, Box } from "@material-ui/core";
import WaitIcon from "@material-ui/icons/PauseCircleOutline";

const useStyles = makeStyles((theme) => ({
  root: (props: { active: boolean }) => ({
    display: "inline-flex",
    cursor: "pointer",
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
}));

interface Props {
  id: string;
  name: string;
  active: boolean;
  onClick: (stageId: string, stageName: string) => void;
}

export const ApprovalStage: FC<Props> = memo(function ApprovalStage({
  id,
  name,
  onClick,
  active,
}) {
  const classes = useStyles({ active });

  function handleOnClick(): void {
    onClick(id, name);
  }

  return (
    <Paper square className={classes.root} onClick={handleOnClick}>
      <Box alignItems="center" display="flex" justifyContent="center" p={2}>
        <WaitIcon className={classes.icon} />
        <Typography variant="subtitle2" className={classes.name}>
          <Box fontFamily="Roboto Mono">{name}</Box>
        </Typography>
      </Box>
    </Paper>
  );
});
